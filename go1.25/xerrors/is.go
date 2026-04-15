package xerrors

import "errors"

func Is[T error](err error) bool {
	var targetErr T
	return errors.As(err, &targetErr)
}
