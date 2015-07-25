package xonotic

import "errors"

// Various errors this library can throw
var (
	ErrInvalidFormat = errors.New("xonotic: invalid server reply")
)
