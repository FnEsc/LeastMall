package models

type Product struct {
	Id              string
	Title           string
	SubTitle        string
	ProductSn       string
	CateId          int
	ClickCount      int
	ProductNumber   int
	Price           float64
	MarketPrice     float64
	RelationProduct string
	ProductAttr     string
	ProductVersion  string
	ProductImg      string
	ProductGift     string
	ProductColor    string
	ProductKeywords string
	ProductDesc     string
	ProductContent  string
	IsDelete        int
	IsHot           int
	IsBest          int
	IsNew           int
	ProductTypeId   int
	Sort            int
	Status          int
	AddTime         int
}

func (Product) TableName() string {
	return "product"
}

func GetProductByCategory(cateId int, productType string, limitNum int) []Product {
	var productList []Product
	var productCateList []ProductCate
	DB.Where("pid=?", cateId).Find(&productCateList)
	var templice []int
	if len(productCateList) > 0 {
		for i := 0; i < len(productCateList); i++ {
			templice = append(templice, productCateList[i].Id)
		}
	}
	templice = append(templice, cateId)
	where := "cate_id in (?)"
	switch productType {
	case "hot":
		where += "AND is_hot=1"
	case "best":
		where += "AND is_best=1"
	case "new":
		where += "AND is_new=1"
	default:
		break
	}
	DB.Where(where, templice).
		Select("id,title,price,product_img,sub_title").
		Limit(limitNum).
		Order("sort desc").
		Find(&productList)
	return productList
}
