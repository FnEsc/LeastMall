package frontend

import (
	"LeastMall/models"
	"strconv"
)

// CartController 购物车结构体
type CartController struct {
	BaseController
}

func (c *CartController) Get() {
	c.BaseInit()

	var cartList []models.Cart
	models.Cookie.Get(c.Ctx, "cartList", &cartList)

	var allPrice float64
	// 执行计算总价
	for i := 0; i < len(cartList); i++ {
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
		}
	}
	c.Data["cartList"] = cartList
	c.Data["allPrice"] = allPrice
	c.TplName = "frontend/cart/cart.html"
}

func (c *CartController) AddCart() {
	c.BaseInit()

	colorId, err1 := c.GetInt("colorId")
	productId, err2 := c.GetInt("productId")

	var product models.Product
	var productColor models.ProductColor
	err3 := models.DB.Where("id=?", productId).Find(&product).Error
	err4 := models.DB.Where("id=?", colorId).Find(&productColor).Error

	// 判断参数错误
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		c.Ctx.Redirect(302, "item_"+strconv.Itoa(product.Id)+".html")
		return
	}

	// 获取增加购物车的商品数据
	currentData := models.Cart{
		Id:             product.Id,
		Title:          product.Title,
		Price:          product.Price,
		ProductVersion: product.ProductVersion,
		Num:            1,
		ProductColor:   productColor.ColorName,
		ProductImg:     product.ProductImg,
		ProductGift:    product.ProductGift,
		ProductAttr:    "",   // 根据需求拓展，作为备注字段
		Checked:        true, // 加入购物车时，默认打勾
	}
	// 判断cookie购物车有没有数据 (购物车存在于cookie中)
	var cartList []models.Cart
	models.Cookie.Get(c.Ctx, "cartList", &cartList)
	if len(cartList) > 0 { // cookie购物车有数据
		if models.CartHasData(cartList, currentData) { // 购物车有当前商品：id+color+attr, 则往购物车该商品，数量+1
			for i := 0; i < len(cartList); i++ {
				if cartList[i].Id == currentData.Id && cartList[i].ProductColor == currentData.ProductColor && cartList[i].ProductAttr == currentData.ProductAttr { // 判断为同一商品
					cartList[i].Num += 1
				}
			}
		} else { // 购物车有商品，但无当前商品，则添加
			cartList = append(cartList, currentData)
		}
		models.Cookie.Set(c.Ctx, "cartList", cartList)
	} else { // 购物车无任何商品，则直接吧当前商品数据写入cookie购物车
		cartList = append(cartList, currentData)
		models.Cookie.Set(c.Ctx, "cartList", cartList)
	}

	c.Data["product"] = product
	c.TplName = "frontend/cart/add_cart_success.html"
}

// DecCart 购物车商品数量-1
func (c *CartController) DecCart() {
	var flag bool
	var allPrice float64
	var currentAllPrice float64
	var num int

	productId, _ := c.GetInt("productId")
	productColor := c.GetString("productColor")
	productAttr := ""

	var cartList []models.Cart
	models.Cookie.Get(c.Ctx, "cartList", &cartList)
	for i := 0; i < len(cartList); i++ {
		if cartList[i].Id == productId && cartList[i].ProductColor == productColor && cartList[i].ProductAttr == productAttr { // 判断为同一商品
			if cartList[i].Num > 1 { // 准备扣减的数量需要大于1
				cartList[i].Num -= 1
			}
			flag = true // 标识是否修改成功
			num = cartList[i].Num
			currentAllPrice = cartList[i].Price * float64(cartList[i].Num) // 计算修改数量商品的总价
		}
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num) // 重新计算购物车的总价
		}
	}
	if flag {
		models.Cookie.Set(c.Ctx, "cartList", cartList)
		c.Data["json"] = map[string]interface{}{
			"success":      true,
			"msg":          "修改数量成功",
			"allPrice":     allPrice,
			"currentPrice": currentAllPrice,
			"num":          num, // 当前购物车该商品的数量
		}
	} else {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "传入扣减数量商品的参数错误",
		}
	}
	c.ServeJSON()
}

