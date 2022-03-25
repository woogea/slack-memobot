package mention

type NotfoundError struct{}

func (err *NotfoundError) Error() string {
	return "Not Found"
}

type InvalidArgumentError struct {
	msg string
}

func (err *InvalidArgumentError) Error() string {
	return err.msg
}
