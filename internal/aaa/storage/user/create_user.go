package user

import (
	"errors"
)

func (s *Storage) CreateUser(email string) (int64, error) {

	var (
		userId int64
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return 0, errors.New("DB conn error: " + err.Error())
	}

	defer s.repo.Close()

	execString := `INSERT INTO users ( email ) 
		VALUES ( $1 )
		RETURNING id`

	if err := conn.QueryRow(execString, email).Scan(&userId); err != nil {
		return 0, errors.New("DB exec error: " + err.Error())
	}
	return userId, nil
}
