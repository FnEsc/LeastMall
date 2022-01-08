package models

type Order struct {
	Id          int
	OrderId     string
	Uid         int
	AllPrice    float64
	Phone       string
	Name        string
	Address     string
	Zipcode     string
	PayStatus   int // 支付状态 0-未支付 1-已支付
	PayType     int // 支付渠道 0-alipay 1-wechat
	OrderStatus int // 订单状态 0-已下单 1-已付款 2-已配货 3-发货 4-交易成功 5-退货 6-取消
	AddTime     int
	OrderItem   []OrderItem `gorm:"foreignkey:OrderId;association_foreignkey:Id"`
}

func (Order) TableName() string {
	return "order"
}
