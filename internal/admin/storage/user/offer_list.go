package user

import (
	"errors"

	"github.com/tmazitov/tracking_backend.git/internal/admin/bl"
)

func (s *Storage) UserOfferList() ([]bl.DB_UserOffer, error) {
	var (
		offers     []bl.DB_UserOffer
		execString string
		err        error
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return offers, errors.New("DB conn error: " + err.Error())
	}
	defer s.repo.Close()

	execString = `
	SELECT 
		id,
		user_id,
		job_type,
		job_experience,
		job_mail
	FROM users_job  
	WHERE accepted_at IS NULL AND canceled_at IS NULL`

	rows, err := conn.Query(execString)
	if err != nil {
		return offers, errors.New("DB exec error: " + err.Error())
	}

	for rows.Next() {
		var offer bl.DB_UserOffer
		err = rows.Scan(&offer.Id, &offer.UserId, &offer.JobType, &offer.JobExperience, &offer.JobMail)
		if err != nil {
			return offers, errors.New("DB get data error: " + err.Error())
		}
		offers = append(offers, offer)
	}

	return offers, nil
}
