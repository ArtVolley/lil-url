package storage

import "errors"

var (
	ErrUrlNotFound = errors.New("url not found")
	ErrLilNotFound = errors.New("lil not found")
	ErrUrlExists   = errors.New("url exists")
)
