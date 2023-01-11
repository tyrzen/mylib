package exc

import "errors"

var ErrDuplicateEmail = errors.New("email is already taken")
var ErrDuplicateID = errors.New("id already exists")
var ErrNoRecord = errors.New("record not found")
var ErrInvalidCredits = errors.New("invalid credentials")
var ErrUnexpected = errors.New("unexpected error")
var ErrDeadline = errors.New("deadline exceeded")
var ErrHashing = errors.New("error hashing")
