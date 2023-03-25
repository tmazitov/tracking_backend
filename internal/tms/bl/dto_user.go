package bl

import "database/sql"

type UserRole int

const (
	Base    UserRole = 1
	Worker  UserRole = 2
	Manager UserRole = 3
	Admin   UserRole = 4
)

type GetUser struct {
	TelNumber string   `json:"telNumber"`
	ShortName string   `json:"shotName"`
	RoleID    UserRole `json:"roleId"`
}

type SelectUserById struct {
	TelNumber sql.NullString
	ShortName sql.NullString
	RoleID    UserRole
}
