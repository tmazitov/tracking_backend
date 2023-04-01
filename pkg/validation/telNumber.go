package validation

import "regexp"

func ValidateTelNumber(tel string) error {
	valid, err := regexp.MatchString(tel, `^(\+7|7|8)?[\s\-]?\(?[489][0-9]{2}\)?[\s\-]?[0-9]{3}[\s\-]?[0-9]{2}[\s\-]?[0-9]{2}$`)
	if err != nil {
		return err
	}
	if !valid {
		return ErrNotValid
	}

	return nil
}
