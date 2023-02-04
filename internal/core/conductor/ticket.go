package conductor

type Ticket struct {
	Code string
}

func (t *Ticket) ValidateCode(code string) error {
	if code != t.Code {
		return ErrInvalidToken
	}

	return nil
}
