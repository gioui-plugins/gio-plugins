package explorer

import (
	"errors"
	"io"
	"runtime"
	"sync"

	"gioui.org/app"
	"git.wow.st/gmp/jni"
)

//go:generate javac -source 8 -target 8  -bootclasspath $ANDROID_HOME/platforms/android-30/android.jar -d explorer_file_android/classes file_android.java
//go:generate jar cf file_android.jar -C explorer_file_android/classes .

type File struct {
	stream    jni.Object
	libObject jni.Object

	sharedBuffer    jni.Object
	sharedBufferLen int
	isClosed        bool
}

var (
	fileManagerOnce sync.Once
	fileManager     jni.Class
	fileMethodRead  jni.MethodID
	fileMethodWrite jni.MethodID
	fileMethodClose jni.MethodID
	fileMethodError jni.MethodID
)

func initFileManager(env jni.Env) error {
	class, err := jni.LoadClass(env, jni.ClassLoaderFor(env, jni.Object(app.AppContext())), "org/gioui/x/explorer/file_android")
	if err != nil {
		return err
	}

	fileManager = jni.Class(jni.NewGlobalRef(env, jni.Object(class)))
	fileMethodRead = jni.GetStaticMethodID(env, fileManager, "fileRead", "(Lorg/gioui/x/explorer/file_android;[B)I")
	fileMethodWrite = jni.GetStaticMethodID(env, fileManager, "fileWrite", "(Lorg/gioui/x/explorer/file_android;[B)Z")
	fileMethodClose = jni.GetStaticMethodID(env, fileManager, "fileClose", "(Lorg/gioui/x/explorer/file_android;)Z")
	fileMethodError = jni.GetStaticMethodID(env, fileManager, "getError", "(Lorg/gioui/x/explorer/file_android;)Ljava/lang/String;")

	return nil
}

func newFile(env jni.Env, stream jni.Object) (*File, error) {
	fileManagerOnce.Do(func() {
		initFileManager(env)
	})

	stream = jni.NewGlobalRef(env, stream)
	f := &File{
		stream: stream,
	}

	runtime.SetFinalizer(f, func(f *File) {
		f.Close()
	})

	class, err := jni.LoadClass(env, jni.ClassLoaderFor(env, jni.Object(app.AppContext())), "org/gioui/x/explorer/file_android")
	if err != nil {
		return nil, err
	}

	obj, err := jni.NewObject(env, class, jni.GetMethodID(env, class, "<init>", `()V`))
	if err != nil {
		return nil, err
	}

	// For some reason, using `f.stream` as argument for a constructor (`public file_android(Object j) {}`) doesn't work.
	if err := jni.CallVoidMethod(env, obj, jni.GetMethodID(env, class, "setHandle", `(Ljava/lang/Object;)V`), jni.Value(f.stream)); err != nil {
		return nil, err
	}

	f.libObject = jni.NewGlobalRef(env, obj)

	return f, nil

}

func (f *File) Read(b []byte) (n int, err error) {
	if f == nil || f.isClosed {
		return 0, io.ErrClosedPipe
	}
	if len(b) == 0 {
		return 0, nil // Avoid unnecessary call to JNI.
	}

	err = jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
		if len(b) != f.sharedBufferLen {
			f.sharedBuffer = jni.Object(jni.NewGlobalRef(env, jni.Object(jni.NewByteArray(env, b))))
			f.sharedBufferLen = len(b)
		}

		size, err := jni.CallStaticIntMethod(env, fileManager, fileMethodRead, jni.Value(f.libObject), jni.Value(f.sharedBuffer))
		if err != nil {
			return err
		}
		if size <= 0 {
			return f.lastError(env)
		}

		n = copy(b, jni.GetByteArrayElements(env, jni.ByteArray(f.sharedBuffer))[:int(size)])
		return nil
	})
	if err != nil {
		return 0, err
	}
	if n == 0 {
		return 0, io.EOF
	}
	return n, err
}

func (f *File) Write(b []byte) (n int, err error) {
	if f == nil || f.isClosed {
		return 0, io.ErrClosedPipe
	}
	if len(b) == 0 {
		return 0, nil // Avoid unnecessary call to JNI.
	}

	err = jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
		ok, err := jni.CallStaticBooleanMethod(env, fileManager, fileMethodWrite, jni.Value(f.libObject), jni.Value(jni.NewByteArray(env, b)))
		if err != nil {
			return err
		}
		if !ok {
			return f.lastError(env)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return len(b), err
}

func (f *File) Close() error {
	if f == nil || f.isClosed {
		return io.ErrClosedPipe
	}

	return jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
		ok, err := jni.CallStaticBooleanMethod(env, fileManager, fileMethodClose, jni.Value(f.libObject))
		if err != nil {
			return err
		}
		if !ok {
			return f.lastError(env)
		}

		f.isClosed = true
		jni.DeleteGlobalRef(env, f.stream)
		jni.DeleteGlobalRef(env, f.libObject)
		if f.sharedBuffer != 0 {
			jni.DeleteGlobalRef(env, f.sharedBuffer)
		}

		return nil
	})
}

func (f *File) lastError(env jni.Env) error {
	message, err := jni.CallStaticObjectMethod(env, fileManager, fileMethodError, jni.Value(f.libObject))
	if err != nil {
		return err
	}
	if err := jni.GoString(env, jni.String(message)); len(err) > 0 {
		return errors.New(err)
	}
	return err
}
