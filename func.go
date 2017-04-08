package smservice

import (
	"errors"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

func VailMobile(mobile string) error {

	if len(mobile) < 11 {
		return errors.New("[mobile]参数不对")
	}
	reg, err := regexp.Compile("^1[3-8][0-9]{9}$")
	if err != nil {
		panic("regexp error")
	}
	if !reg.MatchString(mobile) {
		return errors.New("手机号码[mobile]格式不正确")
	}
	return nil
}

func VailCode(code string) error {

	if len(code) != 4 {
		return errors.New("[code]参数不对")
	}
	c, err := regexp.Compile("^[0-9]{4}$")
	if err != nil {
		panic("regexp error")
	}
	if !c.MatchString(code) {
		return errors.New("验证码[code]格式不正确")
	}
	return nil
}

func GenCode() string {
	return strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(8999) + 1000)
}
