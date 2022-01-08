package models

type User struct {
	Id       int
	phone    string
	Password string
	AddTime  int
	LastIp   string
	Email    string
	Status   int
}

func (User) TableName() string {
	return "user"
}
