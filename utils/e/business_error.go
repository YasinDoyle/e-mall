package e

import "errors"

type BusinessError struct {
	Code int
	Err  error
}

func (e *BusinessError) Error() string {
	if e == nil || e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

func (e *BusinessError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func NewBusinessError(code int, message ...string) error {
	msg := GetMsg(code)
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	return &BusinessError{
		Code: code,
		Err:  errors.New(msg),
	}
}

func CodeFromError(err error) (int, bool) {
	var bizErr *BusinessError
	if errors.As(err, &bizErr) && bizErr != nil {
		return bizErr.Code, true
	}
	return 0, false
}
