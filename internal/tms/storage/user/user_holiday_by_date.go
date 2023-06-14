package user

import (
	"errors"
	"time"
)

func (s *Storage) UserHolidayListByDate(date time.Time) ([]int64, error) {

	var (
		users []int64
		err   error
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return users, errors.New("DB conn error: " + err.Error())
	}
	defer s.repo.Close()

	execString := `
	SELECT worker_id FROM users_holidays
	WHERE date_trunc('day', start_at) = $1::date`

	rows, err := conn.Query(execString, date)
	if err != nil {
		return users, errors.New("DB exec error: " + err.Error())
	}

	for rows.Next() {
		var userId int64
		err = rows.Scan(&userId)
		if err != nil {
			return users, errors.New("DB scan error: " + err.Error())
		}
		users = append(users, userId)
	}

	return users, nil
}
