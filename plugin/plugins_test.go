package plugin

import (
	"testing"
)

func TestRegister(t *testing.T) {
	registeredPlugins = nil
	examplePlugin := NewHandlerFunc(nil, nil, nil, nil)
	Register(examplePlugin)

	if len(registeredPlugins) != 1 {
		t.Errorf("expected 1 registered plugin, got %d", len(registeredPlugins))
	}
}
