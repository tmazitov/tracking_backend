package bl

import "database/sql"

type UserRole int

const (
	Base    UserRole = 0
	Worker  UserRole = 1
	Manager UserRole = 2
	Admin   UserRole = 3
)

type GetUser struct {
	TelNumber string   `json:"telNumber"`
	ShortName string   `json:"shotName"`
	RoleID    UserRole `json:"roleId"`
}

type PutUser struct {
	TelNumber string `json:"telNumber"`
	ShortName string `json:"shotName"`
}

type SelectUserById struct {
	TelNumber sql.NullString
	ShortName sql.NullString
	RoleID    UserRole
}
