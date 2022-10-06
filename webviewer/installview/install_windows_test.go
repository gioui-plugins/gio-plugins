package installview

import (
	"context"
	"testing"
)

func TestInstaller_Install(t *testing.T) {
	installer, err := Download(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}

	if err = installer.Install(); err != nil {
		t.Fatal(err)
	}
}
