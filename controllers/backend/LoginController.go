package backend

import (
	"LeastMall/common"
	"LeastMall/models"
	"strings"
)

type LoginController struct {
	BaseController
}

func (c *LoginController) Get() {
	c.TplName = "backend/login/login.html"
}

func (c *LoginController) GoLogin() {
	var flag = models.Cpt.VerifyReq(c.Ctx.Request)
	if flag {
		username := strings.Trim(c.GetString("username"), "")
		password := common.Md5(strings.Trim(c.GetString("password"), ""))
		var administratorList []models.Administrator
		models.DB.Where("username=? AND password=? AND status=1", username, password).Find(&administratorList)
		if len(administratorList) == 1{
			c.SetSession("userinfo", administratorList[0])
			c.Success("Gologin admin success", "/")
		}else {
			c.Error("GoLogin admin failed 账号密码错误", "/login")
		}
	} else {
		c.Error("GoLogin admin failed 验证码错误", "/login")
	}
}

func (c *LoginController) LogOut()  {
	c.DelSession("userinfo")
	c.Success("Go Logout admin success", "/login")
}
