package idcloudhost

import (
	"errors"
	"net/http"
)

func AuthenticationError() error {
	return errors.New("authentication failed, check authentication token and billing account")
}

func BadRequestError() error {
	return errors.New("bad request error, check request data")
}

func UnknownError() error {
	return errors.New("unknown error")
}

func NotFoundError() error {
	return errors.New("resource not found")
}

func NotImplementedError() error {
	return errors.New("not implemented")
}

func checkError(StatusCode int) error {
	switch StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusForbidden:
		return AuthenticationError()
	case http.StatusUnauthorized:
		return AuthenticationError()
	case http.StatusBadRequest:
		return BadRequestError()
	case http.StatusNotFound:
		return NotFoundError()
	default:
		return UnknownError()
	}
}
