package errorsx

import (
	"errors"
	"testing"
)

func TestDispose(t *testing.T) {
	dispose := func(err error) func() error {
		return func() error {
			return err
		}
	}
	var err error
	Dispose(&err, dispose(nil))
	if err != nil {
		t.Error("expected nil")
	}
	err1 := errors.New("err1")
	err2 := errors.New("err2")
	Dispose(&err, dispose(err1))
	Dispose(&err, dispose(err2))
	if !errors.Is(err, err1) {
		t.Error("expected err1")
	}
	if !errors.Is(err, err2) {
		t.Error("expected err2")
	}
}
