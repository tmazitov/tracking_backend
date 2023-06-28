package user

import (
	"database/sql"
	"errors"
)

func (s *Storage) UserOfferGet(userId int64) (int, error) {
	var (
		offerId    int
		execString string
		err        error
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return offerId, errors.New("DB conn error: " + err.Error())
	}
	defer s.repo.Close()

	execString = `
	SELECT id FROM users_job  
	WHERE user_id=$1 AND accepted_at IS NULL AND canceled_at IS NULL
	LIMIT 1`

	err = conn.QueryRow(execString, userId).Scan(&offerId)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return offerId, errors.New("DB exec error: " + err.Error())
	}

	return offerId, nil
}
