package models

type ProductCate struct {
	Id          int
	Title       string
	CateImg     string
	Link        string
	Template    string
	Pid         int
	SubTitle    string
	Keywords    string
	Description string
	Sort        int
	Status      int
	AddTime     int
	ProductCate []ProductCate `gorm:"foreignkey:Pid;association_foreignkey:Id"`
}

func (ProductCate) TableName() string {
	return "product_cate"
}
