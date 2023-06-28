package user

import (
	"errors"

	"github.com/tmazitov/tracking_backend.git/internal/admin/bl"
)

func (s *Storage) StaffRemove(userId int64) error {
	var (
		err error
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return errors.New("DB conn error: " + err.Error())
	}
	defer s.repo.Close()

	execString := `
	UPDATE users SET 
		role=$2
	WHERE id=$1 `

	if err = conn.QueryRow(execString, userId, bl.Base).Err(); err != nil {
		return errors.New("DB exec error: " + err.Error())
	}

	return nil
}
