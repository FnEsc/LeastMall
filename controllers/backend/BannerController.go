package backend

import (
	"LeastMall/common"
	"LeastMall/models"
	"github.com/astaxie/beego/logs"
	"os"
	"strconv"
)

type BannerController struct {
	BaseController
}

func (c *BannerController) Get() {
	var bannerList []models.Banner
	models.DB.Find(&bannerList)
	c.Data["bannerList"] = bannerList
	c.TplName = "backend/banner/index.html"
}

func (c *BannerController) Add() {
	c.TplName = "backend/banner/add.html"
}

func (c *BannerController) GoAdd() {
	bannerType, err1 := c.GetInt("bannerType")
	sort, err2 := c.GetInt("sort")
	status, err3 := c.GetInt("status")
	title := c.GetString("title")
	link := c.GetString("link")
	if err1 != nil || err2 != nil || err3 != nil {
		c.Error("GoAdd banner 传参不合法", "/banner/add")
		return
	}

	bannerImgSrc, err4 := c.UploadImg("bannerImg")
	if err4 == nil {
		banner := models.Banner{
			Title:      title,
			BannerType: bannerType,
			BannerImg:  bannerImgSrc,
			Link:       link,
			Sort:       sort,
			Status:     status,
			AddTime:    int(common.GetUnix()),
		}
		models.DB.Create(&banner)
		c.Success("GoAdd banner增加轮播图成功", "/banner")
	} else {
		c.Error("GoAdd banner 增加轮播图失败", "/banner/add")
	}
}

func (c *BannerController) Edit() {
	id, err := c.GetInt("id")
	if err != nil {
		c.Error("Edit banner 非法传参", "/banner")
		return
	}
	banner := models.Banner{Id: id}
	models.DB.Find(&banner)
	c.Data["banner"] = banner
	c.TplName = "backend/banner/edit.html"
}

func (c *BannerController) GoEdit() {
	id, err1 := c.GetInt("id")
	bannerType, err2 := c.GetInt("bannerType")
	sort, err3 := c.GetInt("sort")
	status, err4 := c.GetInt("status")
	link := c.GetString("link")
	title := c.GetString("title")
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		c.Error("GoEdit banner 传参不合法", "/banner")
		return
	}
	bannerImgSrc, _ := c.UploadImg("bannerImg")
	banner := models.Banner{Id: id}
	models.DB.Find(&banner)
	banner.Title = title
	banner.BannerType = bannerType
	banner.Link = link
	banner.Sort = sort
	banner.Status = status
	if bannerImgSrc != "" {
		banner.BannerImg = bannerImgSrc
	}
	err5 := models.DB.Save(&banner).Error
	if err5 != nil {
		c.Error("GoEdit banner 修改轮播图失败", "/banner/edit?id="+strconv.Itoa(id))
		return
	}
	c.Success("GoEdit banner 修改轮播图成功", "/banner")
}

func (c *BannerController) Delete() {
	id, err := c.GetInt("id")
	if err != nil {
		c.Error("Delete banner 传参错误", "/banner")
	}
	banner := models.Banner{Id: id}
	models.DB.Find(&banner)
	address := "baseDir" + banner.BannerImg
	deleteImg := os.Remove(address)
	models.DB.Delete(&banner) // 先删除文件，再做返回
	if deleteImg != nil {
		logs.Error(deleteImg)
		c.Error("Delete banner 轮播图成功，文件未清除", "/banner")
	} else {
		c.Success("Delete banner 轮播图成功，文件已清除", "/banner")
	}
}
