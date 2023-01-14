package order

import (
	"errors"

	"github.com/lib/pq"
	"github.com/tmazitov/tracking_backend.git/internal/app/bl"
)

func (s *Storage) InsertOrder(order bl.CreateOrder) error {
	// startAt time.Time
	// points  []Point
	// helpers uint8
	// comment string
	conn, err := s.repo.Conn()
	if err != nil {
		return errors.New("DB conn error: " + err.Error())
	}

	defer s.repo.Close()

	// points := "{"
	// for index, pointID := range order.PointsID {
	// 	points += fmt.Sprintf("%d", pointID)
	// 	if index != len(order.PointsID)-1 {
	// 		points += ","
	// 	}
	// }
	// points = points + "}"

	execString := "insert into orders (startAt, endAt, points, helpers, comment_message) values ($1, $2, $3, $4, $5)"
	_, err = conn.Exec(execString, order.StartAt, order.EndAt, pq.Array(order.PointsID), order.Helpers, order.Comment)
	if err != nil {
		return errors.New("DB exec error: " + err.Error())
	}

	return nil
}
