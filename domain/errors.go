package domain

import "errors"

var ErrNotFound = errors.New("not found")
var ErrAuthentication = errors.New("error authentication")
var ErrDuplicated = errors.New("duplicated")
var ErrNonChangeable = errors.New("non changeable")
var ErrAlreadyApproved = errors.New("already approved")
var ErrUnApproved = errors.New("unapproved")

type ErrInvalidGoogleWorkspaceAccount struct {
	Domain string
}

func (e *ErrInvalidGoogleWorkspaceAccount) Error() string { return e.Domain }