// IncCart 购物车商品数量+1
func (c *CartController) IncCart() {
	var flag bool
	var allPrice float64
	var currentAllPrice float64
	var num int

	productId, _ := c.GetInt("productId")
	productColor := c.GetString("productColor")
	productAttr := ""

	var cartList []models.Cart
	models.Cookie.Get(c.Ctx, "cartList", &cartList)
	for i := 0; i < len(cartList); i++ {
		if cartList[i].Id == productId && cartList[i].ProductColor == productColor && cartList[i].ProductAttr == productAttr { // 判断为同一商品
			cartList[i].Num += 1
			flag = true // 标识是否修改成功
			num = cartList[i].Num
			currentAllPrice = cartList[i].Price * float64(cartList[i].Num) // 计算修改数量商品的总价
		}
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num) // 重新计算购物车的总价
		}
	}

	if flag {
		models.Cookie.Set(c.Ctx, "cartList", cartList)
		c.Data["json"] = map[string]interface{}{
			"success":      true,
			"msg":          "修改数量成功",
			"allPrice":     allPrice,
			"currentPrice": currentAllPrice,
			"num":          num, // 当前购物车该商品的数量
		}
	} else {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "传入增加数量商品的参数错误",
		}
	}
	c.ServeJSON()
}

// ChangeOneCart 反向勾选购物车单一商品
func (c *CartController) ChangeOneCart() {
	var flag bool
	var allPrice float64

	productId, _ := c.GetInt("productId")
	productColor := c.GetString("productColor")
	productAttr := ""

	var cartList []models.Cart
	models.Cookie.Get(c.Ctx, "cartList", &cartList)
	for i := 0; i < len(cartList); i++ {
		if cartList[i].Id == productId && cartList[i].ProductColor == productColor && cartList[i].ProductAttr == productAttr { // 判断为同一商品
			cartList[i].Checked = !cartList[i].Checked
			flag = true // 标识是否修改成功
		}
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num) // 重新计算购物车的总价
		} else {
			allPrice -= cartList[i].Price * float64(cartList[i].Num) // 重新计算购物车的总价
		}
	}

	if flag {
		models.Cookie.Set(c.Ctx, "cartList", cartList)
		c.Data["json"] = map[string]interface{}{
			"success":  true,
			"msg":      "修改购物车商品状态成功",
			"allPrice": allPrice,
		}
	} else {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "传入修改购物车商品状态的参数错误",
		}
	}
	c.ServeJSON()
}

// ChangeAllCart 购物车全部商品状态反选
func (c *CartController) ChangeAllCart() {
	flag, _ := c.GetInt("flag")
	var allPrice float64
	var cartList []models.Cart
	models.Cookie.Get(c.Ctx, "cartList", &cartList)
	for i := 0; i < len(cartList); i++ {
		if flag == 1 { // 全选状态至1
			cartList[i].Checked = true
		} else {
			cartList[i].Checked = false
		}
		// 计算总价
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
		}
	}
	models.Cookie.Set(c.Ctx, "cartList", cartList)

	c.Data["json"] = map[string]interface{}{
		"success":  true,
		"msg":      "全部反选成功",
		"allPrice": allPrice,
	}
	c.ServeJSON()
}

// DelCart 删除购物车的某一件商品
func (c *CartController) DelCart() {
	productId, _ := c.GetInt("productId")
	productColor := c.GetString("productColor")
	productAttr := ""

	var cartList []models.Cart
	models.Cookie.Get(c.Ctx, "cartList", &cartList)
	for i := 0; i < len(cartList); i++ {
		if cartList[i].Id == productId && cartList[i].ProductColor == productColor && cartList[i].ProductAttr == productAttr { // 判断为同一商品
			cartList = append(cartList[:i], cartList[i+1:]...) // 执行删除该商品，做赋值操作
		}
	}
	models.Cookie.Set(c.Ctx, "cartList", cartList)
	c.Redirect("/cart", 302)
}
