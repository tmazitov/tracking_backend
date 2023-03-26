package order

import (
	"errors"
	"fmt"

	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) InsertPoint(points []bl.Point) ([]int64, error) {

	result := []int64{}
	conn, err := s.gis.Conn()
	if err != nil {
		return nil, errors.New("DB conn error: " + err.Error())
	}

	defer s.gis.Close()

	execString := fmt.Sprintf("INSERT INTO points (title, floor, point) VALUES %s RETURNING id", addPoints(points))
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

func addPoints(points []bl.Point) string {
	values := ""
	for index, point := range points {
		values += fmt.Sprintf("('%s', '%d', ST_GeomFromText('POINT(%g %g)',4326))", point.Title, point.Floor, point.Longitude, point.Latitude)

		if index != len(points)-1 {
			values += ", "
		}
	}
	return values
}
