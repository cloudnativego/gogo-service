package cfmgo

import (
	"errors"
	"net/url"

	"github.com/cloudnativego/cfmgo/params"
)

var (
	//ErrInvalidID -- error for invalid id
	ErrInvalidID = errors.New("value is not a properly formatted hex string")

	//ParamsUnfiltered returns a default set of parameters
	ParamsUnfiltered = params.Extract(url.Values{})
)
