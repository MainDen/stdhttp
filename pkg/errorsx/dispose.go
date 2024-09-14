package errorsx

import "errors"

func Dispose(err *error, dispose func() error) {
	*err = errors.Join(*err, dispose())
}
