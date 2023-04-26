package user

import (
	"errors"

	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) UserGetStaffList() ([]bl.DB_GetUser, error) {
	var (
		staff []bl.DB_GetUser
		user  bl.DB_GetUser
		err   error
	)

	staff = []bl.DB_GetUser{}
	conn, err := s.repo.Conn()
	if err != nil {
		return nil, errors.New("DB conn error: " + err.Error())
	}

	defer s.repo.Close()

	execString := `SELECT id, short_name, role FROM users WHERE role=$1 OR role=$2`

	rows, err := conn.Query(execString, bl.Manager, bl.Worker)
	if err != nil {
		return nil, errors.New("DB exec error: " + err.Error())
	}

	for rows.Next() {
		err = rows.Scan(&user.ID, &user.ShortName, &user.RoleID)
		if err != nil {
			return nil, errors.New("DB exec error: " + err.Error())
		}
		staff = append(staff, user)
	}

	return staff, nil
}
