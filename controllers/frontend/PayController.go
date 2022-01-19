package frontend

import (
	"LeastMall/models"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/objcoding/wxpay"
	"github.com/skip2/go-qrcode"
	"github.com/smartwalle/alipay/v3"
	"strconv"
	"strings"
	"time"
)

type PayController struct {
	BaseController
}

func (c *PayController) Alipay() {
	AliId, err1 := c.GetInt("id")
	if err1 != nil {
		c.Redirect(c.Ctx.Request.Referer(), 302)
	}
	var orderItem []models.OrderItem
	models.DB.Where("order_id=?", AliId).Find(&orderItem)
	var privateKey = "xxxxxx" // 必须，上一部使用 RSA 签名验签工具 生成的私钥
	var client, err = alipay.New("202201181915", privateKey, true)
	client.LoadAppPublicCertFromFile("certfile/appCertPublicKey_202201181915.certfile") // 加载应用公钥证书
	client.LoadAliPayRootCertFromFile("certfile/alipayRootCert.certfile")               // 加载支付宝根证书
	client.LoadAliPayPublicCertFromFile("certfile/alipayCertPublicKey_RSA2.certfile")   // 加载支付宝公钥证书

	// 将key的验证调整到初始化阶段
	if err != nil {
		fmt.Println(err)
		return
	}

	// 计算总价格
	var totalAmount float64
	for i := 0; i < len(orderItem); i++ {
		totalAmount += orderItem[i].ProductPrice * float64(orderItem[i].ProductNum)
	}
	var p alipay.TradePagePay
	p.NotifyURL = "xxxxxx"
	p.ReturnURL = "xxxxxx"
	p.TotalAmount = "0.01"
	p.Subject = "订单order——" + time.Now().Format("200601021504")
	p.OutTradeNo = "WF" + time.Now().Format("200601021504")
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	var url, errAlipayPage = client.TradePagePay(p)
	if errAlipayPage != nil {
		fmt.Println(errAlipayPage)
	}
	var payURL = url.String()
	c.Redirect(payURL, 302)
}

// AlipayNotify 支付宝收款，订单状态置为已支付
func (c *PayController) AlipayNotify() {
	var privateKey = "xxxxxx" // 必须，上一部使用 RSA 签名验签工具 生成的私钥
	var client, err = alipay.New("202201181915", privateKey, true)
	client.LoadAppPublicCertFromFile("certfile/appCertPublicKey_202201181915.certfile") // 加载应用公钥证书
	client.LoadAliPayRootCertFromFile("certfile/alipayRootCert.certfile")               // 加载支付宝根证书
	client.LoadAliPayPublicCertFromFile("certfile/alipayCertPublicKey_RSA2.certfile")   // 加载支付宝公钥证书

	if err != nil {
		fmt.Println(err)
		return
	}

	req := c.Ctx.Request
	req.ParseForm()
	ok, err := client.VerifySign(req.Form)
	if !ok || err != nil {
		c.Redirect(c.Ctx.Request.Referer(), 302)
	}
	rep := c.Ctx.ResponseWriter
	var noti, _ = client.GetTradeNotification(req)
	if noti != nil {
		fmt.Println("交易成功：", noti.TradeStatus)
		if string(noti.TradeStatus) == "TRADE_SUCCESS" {
			var order models.Order
			temp := strings.Split(noti.OutTradeNo, "_")[1]
			id, _ := strconv.Atoi(temp)
			models.DB.Where("id=?", id).Find(&order)
			order.PayType = 0
			order.PayStatus = 1
			order.OrderStatus = 1
			models.DB.Save(&order)
		}
	}
	alipay.AckNotification(rep) // 确认收到通知消息
}

// AlipayReturn 支付宝返回
func (c *PayController) AlipayReturn() {
	c.Redirect("/user/order", 302)
}

func (c *PayController) WxPay() {
	WxId, err := c.GetInt("id")
	if err != nil {
		c.Redirect(c.Ctx.Request.Referer(), 302)
	}
	var orderItem []models.OrderItem
	models.DB.Where("order_id=?", WxId).Find(&orderItem)

	// 配置基本信息
	account := wxpay.NewAccount(
		"xxxxxx", // appid
		"xxxxxx", // 商户号
		"xxxxxx", // appKey
		false,
	)
	client := wxpay.NewClient(account)
	var price int64
	for i := 0; i < len(orderItem); i++ {
		price = 1
	}
	// 获取ip地址，订单号等信息
	ip := strings.Split(c.Ctx.Request.RemoteAddr, ":")[0]
	tradeNo := time.Now().Format("200601021504")

	// 调用统一下单
	params := make(wxpay.Params)
	params.SetString("body", "order——"+time.Now().Format("200601021504")).
		SetString("out_trade_no", tradeNo+"_"+strconv.Itoa(WxId)).
		SetInt64("total_fee", price).
		SetString("spbill_create_ip", ip).
		SetString("notify_url", "http://xxxxxx/wxpay/notify"). // 配置回调地址
		SetString("trade_type", "APP").                        // APP端支付
		SetString("trade_type", "NATIVE")                      // 网站支付需要设置要NATIVE

	p, err1 := client.UnifiedOrder(params)
	logs.Info(p)
	if err1 != nil {
		logs.Error(err1)
		c.Redirect(c.Ctx.Request.Referer(), 302)
	}

	// 获取 code_url 生成支付二维码
	var pngObj []byte
	pngObj, _ = qrcode.Encode(p["code_url"], qrcode.Medium, 256)
	c.Ctx.WriteString(string(pngObj))

}

func (c *PayController) WxPayNotify() {
	// 获取表单传过来的xml数据
	xmlStr := string(c.Ctx.Input.RequestBody)
	postParams := wxpay.XmlToMap(xmlStr)
	logs.Info(postParams)

	// 校验签名
	account := wxpay.NewAccount(
		"xxxxxx",
		"xxxxxx",
		"xxxxxx",
		false,
	)
	client := wxpay.NewClient(account)
	isValidate := client.ValidSign(postParams)

	// xml解析
	params := wxpay.XmlToMap(xmlStr)
	if isValidate == true {
		if params["return_code"] == "SUCCESS" {
			idStr := strings.Split(params["out_trade_no"], "_")[1]
			id, _ := strconv.Atoi(idStr)
			var order models.Order
			models.DB.Where("id=?", id).Find(&order)
			order.PayType = 0
			order.PayStatus = 1
			order.OrderStatus = 1
			models.DB.Save(&order)
		} else {
			c.Redirect(c.Ctx.Request.Referer(), 302)
		}
	}
}
