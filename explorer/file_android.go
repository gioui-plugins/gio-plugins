package explorer

import (
	"errors"
	"io"
	"runtime"
	"sync"

	"git.wow.st/gmp/jni"
)

//go:generate javac -source 8 -target 8  -bootclasspath $ANDROID_HOME/platforms/android-30/android.jar -d explorer_file_android/classes file_android.java
//go:generate jar cf file_android.jar -C explorer_file_android/classes .

type File struct {
	VM 		  uintptr
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

func initFileManager(env jni.Env, ctx uintptr) error {
	class, err := jni.LoadClass(env, jni.ClassLoaderFor(env, jni.Object(ctx)), "org/gioui/x/explorer/file_android")
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

func newFile(env jni.Env, stream jni.Object, ctx, vm uintptr) (io.ReadWriteCloser, error) {
	fileManagerOnce.Do(func() {
		initFileManager(env, ctx)
	})

	class, err := jni.LoadClass(env, jni.ClassLoaderFor(env, jni.Object(ctx)), "org/gioui/x/explorer/file_android")
	if err != nil {
		return nil, err
	}

	obj, err := jni.NewObject(env, class, jni.GetMethodID(env, class, "<init>", `()V`))
	if err != nil {
		return nil, err
	}

	// For some reason, using `stream` as argument for a constructor (`public file_android(Object j) {}`) doesn't work.
	if err := jni.CallVoidMethod(env, obj, jni.GetMethodID(env, class, "setHandle", `(Ljava/lang/Object;)V`), jni.Value(stream)); err != nil {
		return nil, err
	}

	f := &File{
		VM:		   vm,
		stream:    jni.NewGlobalRef(env, stream),
		libObject: jni.NewGlobalRef(env, obj),
	}
	runtime.SetFinalizer(f, func(f *File) {
		f.Close()
	})

	return f, nil

}

func (f *File) Read(b []byte) (n int, err error) {
	if f == nil || f.isClosed {
		return 0, io.ErrClosedPipe
	}
	if len(b) == 0 {
		return 0, nil // Avoid unnecessary call to JNI.
	}

	err = jni.Do(jni.JVMFor(f.VM), func(env jni.Env) error {
		if len(b) != f.sharedBufferLen {
			buffer := jni.NewByteArray(env, b)
			defer jni.DeleteLocalRef(env, jni.Object(buffer))
			if f.sharedBuffer != 0 {
				jni.DeleteGlobalRef(env, f.sharedBuffer)
			}
			f.sharedBuffer = jni.NewGlobalRef(env, jni.Object(buffer))
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
	runtime.KeepAlive(f)
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

	err = jni.Do(jni.JVMFor(f.VM), func(env jni.Env) error {
		buffer := jni.NewByteArray(env, b)
		defer jni.DeleteLocalRef(env, jni.Object(buffer))
		ok, err := jni.CallStaticBooleanMethod(env, fileManager, fileMethodWrite, jni.Value(f.libObject), jni.Value(buffer))
		if err != nil {
			return err
		}
		if !ok {
			return f.lastError(env)
		}
		return nil
	})
	runtime.KeepAlive(f)
	if err != nil {
		return 0, err
	}
	return len(b), err
}

func (f *File) Close() error {
	if f == nil || f.isClosed {
		return io.ErrClosedPipe
	}

	err := jni.Do(jni.JVMFor(f.VM), func(env jni.Env) error {
		defer func() {
			f.releaseRefs(env)
			f.isClosed = true
			runtime.SetFinalizer(f, nil)
		}()

		ok, err := jni.CallStaticBooleanMethod(env, fileManager, fileMethodClose, jni.Value(f.libObject))
		if err != nil {
			return err
		}
		if !ok {
			return f.lastError(env)
		}

		return nil
	})
	runtime.KeepAlive(f)
	return err
}

func (f *File) releaseRefs(env jni.Env) {
	if f.stream != 0 {
		jni.DeleteGlobalRef(env, f.stream)
	}
	if f.libObject != 0 {
		jni.DeleteGlobalRef(env, f.libObject)
	}
	if f.sharedBuffer != 0 {
		jni.DeleteGlobalRef(env, f.sharedBuffer)
	}
	f.stream, f.libObject, f.sharedBuffer = 0, 0, 0
	f.sharedBufferLen = 0
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
