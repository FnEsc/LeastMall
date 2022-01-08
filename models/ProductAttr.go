package models

type ProductAttr struct {
	Id              int
	Product         int
	AttributeCateId int
	AttributeId     int
	AttributeTitle  string
	AttributeType   int
	AttributeValue  string
	Sort            int
	AddTime         int
	Status          int
}

func (ProductAttr) TableName() string {
	return "product_attr"
}
