package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/delveper/mylib/app/ent"
)

type Responder struct {
	responseWriter http.ResponseWriter
	request        *http.Request
	logger         ent.Logger
}

type Message struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (resp *Responder) withHandler(rw http.ResponseWriter, req *http.Request) {
	resp.responseWriter = rw
	resp.request = req
}

func (resp *Responder) writeResponse(statusCode int, data any) {
	if data == nil {
		if statusCode != http.StatusNoContent {
			resp.logger.Errorf("Invalid data, expected nil")
		}

		return
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		resp.logger.Errorw("Failed encoding to JSON.",
			"object", data,
			"error", err,
		)
	}

	resp.responseWriter.Header().Set("Content-Type", "application/json")
	resp.responseWriter.WriteHeader(statusCode)

	if _, err := buf.WriteTo(resp.responseWriter); err != nil {
		resp.logger.Errorw("Failed writing response from buffer.",
			"object", data,
			"error", err,
		)
	}
}

func (resp *Responder) decodeBody(data any) error {
	defer func() {
		if err := resp.request.Body.Close(); err != nil {
			resp.logger.Warnf("error closing request body")
		}
	}()

	if err := json.NewDecoder(resp.request.Body).Decode(data); err != nil {
		return fmt.Errorf("error decoding JSON body: %w", err)
	}

	return nil
}
