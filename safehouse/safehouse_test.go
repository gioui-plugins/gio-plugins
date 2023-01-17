package safehouse

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"runtime"
	"testing"
)

func TestSafeHouse_Set(t *testing.T) {
	sh := NewSafeHouse()

	secret := genRandomSecret(512)
	if err := sh.Set(secret); err != nil {
		t.Error(err)
		return
	}

	if x, err := sh.Get(secret.Identifier); err != nil || !bytes.Equal(x.Data, secret.Data) {
		if err == nil {
			err = errors.New("not equal")
		}
		t.Error(err)
		return
	}

	if runtime.GOOS == "JS" {
		return
	}

	if err := sh.Remove(secret.Identifier); err != nil {
		t.Error(err)
		return
	}

	if y, err := sh.Get(secret.Identifier); err == nil || len(y.Data) > 0 {
		t.Error(err)
		return
	}

}

func TestSafeHouse_Remove(t *testing.T) {
	sh := NewSafeHouse()

	secret := genRandomSecret(128)
	if err := sh.Remove(secret.Identifier); err == nil {
		t.Error(err)
	}

	secret = genRandomSecret(128)
	if err := sh.Set(secret); err != nil {
		t.Error(err)
	}

	if r, err := sh.Get(secret.Identifier); err != nil || !bytes.Equal(r.Data, secret.Data) {
		if err == nil {
			errors.New("empty data")
		}
		t.Error(err)
	}

	if runtime.GOOS == "JS" {
		return
	}

	if err := sh.Remove(secret.Identifier); err != nil {
		t.Error(err)
	}

	if _, err := sh.Get(secret.Identifier); err == nil {
		t.Error("not deleted")
	}
}

func TestSafeHouse_Get(t *testing.T) {
	sh := NewSafeHouse()
	secret := genRandomSecret(512)

	if _, err := sh.Get(secret.Identifier); err == nil {
		t.Error(err)
	}
}

func genRandomSecret(n int) Secret {
	return Secret{
		Identifier:  genRandomData[string](64),
		Description: genRandomData[string](128),
		Data:        genRandomData[[]byte](n),
	}
}

func genRandomData[T string | []byte](n int) (r T) {
	b := make([]byte, n)
	rand.Read(b)
	switch ((interface{})(r)).(type) {
	case string:
		return T(base64.URLEncoding.EncodeToString(b))
	case []byte:
		return T(b)
	default:
		return T("")
	}
}
