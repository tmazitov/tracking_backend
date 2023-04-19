package bl

import "database/sql"

type UserRole int

const (
	Base    UserRole = 0
	Worker  UserRole = 1
	Manager UserRole = 2
	Admin   UserRole = 3
)

type R_GetUser struct {
	ID        int64    `json:"id,omitempty"`
	TelNumber string   `json:"telNumber"`
	ShortName string   `json:"shortName"`
	RoleID    UserRole `json:"roleId"`
}

type DB_GetUser struct {
	ID        int64
	TelNumber sql.NullString
	ShortName sql.NullString
	RoleID    UserRole
}

type PutUser struct {
	TelNumber string `json:"telNumber"`
	ShortName string `json:"shotName"`
}
