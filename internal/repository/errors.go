package repository

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrTagNameExists = errors.New("tag name already exists")
)
