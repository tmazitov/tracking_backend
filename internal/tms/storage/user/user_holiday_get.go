package user

import (
	"errors"
	"time"

	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) UserHolidayGet(userId int64, date *time.Time) (*bl.UserHoliday, error) {

	var (
		holiday bl.UserHoliday = bl.UserHoliday{}
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return nil, errors.New("DB conn error: " + err.Error())
	}
	defer s.repo.Close()

	execString := `
	SELECT worker_id, author_id, start_at FROM users_holidays 
	WHERE worker_id=$1 AND date_trunc('day', start_at) = $2::date`

	if err = conn.QueryRow(execString, userId, date).Scan(
		&holiday.WorkerId, &holiday.AuthorId, &holiday.Date,
	); err != nil {
		return nil, errors.New("DB exec error: " + err.Error())
	}
	return &holiday, nil
}
