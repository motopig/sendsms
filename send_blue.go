package smservice

import (
	"blue"
)

func (b *Blue) Send(sms *SMS) error {

	blue.ACCOUNT = config.BlueConf["account"]
	blue.PASSWORD = config.BlueConf["password"]

	err, status := blue.Send(sms.Mobile, sms.Content)
	if status > 0 {
		return err
	} else {
		return nil
	}
}
