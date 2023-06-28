package user

import (
	"database/sql"
	"errors"

	"github.com/tmazitov/tracking_backend.git/internal/admin/bl"
)

func (s *Storage) OfferAccept(offerId int) (*bl.DB_User, error) {
	var (
		userId     int64
		jobType    int8
		userRaw    bl.DB_User
		execString string
		tx         *sql.Tx
		err        error
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return nil, errors.New("DB conn error: " + err.Error())
	}
	defer s.repo.Close()

	tx, err = conn.Begin()
	if err != nil {
		return nil, errors.New("DB transaction begin error: " + err.Error())
	}

	execString = `
	UPDATE users_job SET
		accepted_at=NOW()
	WHERE id=$1
	RETURNING user_id, job_type
	`
	if err = conn.QueryRow(execString, offerId).Scan(&userId, &jobType); err != nil {
		tx.Rollback()
		return nil, errors.New("DB exec error: " + err.Error())
	}

	execString = `
	UPDATE users SET 
		role=$2
	WHERE id=$1
	RETURNING id, short_name, role`

	if err = conn.QueryRow(execString, userId, jobType).Scan(&userRaw.ID, &userRaw.ShortName, &userRaw.RoleID); err != nil {
		tx.Rollback()
		return nil, errors.New("DB exec error: " + err.Error())
	}

	return &userRaw, nil
}
