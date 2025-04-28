package richerr

import "fmt"

const (

	// Client Error-4xx (4000-4999)
	ErrUnauthorized = 4001
	ErrInvalidInput = 4002
	ErrInvalidToken = 4006

	// Server Errors-5xx (5000-5999)
	ErrServer          = 5000
	ErrUnexpected      = 5001
	ErrEntityNotFound  = 5010
	ErrEntityDuplicate = 5011
)

type Operation string

type RichErr struct {
	op         Operation
	wrappedErr error
	message    string
	code       int
	metadata   map[string]interface{}
}

func New(op Operation) RichErr {
	return RichErr{
		op: op,
	}
}

func (re RichErr) WithError(err error) RichErr {
	re.wrappedErr = err
	return re
}

func (re RichErr) WithMessage(msg string) RichErr {
	re.message = msg
	return re
}

func (re RichErr) WithCode(code int) RichErr {
	re.code = code
	return re
}

func (re RichErr) WithMetadata(metadata map[string]interface{}) RichErr {
	re.metadata = metadata
	return re
}

func (err RichErr) Error() string {
	var msg string
	if err.message != "" {
		msg = err.message
	} else if err.wrappedErr != nil {
		msg = err.wrappedErr.Error()
	} else {
		msg = "no message details have been provided!!"
	}

	return fmt.Sprintf("error operational path: %s, message: %s", err.op, msg)
}

func (err RichErr) Code() int {
	if err.code != 0 {

		return err.code
	} else {
		if wErr, ok := err.wrappedErr.(RichErr); ok {

			return wErr.Code()
		} else {

			return 0
		}
	}
}

func (err RichErr) Message() string {
	if err.message != "" {

		return err.message
	} else {
		if wErr, ok := err.wrappedErr.(RichErr); ok {

			return wErr.Message()
		} else {

			return ""
		}
	}
}
