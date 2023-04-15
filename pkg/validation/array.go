package validation

func ValidateArrayMaxMin(items []interface{}, max int, min int) error {
	if len(items) > max || len(items) < min {
		return ErrNotValid
	}
	return nil
}
