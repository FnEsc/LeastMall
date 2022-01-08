package models

type Menu struct {
	Id          int
	Title       string
	Link        string
	Position    string
	IsOpennew   int
	Relation    string
	Sort        int
	Status      int
	AddTime     int
	ProductItem []Product `gorm:"-"`
}

func (Menu) TableNmae() string {
	return "menu"
}
