package frontend

import (
	"LeastMall/common"
	"LeastMall/models"
	"math"
	"strconv"
	"strings"
)

type ProductController struct {
	BaseController
}

func (c *ProductController) CategoryList() {
	c.BaseInit()

	id := c.Ctx.Input.Param(":id")
	cateId, _ := strconv.Atoi(id)
	var currentProductCate models.ProductCate
	var subProductCate []models.ProductCate
	models.DB.Where("id=?", cateId).Find(&currentProductCate)

	// 当前页
	page, _ := c.GetInt("page")
	if page == 0 {
		page = 1
	}
	// 每一页显示的数量
	pageSize := 5

	var tempSlice []int
	if currentProductCate.Pid == 0 { // 顶级分类
		models.DB.Where("pid=?", currentProductCate.Id).Find(&subProductCate) // 二级分类
		for i := 0; i < len(subProductCate); i++ {
			tempSlice = append(tempSlice, subProductCate[i].Id)
		}
	} else {
		// 获取当前二级分类对应的同级分类
		models.DB.Where("pid=?", currentProductCate.Pid).Find(&subProductCate)
	}

	tempSlice = append(tempSlice, cateId)
	where := "cate_id in (?)"
	var productList []models.Product
	models.DB.Where(where, tempSlice).Select("id, title, price, product_img, sub_title").Offset((page - 1) * pageSize).Limit(pageSize).Find(&productList)
	// 查询product表里面的数量
	var count int
	models.DB.Where(where, tempSlice).Table("product").Count(&count)

	c.Data["productList"] = productList
	c.Data["subProductCate"] = subProductCate
	c.Data["currentProductCate"] = currentProductCate
	c.Data["totalPages"] = math.Ceil(float64(count) / float64(pageSize))
	c.Data["page"] = page

	// 制定分类模板
	tpl := currentProductCate.Template
	if tpl == "" {
		tpl = "frontend/product/list.html"
	}
	c.TplName = tpl
}

// ProductItem 获取产品细节
func (c *ProductController) ProductItem() {
	c.BaseInit()

	id := c.Ctx.Input.Param(":id")
	// 获取当前商品信息
	var product models.Product
	models.DB.Where("id=?", id).Find(&product)
	c.Data["product"] = product

	// 获取关联商品 RelationProduct
	var relationProduct []models.Product
	product.RelationProduct = strings.ReplaceAll(product.RelationProduct, "，", ",")
	relationIds := strings.Split(product.RelationProduct, ",")
	models.DB.Where("id in (?)", relationIds).Select("id, title, price, product_version").Find(&relationProduct)
	c.Data["relationProduct"] = relationProduct

	// 获取关联赠品 ProductGift
	var productGift []models.Product
	product.ProductGift = strings.ReplaceAll(product.ProductGift, "，", ",")
	giftIds := strings.Split(product.ProductGift, ",")
	models.DB.Where("id in (?)", giftIds).Select("id, title, price, product_version").Find(&productGift)
	c.Data["productGift"] = productGift

	// 获取关联颜色 ProductColor
	var productColor []models.Product
	product.ProductColor = strings.ReplaceAll(product.ProductColor, "，", ",")
	colorIds := strings.Split(product.ProductColor, ",")
	models.DB.Where("id in (?)", colorIds).Select("id, title, price, product_version").Find(&productColor)
	c.Data["productColor"] = productColor

	// 获取关联配件 ProductFitting
	var productFitting []models.Product
	product.ProductFitting = strings.ReplaceAll(product.ProductFitting, "，", ",")
	fittingIds := strings.Split(product.ProductFitting, ",")
	models.DB.Where("id in (?)", fittingIds).Select("id, title, price, product_version").Find(&productFitting)
	c.Data["productFitting"] = productFitting

	// 获取商品关联的照片 ProductImage
	var productImage []models.Product
	models.DB.Where("product_id=?", product.Id).Find(&productImage)
	c.Data["productImage"] = productImage

	// 获取商品参数信息 ProductAttr
	var productAttr []models.Product
	models.DB.Where("product_id=?", product.Id).Find(&productAttr)
	c.Data["productAttr"] = productAttr

	c.TplName = "frontend/product/item.html"
}

// Collect 收藏商品
func (c ProductController) Collect() {
	productId, err := c.GetInt("productId")
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "收藏商品传参错误",
		}
		c.ServeJSON()
		return
	}

	var user models.User
	ok := models.Cookie.Get(c.Ctx, "userinfo", &user)
	if ok != true {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "收藏商品失败，请先登录",
		}
		c.ServeJSON()
		return
	}

	isExist := models.DB.First(&user)
	if isExist.RowsAffected == 0 {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "收藏商品失败，非法用户",
		}
		c.ServeJSON()
		return
	}

	var goodCollect models.ProductCollect
	isExist = models.DB.Where("user_id=? AND product_id=?", user.Id, productId).First(&goodCollect)
	if isExist.RowsAffected == 0 {
		goodCollect.UserId = user.Id
		goodCollect.ProductId = productId
		goodCollect.AddTime = int(common.GetUnix())
		models.DB.Create(&goodCollect)
		c.Data["json"] = map[string]interface{}{
			"success": true,
			"msg":     "收藏成功",
		}
		c.ServeJSON()
	} else {
		models.DB.Delete(&goodCollect)
		c.Data["json"] = map[string]interface{}{
			"success": true,
			"msg":     "取消收藏成功",
		}
		c.ServeJSON()
	}
}

func (c *ProductController) GetImgList() {
	colorId, err1 := c.GetInt("colorId")
	productId, err2 := c.GetInt("productId")
	// 查询商品图库信息
	var productImage []models.ProductImage
	err3 := models.DB.Where("color_id=? AND product_id=?", colorId, productId).First(&productImage).Error

	if err1 != nil || err2 != nil || err3 != nil {
		c.Data["json"] = map[string]interface{}{
			"success": true,
			"msg":     "获取商品照片失败",
		}
		c.ServeJSON()
	} else {
		if len(productImage) == 0 {
			models.DB.Where("product_id=?", productId).Find(&productImage) // 制定颜色无照片，条件降低
		}
		c.Data["json"] = map[string]interface{}{
			"success": true,
			"result":  productImage,
		}
		c.ServeJSON()
	}

}
