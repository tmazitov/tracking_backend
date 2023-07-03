package order

import (
	"encoding/json"
	"errors"

	"github.com/tmazitov/tracking_backend.git/internal/admin/bl"
)

func (s *Storage) OrderPricelistUpdate(priceList *bl.OrderPriceList) error {
	var (
		execString string
		err        error
	)

	jsonVariables, err := json.Marshal(priceList)
	if err != nil {
		return err
	}

	conn, err := s.repo.Conn()
	if err != nil {
		return errors.New("DB conn error: " + err.Error())
	}
	defer s.repo.Close()

	execString = `UPDATE order_variables SET value = $2 WHERE name=$1`

	if err = conn.QueryRow(execString, "prices", jsonVariables).Err(); err != nil {
		return errors.New("DB exec error: " + err.Error())
	}

	return nil
}
