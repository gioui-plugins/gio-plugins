package auth

import (
	"github.com/gioui-plugins/gio-plugins/auth/providers"
	"github.com/gioui-plugins/gio-plugins/hyperlink"
	"net/url"
)

type driver struct {
	hp *hyperlink.Hyperlink
}

func attachDriver(house *Auth, config Config) {
	house.driver = driver{
		hp: hyperlink.NewHyperlink(hyperlink.Config{}),
	}
	configureDriver(&house.driver, config)
}

func configureDriver(driver *driver, config Config) {}

func (d *driver) open(provider providers.Provider, nonce string) error {
	u, err := url.Parse(provider.URL(nonce))
	if err != nil {
		return err
	}
	return d.hp.Open(u)
}
