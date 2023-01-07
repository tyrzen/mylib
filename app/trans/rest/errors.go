package rest

import "github.com/pkg/errors"

var ErrEncoding = errors.New("error encoding data to buffer")
var ErrDecoding = errors.New("error decoding data from request")
var ErrWritingResponse = errors.New("error writing response from buffer")
var ErrInvalidData = errors.New("invalid data, expected nil")
