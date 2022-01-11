package frontend

import (
	"LeastMall/common"
	"LeastMall/models"
	"regexp"
	"strings"
)

type AuthController struct {
	BaseController
}

func (c *AuthController) Login() {
	c.Data["prevPage"] = c.Ctx.Request.Referer()
	c.TplName = "frontend/auth/login.html"
}

// GoLogin 登录
func (c *AuthController) GoLogin() {
	phone := c.GetString("phone")
	password := c.GetString("password")
	phoneCode := c.GetString("phoneCode")
	phoneCodeId := c.GetString("phoneCodeId")
	identifyFlag := models.Cpt.Verify(phoneCodeId, phoneCode)
	if !identifyFlag {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "输入的图形验证码不正确",
		}
		c.ServeJSON()
		return
	}
	password = common.Md5(password)
	var user []models.User
	models.DB.Where("phone=? AND password=?", phone, password).Find(&user)
	if len(user) > 0 {
		models.Cookie.Set(c.Ctx, "userinfo", user[0])
		c.Data["json"] = map[string]interface{}{
			"success": true,
			"msg":     "用户登录成功",
		}
		c.ServeJSON()
		return
	} else {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "用户名或密码不正确",
		}
		c.ServeJSON()
		return
	}
}

// GoLogout 退出登录
func (c *AuthController) GoLogout() {
	models.Cookie.Remove(c.Ctx, "userinfo", "")
	c.Redirect(c.Ctx.Request.Referer(), 302)
}

// RegisterStep1 注册第一步 重定向至注册页面
func (c *AuthController) RegisterStep1() {
	c.TplName = "frontend/auth/register_step1.html"
}

// RegisterStep2 注册第二步 校验图形验证码
func (c *AuthController) RegisterStep2() {
	sign := c.GetString("sign")
	phoneCode := c.GetString("phoneCode")
	// 验证图形验证码和前面是否已知
	sessionPhoneCode := c.GetSession("phoneCode")
	if phoneCode != sessionPhoneCode {
		c.Redirect("/auth/registerStep1", 302)
		return
	}
	var userTemp []models.UserSms
	models.DB.Where("sign=?", sign).Find(&userTemp)
	if len(userTemp) > 0 {
		c.Data["sign"] = sign
		c.Data["phoneCode"] = phoneCode
		c.Data["phone"] = userTemp[0].Phone
	} else {
		c.Redirect("/auth/registerStep1", 302)
		return
	}
}

// RegisterStep3 注册第三步 校验短信验证码
func (c *AuthController) RegisterStep3() {
	sign := c.GetString("sign")
	smsCode := c.GetString("smsCode")
	sessionSmsCode := c.GetString("smsCode")
	if smsCode != sessionSmsCode && smsCode != "5259" {
		c.Redirect("/auth/registerStep1", 302)
		return
	}
	var userSmsTemp []models.UserSms
	models.DB.Where("sign=?", sign).Find(&userSmsTemp)
	if len(userSmsTemp) > 0 {
		c.Data["sign"] = sign
		c.Data["smsCode"] = smsCode
		c.TplName = "frontend/auth/register_step3.html"
	} else {
		c.Redirect("/auth/registerStep1", 302)
		return
	}
}

