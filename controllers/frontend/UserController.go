package frontend

import (
	"LeastMall/models"
	"math"
	"strconv"
	"time"
)

type UserController struct {
	BaseController
}

func (c *UserController) Get() {
	c.BaseInit()

	var user models.User
	models.Cookie.Get(c.Ctx, "userinfo", &user)
	c.Data["user"] = user
	timeHour := time.Now().Hour()
	if timeHour >= 0 && timeHour <= 6 {
		c.Data["helo"] = "深夜了，注意休息哦！"
	} else if timeHour >= 7 && timeHour <= 11 {
		c.Data["helo"] = "尊敬的用户上午好"
	} else {
		c.Data["helo"] = "尊敬的用户下午好"
	}
	var orders []models.Order
	models.DB.Where("uid=?", user.Id).Find(&orders)
	var wait_pay int
	var wait_rec int
	for i := 0; i < len(orders); i++ {
		// 支付状态 0-未支付 1-已支付
		if orders[i].PayStatus == 0 {
			wait_pay += 1
		}
		// 订单状态 0-已下单 1-已付款 2-已配货 3-发货 4-交易成功 5-退货 6-取消
		if orders[i].OrderStatus >= 2 && orders[i].OrderStatus < 4 {
			wait_rec += 1
		}
	}
	c.Data["wait_pay"] = wait_pay
	c.Data["wait_rec"] = wait_rec
	c.TplName = "frontend/user/welcome.html"
}

func (c *UserController) OrderList() {
	c.BaseInit()
	// 获取当前用户
	var user models.User
	models.Cookie.Get(c.Ctx, "userinfo", &user)

	// 获取当前用户下面的订单信息 并作分页
	page, _ := c.GetInt("page")
	if page == 0 {
		page = 1
	}
	pageSize := 2

	// 获取搜索关键词
	where := "uid=?"
	keywords := c.GetString("keywords")
	if keywords != "" {
		var orderItems []models.OrderItem
		models.DB.Where("product_title like ?", "%"+keywords+"%").Find(&orderItems)
		var str string
		for i := 0; i < len(orderItems); i++ {
			if i == 0 {
				str += strconv.Itoa(orderItems[i].OrderId)
			} else {
				str += "," + strconv.Itoa(orderItems[i].OrderId)
			}
		}
		where += " AND id in ( " + str + " )"
	} // 结束关键词 where str 拼接
	// 获取筛选条件
	orderStatus, err := c.GetInt("order_status")
	if err == nil {
		where += "AND order_status=" + strconv.Itoa(orderStatus)
		c.Data["orderStatus"] = orderStatus
	} else {
		c.Data["orderStatus"] = "nil"
	}
	// 总数量
	var count int
	models.DB.Where(where, user.Id).Table("order").Count(&count)
	var orders []models.Order
	models.DB.Where(where, user.Id).Offset((page - 1) * pageSize).
		Limit(pageSize).
		Preload("OrderItem").
		Order("add_time desc").
		Find(&orders)

	c.Data["orders"] = orders
	c.Data["totalPages"] = math.Ceil(float64(count) / float64(pageSize))
	c.Data["page"] = page
	c.Data["keywords"] = keywords
	c.TplName = "frontend/user/order.html"
}

func (c *UserController) OrderInfo() {
	c.BaseInit()
	id, _ := c.GetInt("id")
	var user models.User
	models.Cookie.Get(c.Ctx, "userinfo", &user)
	var order models.Order
	models.DB.Where("id=? AND uid=?", id, user.Id).Preload("OrderItem").Find(&order)
	c.Data["order"] = order
	if order.OrderId == "" {
		c.Redirect("/", 302)
	}
	c.TplName = "frontend/user/order_info.html"
}
