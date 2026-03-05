//go:build !android && !darwin && !ios

package inapppay

type driver struct {
	config   Config
	inapppay *InAppPay
}

func attachDriver(push *InAppPay, config Config) {
	d := &driver{inapppay: push}
	push.driver = d
	configureDriver(d, config)
}

func configureDriver(d *driver, config Config) {
	d.config = config
}

func (d *driver) listProducts(_ []string) error {
	return ErrNotAvailable
}

func (d *driver) purchase(_ string, _ string, _ bool) error {
	return ErrNotAvailable
}
