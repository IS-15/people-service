package storage

import "errors"

var (
	ErrPersonNotFound = errors.New("person not found")
	ErrPersonExists   = errors.New("person exists")
)
