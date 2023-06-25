package safedata

import (
	"errors"
	"os"
	"path/filepath"

	"git.wow.st/gmp/jni"
)

//go:generate mkdir -p $TEMP/safedata/classes && javac -source 8 -target 8 -bootclasspath $ANDROID_HOME/platforms/android-30/android.jar -d $TEMP/safedata/classes safedata_android.java
//go:generate jar cf safedata_android.jar -C $TEMP/safedata/classes .

type driver struct {
	folder string
	prefix string

	javaVM      uintptr
	javaContext uintptr

	cls        jni.Class
	encryptMid jni.MethodID
	decryptMid jni.MethodID
}

func attachDriver(house *SafeData, config Config) {
	house.driver = driver{}
	configureDriver(&house.driver, config)
}

func configureDriver(driver *driver, config Config) {
	driver.folder = config.Folder
	driver.prefix = config.App

	driver.javaVM = config.VM
	driver.javaContext = config.Context
	initDriver(driver)
}

func initDriver(driver *driver) {
	if driver.javaVM == 0 || driver.javaContext == 0 {
		return
	}

	jni.Do(jni.JVMFor(driver.javaVM), func(env jni.Env) error {
		lib, err := jni.LoadClass(env, jni.ClassLoaderFor(env, jni.Object(driver.javaContext)), "com/inkeliz/safedata_android/safedata_android")
		if err != nil {
			return err
		}

		destroyDriver(driver)
		driver.cls = jni.Class(jni.NewGlobalRef(env, jni.Object(lib)))
		driver.encryptMid = jni.GetStaticMethodID(env, lib, "encrypt", "([BLjava/lang/String;)V")
		driver.decryptMid = jni.GetStaticMethodID(env, lib, "decrypt", "(Ljava/lang/String;)[B")

		return nil
	})
}

func destroyDriver(driver *driver) {
	if driver.cls == 0 {
		return
	}

	jni.DeleteGlobalRef(jni.EnvFor(driver.javaVM), jni.Object(driver.cls))
	driver.cls = 0
	driver.encryptMid = nil
	driver.decryptMid = nil
}

func (d driver) setSecret(secret Secret) error {
	return jni.Do(jni.JVMFor(d.javaVM), func(env jni.Env) error {
		path := jni.JavaString(env, filepath.Join(d.folder, d.keyFor(secret.Identifier)))
		content := jni.NewByteArray(env, secret.Data)

		return jni.CallStaticVoidMethod(env, d.cls, d.encryptMid, jni.Value(content), jni.Value(path))
	})
}

func (d driver) listSecret(looper Looper) error {
	glob, err := filepath.Glob(filepath.Join(d.folder, d.keyFor("*")))
	if err != nil {
		return err
	}

	if len(glob) == 0 {
		return ErrNotFound
	}

	for i := 0; i < len(glob); i++ {
		looper(d.rawKeyFor(filepath.Base(glob[i])))
	}

	return nil
}

func (d driver) getSecret(identifier string, secret *Secret) error {
	if err := d.checkFileExists(identifier); err != nil {
		return err
	}

	return jni.Do(jni.JVMFor(d.javaVM), func(env jni.Env) error {
		path := jni.JavaString(env, filepath.Join(d.folder, d.keyFor(identifier)))

		r, err := jni.CallStaticObjectMethod(env, d.cls, d.decryptMid, jni.Value(path))
		if err != nil {
			return err
		}

		secret.Identifier = identifier
		secret.Description = ""
		secret.Data = jni.GetByteArrayElements(env, jni.ByteArray(r))

		return nil
	})
}

func (d driver) checkFileExists(identifier string) error {
	stat, err := os.Stat(filepath.Join(d.folder, d.keyFor(identifier)))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNotFound
		}
		return err
	}

	if stat.IsDir() {
		return ErrNotFound
	}
	return nil
}

func (d driver) removeSecret(identifier string) error {
	if err := d.checkFileExists(identifier); err != nil {
		return err
	}

	return os.Remove(filepath.Join(d.folder, d.keyFor(identifier)))
}

func (d driver) keyFor(id string) string {
	return d.prefix + id
}

func (d driver) rawKeyFor(id string) string {
	if len(id) < len(d.prefix) {
		return id
	}
	return id[len(d.prefix):]
}
