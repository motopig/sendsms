package smservice

import (
	"net/http"

	"smservice/tool"

	"strings"

	"sync"

	"github.com/gin-gonic/gin"
)

var (
	failed  = 1
	success = 0
	smap    = make(map[string]bool)
	smsmut  sync.Mutex
)

func Router() {

	gin.SetMode(gin.DebugMode) //全局设置环境，此为开发环境，线上环境为gin.ReleaseMode
	router := gin.Default()    //获得路由实例

	//添加中间件
	router.Use(Middleware)

	//注册接口
	routerGroup := router.Group("/sms")
	{
		routerGroup.POST("/send", send)
		routerGroup.POST("/check", check)
		routerGroup.POST("/list", list)
	}

	//监听端口
	router.Run(":8006")
}

func Middleware(c *gin.Context) {

	// 判断是否有权限调用短信接口 sign 为 client 和 key 的 aes加密
	sign := c.PostForm("sign")
	mobile := c.PostForm("mobile")
	code := c.DefaultPostForm("code", "0")

	// 判断throttle 同一个手机号的并发
	smsmut.Lock()
	if _, ok := smap[mobile]; ok {
		smsmut.Unlock()
		c.JSON(http.StatusForbidden, gin.H{"msg": "busy ing!", "code": failed})
	}
	smap[mobile] = true
	smsmut.Unlock()

	// send返回后取消限制
	defer func() {
		smsmut.Lock()
		delete(smap, mobile)
		smsmut.Unlock()
	}()

	aesEnc := tool.NewEnc()
	source, err := aesEnc.Decrypt([]byte(sign))
	if err != nil || sign == "" {
		c.JSON(http.StatusForbidden, gin.H{"msg": "error sign!", "code": failed})
		c.Abort()
	}

	if len(source) != 0 {
		// todo 判断
	}

	// 判断手机号是否合法
	if err := VailMobile(mobile); err != nil {
		c.JSON(http.StatusConflict, gin.H{"msg": "error mobile!", "code": failed})
		c.Abort()
	}

	// 判断code是否合法
	if code != "0" {
		if err := VailCode(code); err != nil {
			c.JSON(http.StatusConflict, gin.H{"msg": "error code!", "code": failed})
			c.Abort()
		}
	}
}

func send(c *gin.Context) {

	// 加载对象
	sms := NewSms()
	sms.Mobile = c.PostForm("mobile")
	sms.ServiceName = c.PostForm("servicename")
	sms.Content = c.PostForm("content")
	sms.BankID = c.DefaultPostForm("bankid", "")
	sms.Fuid = c.DefaultPostForm("fuid", "0")
	sms.UserID = c.DefaultPostForm("userid", "0")
	sms.SendWay = c.PostForm("sendway")
	sms.Type = c.PostForm("type")

	if (sms.Content == "" && !strings.Contains("register,restpwd,getpwd", sms.ServiceName)) || sms.UserID == "0" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "短信内容有误!", "code": failed})
	} else {
		err := sms.send()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"msg": "send sms error!", "code": failed})
		} else {
			c.JSON(http.StatusOK, gin.H{"msg": "send sms success!", "code": success})
		}
	}
}

func check(c *gin.Context) {
	// 加载对象
	sms := NewSms()
	sms.Mobile = c.PostForm("mobile")
	sms.Content = c.PostForm("content")

	code := sms.check()

	if code == "" || len(code) != 4 {
		c.JSON(http.StatusOK, gin.H{"msg": "code check error!", "code": failed})
	} else {
		c.JSON(http.StatusOK, gin.H{"msg": "code check success!", "code": success})
	}
}

func list(c *gin.Context) {
	Select()
}
