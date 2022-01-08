package models

type ProductCollect struct {
	Id        int
	UserId    int
	ProductId int
	AddTime   int
}

func (ProductCollect) TableName() string {
	return "product_collect"
}
