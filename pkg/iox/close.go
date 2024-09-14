package iox

import (
	"errors"
	"io"
)

func Close(streams ...any) error {
	var errs error
	for _, stream := range streams {
		if closer, ok := stream.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				errs = errors.Join(errs, err)
			}
		}
	}
	return errs
}
