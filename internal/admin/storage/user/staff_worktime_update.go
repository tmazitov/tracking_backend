package user

import (
	"encoding/json"
	"errors"

	"github.com/tmazitov/tracking_backend.git/internal/admin/bl"
)

func (s *Storage) StaffWorkTimeUpdate(workTime *bl.StaffWorkTime) error {
	var (
		execString string
		err        error
	)

	jsonVariables, err := json.Marshal(workTime)
	if err != nil {
		return err
	}

	conn, err := s.repo.Conn()
	if err != nil {
		return errors.New("DB conn error: " + err.Error())
	}
	defer s.repo.Close()

	execString = `UPDATE order_variables SET value = $2 WHERE name=$1`

	if err = conn.QueryRow(execString, "work_time", jsonVariables).Err(); err != nil {
		return errors.New("DB exec error: " + err.Error())
	}
	return nil
}
