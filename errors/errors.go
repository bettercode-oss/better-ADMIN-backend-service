package errors

import "github.com/pkg/errors"

var (
	ErrNotFound                  = errors.New("not found")
	ErrAuthentication            = errors.New("error authentication")
	ErrDuplicated                = errors.New("duplicated")
	ErrNonChangeable             = errors.New("non changeable")
	ErrAlreadyApproved           = errors.New("already approved")
	ErrUnApproved                = errors.New("unapproved")
	ErrNotSupportedAccessLogType = errors.New("not supported access log type")
)

type ErrInvalidGoogleWorkspaceAccount struct {
	Domain string
}

func (e *ErrInvalidGoogleWorkspaceAccount) Error() string { return e.Domain }
