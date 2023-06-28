package user

import (
	"errors"

	"github.com/tmazitov/tracking_backend.git/internal/admin/bl"
)

func (s *Storage) OfferGet(offerId int) (*bl.DB_UserOffer, error) {
	var (
		offer      bl.DB_UserOffer
		execString string
		err        error
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return &offer, errors.New("DB conn error: " + err.Error())
	}
	defer s.repo.Close()

	execString = `
	SELECT id, user_id, job_type, job_experience, job_mail FROM users_job  
	WHERE id=$1 AND accepted_at IS NULL AND canceled_at IS NULL`

	err = conn.QueryRow(execString, offerId).Scan(&offer.Id, &offer.UserId, &offer.JobType, &offer.JobExperience, &offer.JobMail)
	if err != nil {
		return nil, errors.New("DB exec error: " + err.Error())
	}

	return &offer, nil
}
