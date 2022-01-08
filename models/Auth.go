package models

type Auth struct {
	Id          int
	ModuleName  string // 模块名称
	ActionName  string // 操作名称
	Type        int    // 节点类型 1-模块 2-菜单 3-操作
	Url         string // 路由跳转地址
	ModuleId    int    // 此module_id和当前模型的id关联 module_id=0 标识模块
	Sort        int
	Description string
	Status      int
	AddTime     int
	AuthItem    []Auth `gorm:"foreignkey:ModuleId;association_foreignkey:Id"`
	Checked     bool   `gorm:"-"`
}

func (Auth) TableName() string {
	return "auth"
}
