package user

import (
	"errors"
)

func (s *Storage) CreateUser(email string) (int, error) {

	var (
		result int
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return 0, errors.New("DB conn error: " + err.Error())
	}

	defer s.repo.Close()

	execString := `INSERT INTO users ( email ) 
		VALUES ( $1 )
		RETURNING id`

	if err := conn.QueryRow(execString, email).Scan(&result); err != nil {
		return 0, errors.New("DB exec error: " + err.Error())
	}
	return result, nil
}
