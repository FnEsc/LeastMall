package models

type Administrator struct {
	Id       int
	Username string
	Passowrd string
	Mobile   string
	Email    string
	Status   int
	RoleId   int `gorm:"roleid"`
	AddTime  int
	IsSuper  int
	Role     Role `gorm:"foreignkey:Id;association_foreignkey:RoleId"`
}

func (Administrator) TableName() string {
	return "administrator"
}
