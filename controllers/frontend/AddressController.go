package frontend

import "LeastMall/models"

type AddressController struct {
	BaseController
}

func (c *AddressController) AddAddress() {
	var user models.User
	models.Cookie.Get(c.Ctx, "userinfo", &user)
	name := c.GetString("name")
	phone := c.GetString("phone")
	address := c.GetString("address")
	zipcode := c.GetString("zipcode")
	var addressCount int
	models.DB.Where("uid=?", user.Id).Table("address").Count(&addressCount)
	if addressCount > 10 {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "添加收货地址失败，收货地址数量超过限制",
		}
		c.ServeJSON()
		return
	}
	models.DB.Table("address").Where("uid=?", user.Id).Updates(map[string]interface{}{"default_address": 0})
	addressResult := models.Address{
		Uid:            user.Id,
		Name:           name,
		Phone:          phone,
		Address:        address,
		Zipcode:        zipcode,
		DefaultAddress: 1,
	}
	models.DB.Create(&addressResult)
	var allAddressResult []models.Address
	models.DB.Where("uid=?", user.Id).Find(&allAddressResult)
	c.Data["json"] = map[string]interface{}{
		"success": false,
		"msg":     "添加收货地址成功",
		"result":  allAddressResult,
	}
	c.ServeJSON()
}

func (c *AddressController) GetOneAddressInfo() {
	addressId, err := c.GetInt("addressId")
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "查询地址失败，参数错误",
		}
		c.ServeJSON()
		return
	}
	var address models.Address
	models.DB.Where("id=?", addressId).Find(&address)
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"msg":     "查询地址成功",
		"result":  address,
	}
	c.ServeJSON()
}

func (c *AddressController) GoEditAddressInfo() {
	var user models.User
	models.Cookie.Get(c.Ctx, "userinfo", &user)
	addressId, err := c.GetInt("addressId")
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "修改地址失败，参数错误",
		}
		c.ServeJSON()
		return
	}
	name := c.GetString("name")
	phone := c.GetString("phone")
	address := c.GetString("address")
	zipcode := c.GetString("zipcode")
	models.DB.Table("address").Where("uid=?", user.Id).Updates(map[string]interface{}{"default_address": 0})
	var addressModel models.Address
	models.DB.Where("id=?", addressId).Find(&addressModel)
	addressModel.Name = name
	addressModel.Phone = phone
	addressModel.Address = address
	addressModel.Zipcode = zipcode
	addressModel.DefaultAddress = 1
	models.DB.Save(&addressModel)

	// 修改完成，查询当前用户的所有收货地址返回
	var allAddressResult []models.Address
	models.DB.Where("uid=?", user.Id).Find(&allAddressResult)
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"msg":     "修改收货地址成功",
		"result":  allAddressResult,
	}
	c.ServeJSON()
}

func (c *AddressController) ChangeDefaultAddress() {
	var user models.User
	models.Cookie.Get(c.Ctx, "userinfo", &user)
	addressId, err := c.GetInt("addressId")
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "修改默认地址失败，参数错误",
		}
		c.ServeJSON()
		return
	}
	models.DB.Where("uid=?", user.Id).Updates(map[string]interface{}{"default_address": 0})
	models.DB.Where("id=?", addressId).Updates(map[string]interface{}{"default_address": 1})
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"msg":     "修改默认地址成功",
	}
	c.ServeJSON()

}
