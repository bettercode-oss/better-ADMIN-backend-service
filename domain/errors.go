package domain

import "errors"

var ErrNotFound = errors.New("not found")
var ErrAuthentication = errors.New("error authentication")
var ErrDuplicated = errors.New("duplicated")
var ErrNonChangeable = errors.New("non changeable")
