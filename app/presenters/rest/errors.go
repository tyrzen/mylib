package rest

import "github.com/pkg/errors"

var ErrEncoding = errors.New("error encoding data to buffer")
var ErrDecoding = errors.New("error decoding data from request")
var ErrWritingResponse = errors.New("error writing response from buffer")
var ErrCreatingFile = errors.New("error creating file")
var ErrInvalidData = errors.New("invalid data, expected nil")
var ErrPermissions = errors.New("error permissions")
var ErrInvalidQuery = errors.New("invalid query")
