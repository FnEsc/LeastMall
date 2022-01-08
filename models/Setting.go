package models

import "reflect"

type Setting struct {
	Id              int    `form:"id,omitempty"`
	SiteTitle       string `form:"site_title"`
	SiteLogo        string `form:"site_logo"`
	SiteDescription string `form:"site_description"`
	NoPicture       string `form:"no_picture"`
	SiteIcp         string `form:"site_icp"`
	SearchKeywords  string `form:"search_keywords"`
	TongjiCode      string `form:"tongji_code"`
	Appid           string `form:"appid"`
	AppSecret       string `form:"app_secret"`
	EndPoint        string `form:"end_point"`
	BucketName      string `form:"bucket_name"`
	OssStatus       int    `form:"oss_status"`
}

func (Setting) TableName() string {
	return "setting"
}

func GetSettingByColumn(columnName string) string {
	setting := Setting{}
	DB.First(&setting)
	// 反射获取column
	v := reflect.ValueOf(setting)
	val := v.FieldByName(columnName).String()
	return val
}
