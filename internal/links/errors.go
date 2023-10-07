package links

import "errors"

var (
	// ErrLinkNotFound is an error that indicates that link with the given properties was not found.
	ErrLinkNotFound = errors.New("link not found")
	// ErrIDNotGenerated is an error that indicates that it was not possible to generate an ID.
	ErrIDNotGenerated = errors.New("cannot generate ID")
	// ErrLinkAlreadyExists is an error that indicates that link with the given properties already exists.
	ErrLinkAlreadyExists = errors.New("link already exists")
)
