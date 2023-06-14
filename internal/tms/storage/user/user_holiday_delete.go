package user

import (
	"errors"
	"time"
)

func (s *Storage) UserHolidayDelete(userId int64, date *time.Time) error {
	conn, err := s.repo.Conn()
	if err != nil {
		return errors.New("DB conn error: " + err.Error())
	}
	defer s.repo.Close()

	execString := `
	DELETE FROM users_holidays
	WHERE worker_id=$1 AND date_trunc('day', start_at) = $2::date`

	if err = conn.QueryRow(execString, userId, date).Err(); err != nil {
		return errors.New("DB exec error: " + err.Error())
	}
	return nil
}
