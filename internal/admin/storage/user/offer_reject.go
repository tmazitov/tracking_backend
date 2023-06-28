package user

import (
	"errors"
)

func (s *Storage) OfferReject(offerId int) error {
	var (
		err error
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return errors.New("DB conn error: " + err.Error())
	}
	defer s.repo.Close()

	execString := `
	UPDATE users_job SET 
		canceled_at=NOW()
	WHERE id=$1 `

	if err = conn.QueryRow(execString, offerId).Err(); err != nil {
		return errors.New("DB exec error: " + err.Error())
	}

	return nil
}
