package exc

import "errors"

var ErrValidation = errors.New("validation error")
var ErrDuplicateEmail = errors.New("email is already taken")
var ErrDuplicateID = errors.New("id already exists")
var ErrNotAuthorized = errors.New("not authorized")
var ErrInvalidHeader = errors.New("invalid header")
var ErrNoRecord = errors.New("record not found")
var ErrInvalidCredits = errors.New("invalid credentials")
var ErrTokenInvalidSigningMethod = errors.New("invalid signing method")
var ErrTokenExpired = errors.New("token expired")
var ErrTokenInvalid = errors.New("token invalid")
var ErrTokenNotFound = errors.New("token not found")
var ErrUnexpected = errors.New("unexpected error")
var ErrDeadline = errors.New("deadline exceeded")
var ErrHashing = errors.New("error hashing")
var ErrComparingHash = errors.New("error comparing hash")
