package cfmgo

import "errors"

var (
	//ErrInvalidID -- error for invalid id
	ErrInvalidID = errors.New("value is not a properly formatted hex string")
)
