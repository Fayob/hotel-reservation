package api

import "net/http"

type Error struct {
	Code int `json:"code"`
	Err string `json:"err"`
}

// Error implement the Error interface
func (e Error) Error() string {
	return e.Err
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err: err,
	}
}

func ErrResourceNotFound(res string) Error {
	return Error{
		Code: http.StatusNotFound,
		Err: res + "resource not found",
	}
}

func ErrBadRequest() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err: "invalid JSON request",
	}
}

func ErrUnAuthorized() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err: "unauthorized request",
	}
}

func ErrInvalidID() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err: "invalid id given",
	}
}