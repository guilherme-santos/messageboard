package messageboard

import (
	"fmt"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewError(code, msg string) error {
	return &Error{
		Code:    code,
		Message: msg,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e Error) Is(target error) bool {
	_, ok := target.(*Error)
	return ok
}
