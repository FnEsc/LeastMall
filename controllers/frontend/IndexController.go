package frontend

import (
	"LeastMall/models"
	"fmt"
	"time"
)

type IndexController struct {
	BaseController
}

func (c *IndexController) Get() {
	// 调用功能
	c.BaseInit()

	// 开始时间
	startTime := time.Now().UnixNano()

	// 获取轮播图，注意获取的时候要写地址
	var bannerList []models.Banner
	if hasBannerList := models.CacheDb.Get("banners", &bannerList); hasBannerList == true {
		c.Data["bannerList"] = bannerList
	} else {
		models.DB.Where("status=1 AND banner_type=1").Order("sort desc").Find(&bannerList)
		c.Data["bannerList"] = bannerList
		models.CacheDb.Set("bannerList", bannerList)
	}

	// 获取手机商品列表
	var productPhoneList []models.Product
	if hasProductPhoneList := models.CacheDb.Get("productPhoneList", &productPhoneList); hasProductPhoneList == true {
		c.Data["productPhoneList"] = productPhoneList
	} else {
		productPhoneList := models.GetProductByCategory(1, "hot", 8)
		c.Data["productPhoneList"] = productPhoneList
		models.CacheDb.Set("productPhoneList", productPhoneList)
	}

	// 获取电视商品列表
	var productTvList []models.Product
	if hasProductTvList := models.CacheDb.Get("productTvList", &productTvList); hasProductTvList == true {
		c.Data["productTvList"] = productTvList
	} else {
		productTvList := models.GetProductByCategory(4, "best", 8)
		c.Data["productTvList"] = productTvList
		models.CacheDb.Set("productTvList", productTvList)
	}

	// 结束时间
	endTime := time.Now().UnixNano()

	fmt.Println("执行时间", endTime-startTime)

	c.TplName = "frontend/index/index.html"
}
