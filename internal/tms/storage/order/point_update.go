package order

import (
	"errors"
	"fmt"

	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) PointsUpdate(points []bl.Point) ([]int64, error) {
	var (
		updatedPointID int64
		pointsID       []int64
		err            error
	)

	conn, err := s.gis.Conn()
	if err != nil {
		return nil, fmt.Errorf("DB conn error: " + err.Error())
	}

	defer s.gis.Close()

	tx, err := conn.Begin()
	if err != nil {
		return nil, fmt.Errorf("DB transaction begin error: " + err.Error())
	}

	execString := `UPDATE points
	SET
		step_id = $2,
		floor   = $3,
		title   = $4,
		point   = ST_GeomFromText($5, 4326)
	WHERE id=$1
	RETURNING id`

	for _, point := range points {
		if err = tx.QueryRow(execString, point.ID, point.StepID,
			point.Floor, point.Title, fmt.Sprintf("POINT(%g %g)", point.Longitude, point.Latitude)).Scan(&updatedPointID); err != nil {
			tx.Rollback()
			return nil, errors.New("DB exec error: " + err.Error())
		}
		pointsID = append(pointsID, updatedPointID)
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("DB transaction commit error: " + err.Error())
	}

	fmt.Println(pointsID)

	return pointsID, nil
}
