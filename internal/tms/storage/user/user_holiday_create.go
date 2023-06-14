package user

import (
	"errors"
	"time"
)

func (s *Storage) UserHolidayCreate(userId int64, authorId int64, date *time.Time) error {
	conn, err := s.repo.Conn()
	if err != nil {
		return errors.New("DB conn error: " + err.Error())
	}
	defer s.repo.Close()

	execString := `
	INSERT INTO users_holidays (worker_id, author_id, start_at) 
	VALUES (
		$1,
		$2,
		$3
	)`

	if err = conn.QueryRow(execString, userId, authorId, date).Err(); err != nil {
		return errors.New("DB exec error: " + err.Error())
	}
	return nil
}
