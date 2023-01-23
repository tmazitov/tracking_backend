package order

import (
	"errors"
	"fmt"

	"github.com/tmazitov/tracking_backend.git/internal/app/bl"
)

func (s *Storage) InsertPoint(points []bl.Point) ([]int64, error) {

	result := []int64{}
	conn, err := s.gis.Conn()
	if err != nil {
		return nil, errors.New("DB conn error: " + err.Error())
	}

	defer s.gis.Close()

	execString := "insert into points (title, floor, point) values"
	values := ""

	for index, point := range points {
		values += fmt.Sprintf("('%s', '%d', 'SRID=4326;POINT(%g %g)')", point.Title, point.Floor, point.Longitude, point.Latitude)

		if index != len(points)-1 {
			values += ", "
		}
	}

	execString = execString + values + " returning id"

	rows, err := conn.Query(execString)
	if err != nil {
		return nil, errors.New("DB exec error: " + err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, errors.New("DB get data error: " + err.Error())
		}
		result = append(result, id)
	}

	return result, nil
}
