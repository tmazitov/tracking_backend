package bl

import (
	"database/sql"
	"time"
)

type UserRole int

const (
	Base    UserRole = 0
	Worker  UserRole = 1
	Manager UserRole = 2
	Admin   UserRole = 3
)

type R_GetUser struct {
	ID        int64    `json:"id,omitempty"`
	ShortName string   `json:"shortName"`
	RoleID    UserRole `json:"roleId"`
}

type DB_GetUser struct {
	ID        sql.NullInt64
	ShortName sql.NullString
	RoleID    sql.NullInt32
}

type PutUser struct {
	TelNumber string `json:"telNumber"`
	ShortName string `json:"shotName"`
}

type UserHoliday struct {
	WorkerId int64
	AuthorId int64
	Date     *time.Time
}

type UserJob struct {
	JobType       uint8  `json:"jobType" binding:"required"`
	JobExperience uint8  `json:"jobExperience" binding:"gte=0"`
	JobMail       string `json:"jobMail,omitempty"`
}
