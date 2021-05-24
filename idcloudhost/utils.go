package idcloudhost

import (
	"errors"
	"net/http"
)

func AuthenticationError() error {
	return errors.New("Authentication failed.")
}

func UnknownError() error {
	return errors.New("Unknown Error.")
}

func NotImplementedError() error {
	return errors.New("Not Implemented.")
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
