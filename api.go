package main

import "errors"

type Authentication interface {
	setAuthToken()
}

func AuthenticationError() error {
	return errors.New("Authentication failed.")
}

func UnknownError() error {
	return errors.New("Unknown Error.")
}
