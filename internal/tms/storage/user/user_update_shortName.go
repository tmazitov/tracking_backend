package user

import "errors"

func (s *Storage) UpdateUserShortName(userId int, shortName string) error {
	conn, err := s.repo.Conn()
	if err != nil {
		return errors.New("DB conn error: " + err.Error())
	}
	defer s.repo.Close()

	execString := `
	UPDATE users SET 
		short_name=$1,
		edited_at=now()
	WHERE id=$2`

	if err = conn.QueryRow(execString, shortName, userId).Err(); err != nil {
		return errors.New("DB exec error: " + err.Error())
	}
	return nil
}
