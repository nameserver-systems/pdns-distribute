package certificate

import "errors"

var (
	errCertPathNotAbsolute = errors.New("certificate or keypath are not absolute")
)
