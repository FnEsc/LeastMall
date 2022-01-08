package models

type RoleAuth struct {
	AuthId int
	RoleId int
}

func (RoleAuth) TableName() string {
	return "role_auth"
}
