//go:build windows

package pushnotification

import (
	"github.com/gioui-plugins/gio-plugins/pushnotification/internal"
	"github.com/go-ole/go-ole"
)

type driver struct {
	config Config
	push   *Push
}

func attachDriver(push *Push, config Config) {
	d := &driver{push: push}
	push.driver = d
	configureDriver(d, config)
}

func configureDriver(d *driver, config Config) {
	d.config = config
	if id := ole.NewGUID(d.config.WindowsAzureConfig.ObjectID); id != nil {
		d.config.WindowsAzureConfig.ObjectID = id.String()
	}
}

func (d *driver) requestToken() (Token, error) {
	token, err := internal.GetChannelURI(d.config.WindowsAzureConfig.ObjectID)
	if err != nil {
		return Token{}, err
	}
	return Token{Token: token}, nil
}
