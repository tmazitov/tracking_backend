package user

import (
	"errors"
)

func (s *Storage) CheckUserByEmail(email string) (int, int, error) {

	var (
		// result []int
		userId int
		roleId int
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return -1, -1, errors.New("DB conn error: " + err.Error())
	}

	defer s.repo.Close()

	execString := `SELECT id, role_id FROM users WHERE email=$1 `
	rows, err := conn.Query(execString, email)
	if err != nil {
		return -1, -1, errors.New("DB exec error: " + err.Error())
	}
	for rows.Next() {
		err = rows.Scan(&userId, &roleId)
		if err != nil {
			return -1, -1, err
		}
	}
	return userId, roleId, nil
}
