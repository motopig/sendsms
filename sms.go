package smservice

import (
	"errors"
	"time"

	"fmt"
	"strings"

	log "github.com/golang/glog"
	uuid "github.com/satori/go.uuid"
)

type msg struct {
	message string
	code    int
}

type SMS struct {
	Mobile      string
	Content     string
	ServiceName string
	Config      *ServiceConfig
	NowTime     time.Time
	BankID      string
	UserID      string
	Fuid        string
	SendWay     string
	Type        string
	ReadStatus  int
	AppType     int
}

type Sender interface {
	Send(sms *SMS) error
}

func NewSms() *SMS {
	sms := &SMS{}
	sms.NowTime = time.Now()
	return sms
}

// send sms
func (sms *SMS) send() error {
	var s Sender
	// 根据serviceName切换通道和配置
	agent := config.ServiceList[sms.ServiceName].Agent
	terr := sms.renderTpl()

	if terr != nil {
		return terr
	}
	switch agent {
	case "blue":
		s = &Blue{sms: sms}
	//case "xxx":
	//	s = &xxx{sms: sms}
	default:
		log.Warning("短信服务商设置错误")
		return errors.New("短信通道设置不正确")
	}

	err := s.Send(sms)
	if err != nil {
		return err
	}

	// 如果是发送验证码 需要保存
	if sms.ServiceName == "register" || sms.ServiceName == "restpwd" || sms.ServiceName == "getpwd" {
		// store the code
		status := storeCode(sms.Mobile, sms.Content)
		if status == false {
			return errors.New("发送失败")
		}
	} else {
		// 保存短信数据到数据库
		var um UserMessage
		um.Id = uuid.NewV4().String()
		um.Content = sms.Content
		um.UserID = sms.UserID
		um.ReadStatus = 2
		um.FromUid = sms.Fuid
		um.Type = sms.Type
		um.AppType = 1
		um.Title = sms.ServiceName
		um.SendWay = sms.SendWay
		um.BankId = sms.BankID
		um.CreatedAt = time.Now().Unix()
		um.UpdatedAt = time.Now().Unix()

		row := Add(um)

		if row != 1 {
			return errors.New("发送失败")
		}
	}

	return nil
}

func (sms *SMS) check() string {
	code := getCode(sms.Mobile)
	if code == "none" {
		return ""
	}
	return code
}

func (sms *SMS) renderTpl() error {

	if strings.Contains("register,restpwd,getpwd", sms.ServiceName) {
		sms.Content = GenCode()
	} else {
		// 处理短信模板
		tpl := config.ServiceList[sms.ServiceName].Tpl
		sp := strings.Split(sms.Content, `,`)
		lensp := len(sp)
		switch lensp {
		case 1:
			sms.Content = fmt.Sprintf(tpl, sp[0])
		case 2:
			sms.Content = fmt.Sprintf(tpl, sp[0], sp[1])
		case 3:
			sms.Content = fmt.Sprintf(tpl, sp[0], sp[1], sp[2])
		case 4:
			sms.Content = fmt.Sprintf(tpl, sp[0], sp[1], sp[2], sp[3])
		case 5:
			sms.Content = fmt.Sprintf(tpl, sp[0], sp[1], sp[2], sp[3], sp[4])
		default:
			sms.Content = ""
		}

	}

	if sms.Content == "" {
		return errors.New(sms.ServiceName + ":短信模板获取失败")
	} else {
		return nil
	}

}
