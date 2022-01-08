package common

import (
	"LeastMall/models"
	"github.com/astaxie/beego/context"
)

func FrontendAuth(ctx *context.Context) {
	// 判断前台用户有没有登录
	user := models.User{}
	models.Cookie.Get(ctx, "userinfo", &user)
	if len(user.Phone) != 11 {
		ctx.Redirect(302, "/auth/login")
	}
}
