package auth

import "errors"

// ErrMissingCredentials if client id or secret are missing
var ErrMissingCredentials = errors.New("empty client id or secret")
