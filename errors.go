package certrenewal

import "errors"

var (
	// ErrConfig will be returned if the config has a problem.
	ErrConfig = errors.New("cannot read config")
	// ErrLogin reflects a login problem with the vault server.
	ErrLogin = errors.New("cannot login")
	// ErrIssue reflects a problem with the cert renewal.
	ErrIssue = errors.New("cannot issue new certificate")
	// ErrCert reflects a problem with the cert.
	ErrCert       = errors.New("problem with cert")
	ErrHookFailed = errors.New("hook failed")
)
