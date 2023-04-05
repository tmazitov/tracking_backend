package order

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) PointsCreate(orderID int, points []bl.Point) ([]int64, error) {

	conn, err := s.gis.Conn()
	if err != nil {
		return nil, errors.New("DB conn error: " + err.Error())
	}

	tx, err := conn.Begin()
	if err != nil {
		return nil, errors.New("DB transaction begin err: " + err.Error())
	}

	defer s.gis.Close()

	var (
		pointsID            []int64
		orderToPointValues  []interface{}
		createOrderValues   []interface{}
		unknownValuesString string
	)
	createOrderValues, unknownValuesString = getCreatablePointsValues(points)

	execString := fmt.Sprintf(`INSERT INTO points 
	(title, step_id, floor, point) 
	VALUES %s
	RETURNING id`, unknownValuesString)

	rows, err := conn.Query(execString, createOrderValues...)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("DB exec error: " + err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			tx.Rollback()
			return nil, errors.New("DB get data error: " + err.Error())
		}
		pointsID = append(pointsID, id)
	}

	var counter int = 1
	unknownValuesString = ""
	for index, pointID := range pointsID {
		orderToPointValues = append(orderToPointValues, orderID, pointID)
		unknownValuesString += fmt.Sprintf("($%d, $%d)", counter, counter+1)
		if index != len(points)-1 {
			unknownValuesString += ", "
		}
		counter += 2
	}

	execString = fmt.Sprintf(`INSERT INTO points_to_orders (order_id, point_id) VALUES %s`, unknownValuesString)
	if err = tx.QueryRow(execString, orderToPointValues...).Err(); err != nil {
		tx.Rollback()
		return nil, errors.New("DB exec error: " + err.Error())
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return nil, errors.New("DB transaction continue err: " + err.Error())
	}

	return pointsID, nil
}

func getCreatablePointsValues(points []bl.Point) ([]interface{}, string) {
	var (
		values              []interface{}
		unknownValuesString string
	)
	unknownValuesString = ""
	for index, point := range points {
		values = append(values, point.Title, point.StepID, point.Floor)
		unknownValuesString += getCreatableUnknownValueItem(index*3+1, point)
		if index != len(points)-1 {
			unknownValuesString += ", "
		}
	}
	return values, unknownValuesString
}

func getCreatableUnknownValueItem(startWith int, point bl.Point) string {
	var unknownValueItem = "("
	for i := startWith; i < startWith+3; i++ {
		unknownValueItem += "$" + strconv.Itoa(i)
		unknownValueItem += ", "
	}
	unknownValueItem += fmt.Sprintf("ST_GeomFromText('POINT(%g %g)',4326)", point.Longitude, point.Latitude) + ")"
	return unknownValueItem
}
