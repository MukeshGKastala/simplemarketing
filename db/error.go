package db

import "errors"

var ErrDuplicateEntry = errors.New("failed uniqueness constraint")
