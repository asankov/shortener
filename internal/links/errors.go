package links

import "errors"

var (
	ErrLinkNotFound      = errors.New("link not found")
	ErrIDNotGenerated    = errors.New("cannot generate ID")
	ErrLinkAlreadyExists = errors.New("link already exists")
)
