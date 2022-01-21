package backend

import (
	"LeastMall/common"
	"LeastMall/models"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/astaxie/beego"
	"os"
	"path"
	"strconv"
	"strings"
)

type BaseController struct {
	beego.Controller
}

func (c *BaseController) Success(msg string, redirect string) {
	c.Data["msg"] = msg
	if strings.Contains(redirect, "http") {
		c.Data["redirect"] = redirect
	} else {
		c.Data["redirect"] = "/" + beego.AppConfig.String("adminPath") + redirect
	}
	c.TplName = "backend/public/success.html"
}

func (c *BaseController) Error(msg string, redirect string) {
	c.Data["msg"] = msg
	if strings.Contains(redirect, "http") {
		c.Data["redirect"] = redirect
	} else {
		c.Data["redirect"] = "/" + beego.AppConfig.String("adminPath") + redirect
	}
	c.TplName = "backend/public/error.html"
}

func (c *BaseController) Goto(redirect string) {
	c.Redirect("/"+beego.AppConfig.String("adminPath")+redirect, 302)
}

func (c *BaseController) UploadImg(picName string) (string, error) {
	ossStatus, _ := beego.AppConfig.Bool("ossStatus")
	if ossStatus == true {
		return c.OssUploadImg(picName)
	}
	return c.LocalUploadImg(picName)
}

func (c BaseController) LocalUploadImg(picName string) (string, error) {
	f, h, err := c.GetFile(picName)
	if err != nil {
		return "LocalUploadImg 获取文件失败", err
	}
	defer func() {
		_ = f.Close() // 关闭文件流
	}()

	// 检查文件后缀
	extName := path.Ext(h.Filename)
	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".gif":  true,
		".jpeg": true,
	}
	if _, ok := allowExtMap[extName]; !ok {
		return "LocalUploadImg 图片后缀名不合法", err
	}

	// 创建图保存目录 static/upload/20210121
	day := common.FormatDay()
	dir := "static/upload/" + day
	if err := os.MkdirAll(dir, 0666); err != nil {
		return "LocalUploadImg 创建文件目录失败", err
	}

	// 生成文件名称
	fileUnixName := strconv.FormatInt(common.GetUnixNano(), 10)
	saveDir := path.Join(dir, fileUnixName+extName)

	// 保存图片
	err = c.SaveToFile(picName, saveDir)
	return saveDir, err
}

func (c BaseController) OssUploadImg(picName string) (string, error) {
	var setting models.Setting
	models.DB.First(&setting)
	f, h, err := c.GetFile(picName)
	if err != nil {
		return "LocalUploadImg 获取文件失败", err
	}

	defer func() {
		_ = f.Close() // 关闭文件流
	}()

	// 检查文件后缀
	extName := path.Ext(h.Filename)
	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".gif":  true,
		".jpeg": true,
	}
	if _, ok := allowExtMap[extName]; !ok {
		return "OssUploadImg 图片后缀名不合法", err
	}

	// 创建OSS实例
	client, err := oss.New(setting.EndPoint, setting.Appid, setting.AppSecret)
	if err != nil {
		return "OssUploadImg 创建OSS实例失败", err
	}

	// 获取存储空间
	bucket, err := client.Bucket(setting.BucketName)
	if err != nil {
		return "OssUploadImg 获取存储空间失败", err
	}

	// 创建图保存目录 static/upload/20210121
	day := common.FormatDay()
	dir := "static/upload/" + day
	fileUnixName := strconv.FormatInt(common.GetUnixNano(), 10)
	saveDir := path.Join(dir, fileUnixName+extName)

	// 上传文件流
	err = bucket.PutObject(saveDir, f)
	if err != nil {
		return "OssUploadImg 上传文件流失败", err
	}
	return saveDir, nil
}

func (c BaseController) GetSetting() models.Setting {
	setting := models.Setting{Id: 1}
	models.DB.First(&setting)
	return setting
}
