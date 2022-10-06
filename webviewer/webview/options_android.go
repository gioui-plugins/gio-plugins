package webview

import (
	"encoding/base64"

	"git.wow.st/gmp/jni"
)

func (r *driver) setProxy() error {
	options.Lock()
	defer options.Unlock()

	if options.proxy.ip == "" && options.proxy.port == "" {
		return nil
	}

	ok, err := r.callBooleanArgs("webview_proxy", "(Ljava/lang/String;Ljava/lang/String;)Z", func(env jni.Env) []jni.Value {
		return []jni.Value{
			jni.Value(jni.JavaString(env, options.proxy.ip)),
			jni.Value(jni.JavaString(env, options.proxy.port)),
		}
	})

	if err != nil || !ok {
		return ErrInvalidProxy
	}

	return nil
}

func (r *driver) setCerts() error {
	options.Lock()
	defer options.Unlock()

	if options.certs == nil {
		return nil
	}

	var jcerts string
	for _, c := range options.certs {
		jcerts += base64.StdEncoding.EncodeToString(c.Raw) + ";"
	}

	ok, err := r.callBooleanArgs("webview_certs", "(Ljava/lang/String;)Z", func(env jni.Env) []jni.Value {
		return []jni.Value{
			jni.Value(jni.JavaString(env, jcerts)),
		}
	})

	if err != nil || !ok {
		return ErrInvalidCert
	}

	return nil
}
