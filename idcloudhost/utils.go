package idcloudhost

import (
	"errors"
	"net/http"
)

func AuthenticationError() error {
	return errors.New("authentication failed")
}

func UnknownError() error {
	return errors.New("unknown error")
}

func NotImplementedError() error {
	return errors.New("not implemented")
}

func checkError(StatusCode int) error {
	if StatusCode != http.StatusOK {
		if StatusCode == http.StatusForbidden {
			return AuthenticationError()
		}
		if StatusCode == http.StatusUnauthorized {
			return AuthenticationError()
		}
		return UnknownError()
	}
	return nil
}
