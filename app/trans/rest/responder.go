package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/delveper/mylib/app/ent"
)

type responder struct{ ent.Logger }

type response struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (r responder) decodeBody(req *http.Request, data any) (err error) {
	defer func() {
		if e := req.Body.Close(); e != nil {
			r.Errorw("error while closing request body", "error", err)
		}
	}()

	if err := json.NewDecoder(req.Body).Decode(data); err != nil {
		return fmt.Errorf("error decoding JSON body: %w", err)
	}

	return nil
}

func (r responder) writeResponse(rw http.ResponseWriter, req *http.Request, code int, data any) {
	if data == nil && code != http.StatusNoContent {
		r.Errorw("Failed writing response due nil data.",
			"object", nil,
			"error", ErrInvalidData,
		)
		r.writeResponse(rw, req, http.StatusBadRequest, ErrInvalidData)

		return
	}

	if err, ok := data.(error); ok {
		data = response{Message: http.StatusText(code), Details: err.Error()}
	}

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(data)
	if err != nil {
		r.Errorw("Failed encoding data to JSON.",
			"object", data,
			"error", err)
		r.writeResponse(rw, req, http.StatusInternalServerError, ErrEncoding)

		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)

	if _, err := buf.WriteTo(rw); err != nil {
		r.Errorw("Failed writing response from buffer.",
			"object", data,
			"error", fmt.Errorf("%w: %v", ErrWritingResponse, err),
		)
	}
}
