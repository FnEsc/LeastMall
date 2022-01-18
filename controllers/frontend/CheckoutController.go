package frontend

import (
	"LeastMall/common"
	"LeastMall/models"
	"fmt"
	"strconv"
)

type CheckoutController struct {
	BaseController
}

func (c CheckoutController) Checkout() {
	c.BaseInit()

	// 1. 获取要结算的商品
	var cartList []models.Cart  // 购物车的商品
	var orderList []models.Cart // 要结算的商品
	models.Cookie.Get(c.Ctx, "cartList", &cartList)

	// 2. 计算总价
	var allPrice float64
	for i := 0; i < len(cartList); i++ {
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
			orderList = append(orderList, cartList[i])
		}
	}

	// 3. 判断结算商品
	if len(orderList) == 0 {
		c.Redirect("/", 302)
		return
	}
	c.Data["orderList"] = orderList
	c.Data["allPrice"] = allPrice

	// 4. 获取收货地址
	var user models.User
	models.Cookie.Get(c.Ctx, "userinfo", &user)
	var addressList models.Address
	models.DB.Where("uid=?", user.Id).Order("default_address desc").Find(&addressList)
	c.Data["addressList"] = addressList

	// 5. 生成签名 防止重复提交订单
	orderSign := common.Md5(common.GetRandomNum())
	c.SetSession("orderSign", orderSign)
	c.Data["orderSign"] = orderSign

	c.TplName = "frontend/buy/checkout.html"
}

// GoOrder 提交订单
func (c *CheckoutController) GoOrder() {
	// 0. 防止重复提交订单
	orderSign := c.GetString("orderSign")
	sessionOrderSign := c.GetSession("orderSign")
	if orderSign != sessionOrderSign {
		c.Redirect("/", 302)
		return
	}
	c.DelSession("orderSign")

	// 1. 获取收货地址信息
	var user models.User
	models.Cookie.Get(c.Ctx, "userinfo", &user)

	var addressList []models.Address
	models.DB.Where("uid=? and default_address=1", user.Id).Find(&addressList)

	if len(addressList) > 0 {
		// 2. 获取购物的商品信息 => orderList
		var cartList []models.Cart
		var orderList []models.Cart // 要结算的商品
		models.Cookie.Get(c.Ctx, "cartList", &cartList)
		var allPrice float64
		for i := 0; i < len(cartList); i++ {
			if cartList[i].Checked {
				allPrice += cartList[i].Price * float64(cartList[i].Num)
				orderList = append(orderList, cartList[i])
			}
		}

		// 3. 把订单信息放在订单表，把商品信息放在商品表
		order := models.Order{
			OrderId:     common.GenerateOrderId(),
			Uid:         user.Id,
			AllPrice:    allPrice,
			Phone:       addressList[0].Phone,
			Name:        addressList[0].Name,
			Address:     addressList[0].Address,
			Zipcode:     addressList[0].Zipcode,
			PayStatus:   0,
			PayType:     0,
			OrderStatus: 0,
			AddTime:     int(common.GetUnix()),
		}
		err := models.DB.Create(&order).Error
		if err != nil {
			for i := 0; i < len(orderList); i++ {
				orderItem := models.OrderItem{
					OrderId:        order.Id,
					Uid:            user.Id,
					ProductTile:    orderList[i].Title,
					ProductId:      orderList[i].Id,
					ProductImg:     orderList[i].ProductImg,
					ProductPrice:   orderList[i].Price,
					ProductNum:     orderList[i].Num,
					ProductVersion: orderList[i].ProductVersion,
					ProductColor:   orderList[i].ProductColor,
					AddTime:        int(common.GetUnix()),
				}
				err := models.DB.Create(&orderItem).Error
				if err != nil {
					fmt.Println(err)
				}
			}

			// 4. 购物车中删除已结算商品
			var noSelectCartList []models.Cart
			for i := 0; i < len(cartList); i++ {
				if !cartList[i].Checked {
					noSelectCartList = append(noSelectCartList, cartList[i])
				}
			}
			models.Cookie.Set(c.Ctx, "cartList", noSelectCartList)
			c.Redirect("/buy/comfirm?id="+strconv.Itoa(order.Id), 302)
		} else {
			// 创建主订单失败
			c.Redirect("/", 302)
		}
	} else {
		// 无主地址
		c.Redirect("/", 302)
	}
}

// Confirm 确认结算订单
func (c *CheckoutController) Confirm() {
	c.BaseInit()
	id, err := c.GetInt("id")
	if err != nil {
		c.Redirect("/", 302)
		return
	}
	// 获取用户信息
	var user models.User
	models.Cookie.Get(c.Ctx, "userinfo", &user)

	// 获取主订单信息
	var order models.Order
	models.DB.Where("id=?", id).Find(&order)
	c.Data["order"] = order
	// 判断订单是为当前用户
	if user.Id != order.Uid {
		c.Redirect("/", 302)
		return
	}

	// 获取主订单下的商品信息
	var orderItem []models.OrderItem
	models.DB.Where("order_id=?", order.Id).Find(&orderItem)
	c.Data["orderItem"] = orderItem

	c.TplName = "frontend/buy/confirm.html"
}

// OrderPayStat 获取订单支付状态
func (c CheckoutController) orderPayStatus() {
	// 获取订单号
	id, err := c.GetInt("id")
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "未传入订单id",
		}
		c.ServeJSON()
		return
	}

	// 查询订单
	var order models.Order
	models.DB.Where("id=?", id).Find(&order)

	// 判断当前用户权限
	var user models.User
	models.Cookie.Get(c.Ctx, "userinfo", &user)
	if user.Id != order.Uid {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "订单用户权限校验失败",
		}
		c.ServeJSON()
		return
	}

	// 判断订单，支付状态=已支付 && 订单状态=已付款
	if order.PayStatus == 1 && order.OrderStatus == 1 {
		c.Data["json"] = map[string]interface{}{
			"success": true,
			"msg":     "订单已支付",
		}
		c.ServeJSON()
		return
	} else {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "订单未支付",
		}
		c.ServeJSON()
		return
	}

}
