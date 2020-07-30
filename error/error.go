package error

import "errors"

var(
	ErrReflectNil  = errors.New("can't reflect nil pointer")
	ErrReflectNonSlice = errors.New("can't reflect not slice object")
)