package errors

type StandardError struct {
	msg string
}

func (e StandardError) Error() string {
	return e.msg
}
