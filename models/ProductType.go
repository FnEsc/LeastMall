package models

type ProductType struct {
	Id          int
	Title       string
	Description string
	Status      int
	AddTime     int
}

func (ProductType) TableName() string {
	return "product_type"
}
