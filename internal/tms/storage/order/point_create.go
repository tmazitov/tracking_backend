package order

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) CreatePoints(points []bl.Point) ([]int64, error) {

	result := []int64{}
	conn, err := s.gis.Conn()
	if err != nil {
		return nil, errors.New("DB conn error: " + err.Error())
	}

	defer s.gis.Close()

	var (
		values              []interface{}
		unknownValuesString string
	)
	values, unknownValuesString = getPointsValues(points)

	execString := fmt.Sprintf(`INSERT INTO points 
	(title, step_id, floor, point) 
	VALUES %s
	RETURNING id`, unknownValuesString)

	rows, err := conn.Query(execString, values...)
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

func getPointsValues(points []bl.Point) ([]interface{}, string) {
	var (
		values              []interface{}
		unknownValuesString string
	)
	unknownValuesString = ""
	for index, point := range points {
		values = append(values, point.Title, point.StepID, point.Floor)
		unknownValuesString += getUnknownValueItem(index*3+1, point)
		if index != len(points)-1 {
			unknownValuesString += ", "
		}
	}
	return values, unknownValuesString
}

func getUnknownValueItem(startWith int, point bl.Point) string {
	var unknownValueItem = "("
	for i := startWith; i < startWith+3; i++ {
		unknownValueItem += "$" + strconv.Itoa(i)
		unknownValueItem += ", "
	}
	unknownValueItem += fmt.Sprintf("ST_GeomFromText('POINT(%g %g)',4326)", point.Longitude, point.Latitude) + ")"
	return unknownValueItem
}
