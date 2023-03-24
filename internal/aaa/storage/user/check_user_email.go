package user

import (
	"errors"
)

func (s *Storage) CheckUserByEmail(email string) (bool, error) {

	var (
		result bool
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return false, errors.New("DB conn error: " + err.Error())
	}

	defer s.repo.Close()

	execString := `SELECT EXISTS ( 
		SELECT 1 FROM users * WHERE email=$1 
	)`
	rows, err := conn.Query(execString, email)
	if err != nil {
		return false, errors.New("DB exec error: " + err.Error())
	}
	// result = strconv.Atoi(rows[0])
	for rows.Next() {
		err := rows.Scan(&result)
		if err != nil {
			return false, errors.New("DB exec error: " + err.Error())
		}
	}
	return result, nil
}
