package models

type ProductImage struct {
	Id        int
	ProductId int
	ImgUrl    string
	ColorId   int
	Sort      int
	AddTime   int
	Status    int
}

func (ProductImage) TableName() string {
	return "product_image"
}
