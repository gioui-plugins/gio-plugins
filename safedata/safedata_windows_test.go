package safedata

import (
	"testing"
)

func TestSafeData_SetInvalidChar(t *testing.T) {
	sh := NewSafeData(Config{App: "Test"})

	secret := genRandomSecret(64)
	secret.Identifier = secret.Identifier[:42] + "\x00" + secret.Identifier[42:]

	if err := sh.Set(secret); err != ErrMalformedMetadata {
		t.Error("must fail, since identifier has invalid char")
	}

	secret = genRandomSecret(64)
	secret.Description = secret.Description[:42] + "\x00" + secret.Description[42:]

	if err := sh.Set(secret); err != ErrMalformedMetadata {
		t.Error("must fail, since description has invalid char")
	}
}

func TestSafeData_SetLongName(t *testing.T) {
	sh := NewSafeData(Config{App: "Test"})

	secret := genRandomSecret(14)
	secret.Identifier = genRandomData[string](_Cred_Max_Target_Length + 1)
	secret.Description = ""

	if err := sh.Set(secret); err != ErrMetadataMaxLength {
		t.Error(err)
	}
}
