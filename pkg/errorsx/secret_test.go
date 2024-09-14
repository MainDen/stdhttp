package errorsx

import (
	"errors"
	"testing"
)

func TestSecret(t *testing.T) {
	var err error
	err = WrapSecret("secret_value1", err)
	if err != nil {
		t.Error("expected nil")
	}
	err1 := errors.New("error with secret_value1")
	err = WrapSecret("secret_value1", err1)
	if !errors.Is(err, err1) {
		t.Error("expected err1")
	}
	if err.Error() != "error with SECRET" {
		t.Error("expected error with SECRET")
	}
	err2 := errors.New("error with secret_value1 and secret_value2")
	err = WrapSecret("secret_value1", WrapSecret("secret_value2", err2))
	if !errors.Is(err, err2) {
		t.Error("expected err2")
	}
	if err.Error() != "error with SECRET and SECRET" {
		t.Error("expected error with SECRET and SECRET")
	}
}
