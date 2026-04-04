package ordererrors

import "errors"

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrOrderNotFound   = errors.New("order not found")
)
