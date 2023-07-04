package user

import (
	"encoding/json"
	"errors"

	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) UserStaffGetWorkTime() (*bl.StaffWorkTime, error) {
	var (
		workTime   bl.StaffWorkTime
		jsonRaw    []byte
		execString string
		err        error
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return nil, errors.New("DB conn error: " + err.Error())
	}
	defer s.repo.Close()

	execString = `SELECT value FROM order_variables WHERE name=$1`

	if err = conn.QueryRow(execString, "work_time").Scan(&jsonRaw); err != nil {
		return nil, errors.New("DB exec error: " + err.Error())
	}

	if err = json.Unmarshal(jsonRaw, &workTime); err != nil {
		return nil, err
	}

	return &workTime, nil
}