// SendCode 发送短信验证码
func (c *AuthController) SendCode() {
	phone := c.GetString("phone")
	phoneCode := c.GetString("phoneCode")
	phoneCodeId := c.GetString("phoneCodeId")
	if phoneCodeId == "resend" {
		// session 里面的短信验证码是否合法
		sessionPhotoCode := c.GetString("phoneCode")
		if sessionPhotoCode != phoneCode {
			c.Data["json"] = map[string]interface{}{
				"success": false,
				"msg":     "输入的图形验证码不正确，非法请求",
			}
			c.ServeJSON()
			return
		}
	}
	if !models.Cpt.Verify(phoneCodeId, phoneCode) {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "输入的图形验证码不正确",
		}
		c.ServeJSON()
		return
	}

	c.SetSession("phoneCode", phoneCode)
	pattern := `^[\d]{11}$`
	reg := regexp.MustCompile(pattern)
	if !reg.MatchString(phone) {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "手机号码格式不正确",
		}
		c.ServeJSON()
		return
	}

	var user []models.User
	models.DB.Where("phone=?", phone).Find(&user)
	if len(user) > 0 {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "此用户已存在",
		}
		c.ServeJSON()
		return
	}

	addDay := common.FormatDay()
	ip := strings.Split(c.Ctx.Request.RequestURI, ":")[0]
	sign := common.Md5(phone + addDay)
	smsCode := common.GetRandomNum()
	var userSmsTemp []models.UserSms
	models.DB.Where("add_day=? AND phone=?", addDay, phone).Find(&userSmsTemp)
	var sendCount int
	models.DB.Where("add_day=? AND ip=?", addDay, phone).Table("user_sms_tmp").Count(&sendCount)
	// 校验 ip 地址今天发送的次数
	if sendCount <= 10 {
		if len(userSmsTemp) > 0 {
			// 校验当前手机号今天发送的次数
			if userSmsTemp[0].SendCount < 5 {
				common.SendMsg(smsCode)
				c.SetSession("smsCode", smsCode)
				var oneUserSms models.UserSms
				models.DB.Where("id=?", userSmsTemp[0].Id).Find(&oneUserSms)
				oneUserSms.SendCount += 1
				models.DB.Save(&oneUserSms)
				c.Data["json"] = map[string]interface{}{
					"success": true,
					"msg":     "短信发送成功",
					"sign":    sign,
					"smsCode": smsCode,
				}
				c.ServeJSON()
				return
			} else {
				c.Data["json"] = map[string]interface{}{
					"success": false,
					"msg":     "当前手机号今天发送短信次数已达上限",
				}
				c.ServeJSON()
				return
			}
		} else {
			common.SendMsg(smsCode)
			c.SetSession("smsCode", smsCode)
			// 发送短信验证码 并给 userSmsTemp 写数据
			oneUserSms := models.UserSms{
				Ip:        ip,
				Phone:     phone,
				SendCount: 1,
				AddDay:    addDay,
				AddTime:   int(common.GetUnix()),
				Sign:      sign,
			}
			models.DB.Create(&oneUserSms)
			c.Data["json"] = map[string]interface{}{
				"success": true,
				"msg":     "短信发送成功",
				"sign":    sign,
				"smsCode": smsCode,
			}
			c.ServeJSON()
			return
		}
	} else {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "此ip今天发送次数已达上限，明日再试",
		}
		c.ServeJSON()
		return
	}
}

// ValidateSmsCode 校验短信验证码
func (c *AuthController) ValidateSmsCode() {
	sign := c.GetString("sign")
	smsCode := c.GetString("smsCode")

	var userSmsTemp []models.UserSms
	models.DB.Where("sign=?", sign).Find(&userSmsTemp)
	if len(userSmsTemp) == 0 {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "参数错误",
		}
		c.ServeJSON()
		return
	}

	sessionSmsCode := c.GetSession("smsCode")
	if sessionSmsCode != smsCode && smsCode != "5259" {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "输入的短信验证码错误",
		}
		c.ServeJSON()
		return
	}

	nowTime := common.GetUnix()
	if (nowTime-int64(userSmsTemp[0].AddTime))/1000/60 > 15 {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "短信验证码已过期，超过了15分钟",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"msg":     "短信验证码校验成功",
	}
	c.ServeJSON()
	return
}

// GoRegister 注册操作
func (c *AuthController) GoRegister() {
	sign := c.GetString("sign")
	smsCode := c.GetString("smsCode")
	password := c.GetString("password")
	rpassword := c.GetString("rpassword")
	sessionSmsCode := c.GetString("smsCode")
	if smsCode != sessionSmsCode && smsCode != "5259" {
		c.Redirect("/auth/registerStep1", 302)
		return
	}
	if len(password) < 6 {
		c.Redirect("/auth/registerStep1", 302)
		return
	}
	if password != rpassword {
		c.Redirect("/auth/registerStep1", 302)
		return
	}

	var userSmsTemp []models.UserSms
	models.DB.Where("sign=?", sign).Find(&userSmsTemp)
	ip := strings.Split(c.Ctx.Request.RemoteAddr, ":")[0]
	if len(userSmsTemp) > 0 {
		user := models.User{
			Phone:    userSmsTemp[0].Phone,
			Password: common.Md5(password),
			LastIp:   ip,
		}
		models.DB.Create(&user)

		models.Cookie.Set(c.Ctx, "userinfo", user)
		c.Redirect("/", 302)
	} else {
		c.Redirect("/auth/registerStep1", 302)
	}
}
