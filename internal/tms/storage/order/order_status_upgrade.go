package order

import "errors"

func (s *Storage) OrderStatusUpgrade(orderId int) (int, error) {
	conn, err := s.repo.Conn()
	if err != nil {
		return 0, errors.New("DB conn error: " + err.Error())
	}
	defer s.repo.Close()

	var currentStatusId int
	execString := `
	UPDATE orders SET 
		status_id=CASE status_id
			WHEN 0 THEN 1
			WHEN 1 THEN 2
			WHEN 2 THEN 3
			WHEN 3 THEN 4
		END,
		edited_at=now()
	WHERE id=$1 RETURNING status_id`

	if err = conn.QueryRow(execString, orderId).Scan(&currentStatusId); err != nil {
		return 0, errors.New("DB exec error: " + err.Error())
	}

	return currentStatusId, nil
}
