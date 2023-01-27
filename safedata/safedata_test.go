package safedata

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"testing"
)

func init() {
	for _, n := range []string{"_Test_", "Test"} {
		sh := NewSafeData(Config{App: n})
		sh.List(func(identifier string) (next bool) {
			if err := sh.Remove(identifier); err != nil {
				panic(err)
			}
			return true
		})
	}
}

func TestSafeData_Set(t *testing.T) {
	sh := NewSafeData(Config{App: "Test"})

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

	if err := sh.Remove(secret.Identifier); err != nil {
		t.Error(err)
		return
	}

	if y, err := sh.Get(secret.Identifier); err == nil || len(y.Data) > 0 {
		t.Error(err)
		return
	}
}

func TestSafeData_SetUpdate(t *testing.T) {
	sh := NewSafeData(Config{App: "Test"})

	secret := genRandomSecret(512)
	if err := sh.Set(secret); err != nil {
		t.Error(err)
		return
	}

	secret.Data = genRandomData[[]byte](1024)
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
}

func TestSafeData_SetLongSize(t *testing.T) {
	sh := NewSafeData(Config{App: "Test"})

	secret := genRandomSecret(64)
	secret.Data = genRandomData[[]byte](3_000)

	if err := sh.Set(secret); err != nil {
		t.Error(err)
	}

	if x, err := sh.Get(secret.Identifier); err != nil || !bytes.Equal(x.Data, secret.Data) {
		if err == nil {
			err = errors.New("not equal")
		}
		t.Error(err)
	}
}

func TestSafeData_Remove(t *testing.T) {
	sh := NewSafeData(Config{App: "Test"})

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
			err = errors.New("empty data")
		}
		t.Error(err)
	}

	if err := sh.Remove(secret.Identifier); err != nil {
		t.Error(err)
	}

	if _, err := sh.Get(secret.Identifier); err == nil {
		t.Error("not deleted")
	}
}

func TestSafeData_Get(t *testing.T) {
	sh := NewSafeData(Config{App: "Test"})
	secret := genRandomSecret(512)

	if _, err := sh.Get(secret.Identifier); err == nil {
		t.Error(err)
	}
}

func TestSafeData_ListRemoveAll(t *testing.T) {
	for _, n := range []string{"_Test_", "Test"} {
		sh := NewSafeData(Config{App: n})
		for i := 0; i < 10; i++ {
			sh.Set(genRandomSecret(100))
		}

		err := sh.List(func(identifier string) (next bool) {
			if err := sh.Remove(identifier); err != nil {
				t.Error(err)
			}
			return true
		})

		if err != nil {
			t.Error(err)
		}

		i := 0
		sh.List(func(identifier string) (next bool) {
			i++
			return false
		})

		if i > 0 {
			t.Error("not deleted")
		}
	}

}

func genRandomSecret(n int) Secret {
	return Secret{
		Identifier:  genRandomData[string](64),
		Description: genRandomData[string](92),
		Data:        genRandomData[[]byte](n),
	}
}

func genRandomData[T string | []byte](n int) (r T) {
	switch ((interface{})(r)).(type) {
	case string:
		b := make([]byte, base64.URLEncoding.DecodedLen(n))
		rand.Read(b)
		return T(base64.URLEncoding.EncodeToString(b))
	case []byte:
		b := make([]byte, n)
		rand.Read(b)
		return T(b)
	default:
		return T("")
	}
}
