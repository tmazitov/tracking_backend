package order

import (
	"encoding/json"
	"errors"

	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) OrderPriceList() (*bl.OrderPriceList, error) {

	var (
		priceList     bl.OrderPriceList
		priceListJson []byte
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return nil, errors.New("DB conn error: " + err.Error())
	}

	defer s.repo.Close()

	execString := `SELECT value FROM order_variables WHERE name = 'prices'`

	if err := conn.QueryRow(execString).Scan(&priceListJson); err != nil {
		return nil, errors.New("DB exec error: " + err.Error())
	}

	if err = json.Unmarshal(priceListJson, &priceList); err != nil {
		return nil, nil
	}

	return &priceList, nil
}
