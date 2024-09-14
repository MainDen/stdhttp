package errorsx

import "strings"

const SecretString = "SECRET"

type secretError struct {
	secret string
	err    error
}

func (e *secretError) Error() string {
	return strings.ReplaceAll(e.err.Error(), e.secret, SecretString)
}

func (e *secretError) Unwrap() error {
	return e.err
}

func WrapSecret(secret string, err error) error {
	if err == nil {
		return nil
	}
	return &secretError{secret: secret, err: err}
}
