package backend

import (
	"LeastMall/common"
	"LeastMall/models"
	"strings"
)

type AdministratorController struct {
	BaseController
}

func (c *AdministratorController) Get() {
	var administratorList []models.Administrator
	models.DB.Preload("Role").Find(&administratorList)
	c.Data["administratorList"] = administratorList
	c.TplName = "backend/administrator/index.html"
}

func (c *AdministratorController) Add() {
	var roleList []models.Role
	models.DB.Find(&roleList)
	c.Data["roleList"] = roleList
	c.TplName = "backend/administrator/add.html"
}

func (c *AdministratorController) GoAdd() {
	username := strings.Trim(c.GetString("username"), "")
	password := strings.Trim(c.GetString("password"), "")
	mobile := strings.Trim(c.GetString("mobile"), "")
	email := strings.Trim(c.GetString("email"), "")
	roleId, err1 := c.GetInt("roleId")
	if err1 != nil {
		c.Error("GoAdd 传入roleId参数不合法", "/administrator/add")
		return
	}
	if len(username) < 2 || len(password) < 6 {
		c.Error("GoAdd 传入账号密码参数不合法", "/administrator/add")
		return
	} else if common.VerifyEmail(email) == false {
		c.Error("GoAdd 传入邮箱参数不合法", "/administrator/add")
		return
	}
	var administratorList []models.Administrator
	models.DB.Where("username=?", username).Find(&administratorList)
	if len(administratorList) > 0 {
		c.Error("GoAdd 传入用户名已存在", "/administrator/add")
		return
	}

	var administrator models.Administrator
	administrator.Username = username
	administrator.Passowrd = password
	administrator.Mobile = mobile
	administrator.Email = email
	administrator.Status = 1
	administrator.AddTime = int(common.GetUnix())
	administrator.RoleId = roleId
	err := models.DB.Create(&administrator).Error
	if err != nil {
		c.Error("GoAdd 创建数据库记录失败", "administrator/add")
		return
	}
	c.Success("GoAdd 添加管理员成功", "administrator")
}

func (c *AdministratorController) Edit() {
	id, err := c.GetInt("id")
	if err != nil {
		c.Error("Edit传入id参数错误", "/administrator")
		return
	}
	administrator := models.Administrator{Id: id}
	models.DB.Find(&administrator)
	c.Data["administrator"] = administrator
	var roleList []models.Role
	models.DB.Find(&roleList)
	c.Data["roleList"] = roleList
	c.TplName = "backend/administrator/edit.html"
}

func (c *AdministratorController) GoEdit() {

}
