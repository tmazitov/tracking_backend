package order

import (
	"database/sql"
	"errors"
	"fmt"
)

func (s *Storage) PointsDelete(pointsID []int64) error {

	var (
		tx           *sql.Tx
		execString   string
		pointsString string
		err          error
	)

	conn, err := s.gis.Conn()
	if err != nil {
		return errors.New("DB conn error: " + err.Error())
	}

	tx, err = conn.Begin()
	if err != nil {
		return errors.New("DB transaction begin error: " + err.Error())
	}

	defer s.gis.Close()

	pointsString = getPointsIdString(pointsID)
	execString = fmt.Sprintf("DELETE FROM points_to_orders WHERE point_id IN ( %s)", pointsString)
	if err = tx.QueryRow(execString).Err(); err != nil {
		tx.Rollback()
		return errors.New("DB exec error: " + err.Error())
	}

	execString = fmt.Sprintf("DELETE FROM points WHERE id IN ( %s)", pointsString)

	if err = tx.QueryRow(execString).Err(); err != nil {
		tx.Rollback()
		return errors.New("DB exec error: " + err.Error())
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return errors.New("DB transaction commit error: " + err.Error())
	}

	return nil
}
