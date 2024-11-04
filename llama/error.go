package llama

import "errors"

// LastError returns the last error that occurred if one did.
func LastError() error {
	if did_error() {
		str := get_last_error()
		if str == "" {
			return errors.New("unknown error")
		}
		return errors.New(str)
	}
	return nil
}
