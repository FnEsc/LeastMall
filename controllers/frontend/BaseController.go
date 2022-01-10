package frontend

import (
	"LeastMall/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	"net/url"
	"strings"
)

type BaseController struct {
	beego.Controller
}

func (c *BaseController) BaseInit() {

	// 获取顶部导航
	var topMenu models.Menu
	if hasTopMenu := models.CacheDb.Get("topMenu", &topMenu); hasTopMenu == true {
		c.Data["topMenuList"] = topMenu
	} else {
		models.DB.Where("status=1 AND position=1").Order("sort desc").Find(&topMenu)
		c.Data["topMenuList"] = topMenu
		models.CacheDb.Set("topMenu", topMenu)
	}

	// 左侧分类（预加载）
	var productCateList []models.ProductCate
	if hasProductCate := models.CacheDb.Get("productCateList", &productCateList); hasProductCate == true {
		c.Data["productCateList"] = productCateList
	} else {
		models.DB.Preload(
			"ProductCateItem",
			func(db *gorm.DB) *gorm.DB {
				return db.Where("product_cate.status=1").Order("product_cate.sort Desc")
			}).Where("pid=0 and status=1").Order("sort desc", true).Find(&productCateList)
		c.Data["productCateList"] = productCateList
		models.CacheDb.Set("productCateList", productCateList)
	}

	// 获取中间导航的数据
	var middleMenuList []models.Menu
	if hasMiddleMenu := models.CacheDb.Get("middleMenuList", &middleMenuList); hasMiddleMenu == true {
		c.Data["hasMiddleMenu"] = middleMenuList
	} else {
		models.DB.Where("status=1 AND position=2").Order("sort desc").Find(&middleMenuList)
		for i := 0; i < len(middleMenuList); i++ {
			// 获取关联商品
			middleMenuList[i].Relation = strings.ReplaceAll(middleMenuList[i].Relation, "，", ",")
			relation := strings.Split(middleMenuList[i].Relation, ",")
			var products []models.Product
			models.DB.Where("id in (?)", relation).Limit(6).Order("sort asc").
				Select("id,title,product_img,price").Find(&products)
			middleMenuList[i].ProductItem = products
		}
		c.Data["middleMenuList"] = middleMenuList
		models.CacheDb.Set("middleMenuList", middleMenuList)
	}

	// 判断用户是否登录
	var user models.User
	models.Cookie.Get(c.Ctx, "userinfo", &user)
	if len(user.Phone) == 11 {
		str := fmt.Sprintf(`
			<ul>
				<li class="userinfo">
					<a href="#">%v</a>
					<i class="i"></i>
					<ol>
						<li><a href="/user">个人中心</a></li>
						<li><a href="#">我的收藏</a></li>
						<li><a href="/auth/loginOut">退出登录</a></li>
					</ol>
				</li>
			</ul>`, user.Phone)
		c.Data["userInfo"] = str
	} else {
		str := fmt.Sprintf(`
			<ul>
				<li><a href="/auth/login" target="_blank">登录</a></li>
				<li>|</li>
				<li><a href="/auth/registerStep1" target="_blank" >注册</a></li>
			</ul>`)
		c.Data["userinfo"] = str
	}
	urlPath, _ := url.Parse(c.Ctx.Request.URL.String())
	c.Data["pathName"] = urlPath.Path

}
