package exc

import "errors"

var ErrDuplicateEmail = errors.New("email is already taken")
var ErrDuplicateID = errors.New("id is already taken")
var ErrUnexpected = errors.New("unexpected error")
