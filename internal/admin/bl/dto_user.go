package bl

import "database/sql"

type UserRole int

const (
	Base    UserRole = 0
	Worker  UserRole = 1
	Manager UserRole = 2
	Admin   UserRole = 3
)

type R_User struct {
	ID        int64    `json:"id,omitempty"`
	ShortName string   `json:"shortName"`
	RoleID    UserRole `json:"roleId"`
}

type DB_User struct {
	ID        sql.NullInt64
	ShortName sql.NullString
	RoleID    sql.NullInt32
}

func (u *DB_User) ToReal() *R_User {
	return &R_User{
		ID:        u.ID.Int64,
		ShortName: u.ShortName.String,
		RoleID:    UserRole(u.RoleID.Int32),
	}
}

type DB_UserOffer struct {
	Id            int
	UserId        int64
	JobType       uint8
	JobExperience uint8
	JobMail       sql.NullString
}

func (o *DB_UserOffer) ToReal() *R_UserOffer {
	return &R_UserOffer{
		Id:            o.Id,
		UserId:        o.UserId,
		JobType:       o.JobType,
		JobExperience: o.JobExperience,
		JobMail:       o.JobMail.String,
	}
}

type R_UserOffer struct {
	Id            int    `json:"id"`
	UserId        int64  `json:"userId"`
	JobType       uint8  `json:"jobType"`
	JobExperience uint8  `json:"jobExperience"`
	JobMail       string `json:"jobMail,omitempty"`
}
