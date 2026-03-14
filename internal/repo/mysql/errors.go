package mysql

import "errors"

var (
	ErrNilDB           = errors.New("mysql repo: nil db")
	ErrInvalidArgument = errors.New("mysql repo: invalid argument")
	ErrNotFound        = errors.New("mysql repo: record not found")
	ErrConflict        = errors.New("mysql repo: conflict")
)
