package auth

import "errors"

var ErrMissingCredentials = errors.New("empty client id or secret")
