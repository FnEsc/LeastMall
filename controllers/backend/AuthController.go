package backend

import (
	"LeastMall/models"
	"strconv"
)

type AuthController struct {
	BaseController
}

func (c *AuthController) Get() {
	var authList []models.Auth
	models.DB.Preload("AuthItem").Where("module_id=0").Find(&authList)
	c.Data["authList"] = authList
	c.TplName = "backend/auth/index.html"
}

func (c *AuthController) Add() {
	var authList []models.Auth
	models.DB.Where("module_id=0").Find(&authList)
	c.Data["authList"] = authList
	c.TplName = "backend/auth/add.html"
}

func (c *AuthController) GoAdd() {
	moduleName := c.GetString("module_name")
	iType, err1 := c.GetInt("type")
	actionName := c.GetString("actionName")
	url := c.GetString("url")
	moduleId, err2 := c.GetInt("module_id")
	sort, err3 := c.GetInt("sort")
	description := c.GetString("description")
	status, err4 := c.GetInt("status")
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		c.Error("GoAdd Auth 传入参数不合法", "/auth/add")
		return
	}
	auth := models.Auth{
		ModuleName:  moduleName,
		Type:        iType,
		ActionName:  actionName,
		Url:         url,
		ModuleId:    moduleId,
		Sort:        sort,
		Description: description,
		Status:      status,
	}
	err := models.DB.Create(&auth).Error
	if err != nil {
		c.Error("GoAdd auth 添加失败", "/auth/add")
		return
	} else {
		c.Success("GoAdd auth 增加数据成功", "/auth")
	}
}

func (c *AuthController) Edit() {
	id, err1 := c.GetInt("id")
	if err1 != nil {
		c.Error("Edit auth 传入参数错误", "/auth")
		return
	}
	auth := models.Auth{Id: id}
	models.DB.Find(&auth)
	var authList []models.Auth
	models.DB.Where("module_id=0").Find(&authList)
	c.Data["authList"] = authList
	c.TplName = "backend/auth/edit.html"
}

func (c *AuthController) GoEdit() {
	id, err1 := c.GetInt("id")
	iType, err2 := c.GetInt("type")
	sort, err3 := c.GetInt("sort")
	moduleId, err4 := c.GetInt("module_id")
	status, err5 := c.GetInt("status")
	moduleName := c.GetString("moduleName")
	actionName := c.GetString("actionName")
	description := c.GetString("description")
	url := c.GetString("url")
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		c.Error("GoEdit Auth 传惨错误", "/auth")
		return
	}
	auth := models.Auth{Id: id}
	models.DB.Find(&auth)
	auth.ModuleName = moduleName
	auth.Type = iType
	auth.Sort = sort
	auth.ModuleId = moduleId
	auth.Status = status
	auth.ModuleName = moduleName
	auth.ActionName = actionName
	auth.Description = description
	auth.Url = url
	err6 := models.DB.Save(&auth).Error
	if err6 != nil {
		c.Error("GoEdit auth修改权限失败", "/auth/edit?id="+strconv.Itoa(id))
		return
	}
	c.Success("GoEdit auth 修改权限成功", "/auth")
}

func (c *AuthController) Delete() {
	id, err := c.GetInt("id")
	if err != nil {
		c.Error("Delete auth id传惨错误", "/role")
		return
	}
	auth := models.Auth{Id: id}
	models.DB.Find(&auth)
	if auth.ModuleId == 0 {
		var auth2 []models.Auth
		models.DB.Where("module_id=?", auth.Id).Find(&auth2)
		if len(auth2) > 0 {
			c.Error("该菜单下有其他子目录，请先删除子层", "/auth")
			return
		}
	}
	models.DB.Delete(&auth)
	c.Success("auth delete success", "/auth")
}
