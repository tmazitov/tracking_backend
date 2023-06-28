package user

import (
	"errors"

	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) UserOfferCreate(userId int64, job bl.UserJob) (int, error) {

	var (
		offerId int
		err     error
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return offerId, errors.New("DB conn error: " + err.Error())
	}
	defer s.repo.Close()

	execString := `
	INSERT INTO 
	users_job (user_id, job_type, job_experience, job_mail) 
	VALUES ( $1, $2, $3, $4 )
	RETURNING id`

	if err = conn.QueryRow(execString, userId, job.JobType, job.JobExperience, job.JobMail).Scan(&offerId); err != nil {
		return offerId, errors.New("DB exec error: " + err.Error())
	}
	return offerId, nil
}
