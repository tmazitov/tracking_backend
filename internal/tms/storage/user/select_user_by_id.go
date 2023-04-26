package user

import (
	"errors"

	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) GetUserInfo(userId int64) (bl.DB_GetUser, error) {
	var (
		result bl.DB_GetUser
		err    error
	)

	result = bl.DB_GetUser{}
	conn, err := s.repo.Conn()
	if err != nil {
		return result, errors.New("DB conn error: " + err.Error())
	}

	defer s.repo.Close()

	execString := `SELECT  short_name, role FROM users WHERE id=$1`

	if err := conn.QueryRow(execString, userId).Scan(&result.ShortName, &result.RoleID); err != nil {
		return result, errors.New("DB exec error: " + err.Error())
	}

	return result, err
}
