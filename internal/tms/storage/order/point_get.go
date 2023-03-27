package order

import (
	"errors"
	"fmt"

	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func getPointsIdString(pointsID []int64) string {

	result := ""
	for index, id := range pointsID {
		result += fmt.Sprintf("%d", id)
		if index != len(pointsID)-1 {
			result += ", "
		}
	}
	return result
}

func (s *Storage) GetPoints(pointsID []int64) ([]bl.Point, error) {
	var (
		result []bl.Point
		err    error
	)

	result = []bl.Point{}
	conn, err := s.gis.Conn()
	if err != nil {
		return result, fmt.Errorf("DB conn error: " + err.Error())
	}

	defer s.gis.Close()

	pointsString := getPointsIdString(pointsID)
	execString := fmt.Sprintf("SELECT floor, title, ST_X(ST_AsText(point)), ST_Y(ST_AsText(point)) FROM points WHERE id IN ( %s)", pointsString)

	rows, err := conn.Query(execString)
	if err != nil {
		return result, errors.New("DB exec error: " + err.Error())
	}
	for rows.Next() {
		p := bl.Point{}
		err := rows.Scan(
			&p.Floor,
			&p.Title,
			&p.Latitude,
			&p.Longitude,
		)
		if err != nil {
			return result, errors.New("DB read error: " + err.Error())
		}
		result = append(result, p)
	}
	return result, err
}
