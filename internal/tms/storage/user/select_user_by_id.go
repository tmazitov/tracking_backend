package user

import (
	"errors"

	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) GetUserInfo(userId int) (bl.SelectUserById, error) {
	var (
		result bl.SelectUserById
		err    error
	)

	result = bl.SelectUserById{}
	conn, err := s.repo.Conn()
	if err != nil {
		return result, errors.New("DB conn error: " + err.Error())
	}

	defer s.repo.Close()

	execString := `SELECT telephone, short_name, role FROM users WHERE id=$1`

	if err := conn.QueryRow(execString, userId).Scan(&result.TelNumber, &result.ShortName, &result.RoleID); err != nil {
		return result, errors.New("DB exec error: " + err.Error())
	}

	return result, err
}
