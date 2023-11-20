package service

import "errors"

// Validation error, meant to be combined with "errors.Join()" with
// whatever validation error that occurred.
var ErrValidation = errors.New("error validating input: ")
