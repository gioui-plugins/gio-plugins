//go:build !android && !darwin && !ios

package inapppay

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
}

func (d *driver) requestToken() error {
	return ErrNotAvailable
}
