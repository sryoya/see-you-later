package syl

import "errors"

var errInvalidDuration = errors.New("invalid duration to open was provided")

var errUnsupportedOS = errors.New("unsupported platform")
