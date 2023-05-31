package order

import (
	"errors"

	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) OrderDefaultPrices() ([]bl.DefaultPriceItems, error) {

	var (
		defaultPrices []bl.DefaultPriceItems
		item          bl.DefaultPriceItems
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return nil, errors.New("DB conn error: " + err.Error())
	}

	defer s.repo.Close()

	execString := `
		SELECT name, value FROM default_prices
	`

	rows, err := conn.Query(execString)
	if err != nil {
		return nil, errors.New("DB exec error: " + err.Error())
	}

	for rows.Next() {
		item = bl.DefaultPriceItems{}
		err = rows.Scan(&item.Name, &item.Value)
		if err != nil {
			return nil, errors.New("DB read error: " + err.Error())
		}
		defaultPrices = append(defaultPrices, item)
	}

	return defaultPrices, nil
}
