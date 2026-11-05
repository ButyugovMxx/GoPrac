package domain

import "errors"

var ErrInvalidRange = errors.New("invalid movie range")
var NotFound = errors.New("film not found")
