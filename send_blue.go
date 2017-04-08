package smservice

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// todo config file
const (
	SMSURL      = "http://222.73.117.158/msg/HttpBatchSendSM"
	SMSACCOUNT  = "xxx"
	SMSPASSWORD = "yyy"
	MERGEURL    = SMSURL + "?account=" + SMSACCOUNT + "&pswd=" + SMSPASSWORD + "&mobile=%s&msg=%s"
)

type Blue struct {
	sms *SMS
}

func (b *Blue) Send(sms *SMS) error {

	if b.Call(sms) == nil {
		return errors.New("发送失败")
	}
	return nil
}

func (b *Blue) Call(sms *SMS) []byte {

	strUrl := fmt.Sprintf(MERGEURL, sms.Mobile, sms.Content)
	r, err := http.NewRequest("GET", strUrl, nil)
	if err != nil {
		fmt.Println("http.NewRequest: ", err.Error())
		return nil
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		fmt.Println("http.DefaultClient.Do: ", err.Error())
		return nil
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		// 超过运营商条数限制
		fmt.Println("resp.StatusCode!=http.StatusOK: ", resp.StatusCode)
		return nil
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil && err != io.EOF {
		fmt.Println("ioutil.ReadAll: ", err.Error())
		return nil
	}

	ret := strings.Split(string(data), ",")

	if len(ret[1]) > int(1) {
		return nil
	}

	return data
}
