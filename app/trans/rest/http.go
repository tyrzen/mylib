package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/delveper/mylib/app/ent"
)

type response struct {
	writer     http.ResponseWriter
	request    *http.Request
	Logger     ent.Logger
	statusCode int
	data       any
}

type Message struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func newResponse(rw http.ResponseWriter, req *http.Request, logger ent.Logger) response {
	return response{
		writer:  rw,
		request: req,
		Logger:  logger,
	}
}

func (resp *response) Set(code int, data any) {
	resp.statusCode = code
	resp.data = data
}

func (resp *response) Send() {
	resp.writer.Header().Set("Content-Type", "application/json")
	resp.writer.WriteHeader(resp.statusCode)

	if resp.data == nil {
		if resp.statusCode != http.StatusNoContent {
			resp.Logger.Errorf("Invalid data, expected nil")
		}

		return
	}

	if err := json.NewEncoder(resp.writer).Encode(resp.data); err != nil {
		resp.Logger.Errorf("Failed encoding to JSON %+v; with status code %d: %+v\n", resp.data, resp.statusCode, err)
	}
}

func decodeBody(req *http.Request, val any) (err error) {
	defer func() {
		if e := req.Body.Close(); e != nil {
			err = fmt.Errorf("error closing request body")
		}
	}()

	if err := json.NewDecoder(req.Body).Decode(val); err != nil {
		return fmt.Errorf("error decoding JSON body: %w", err)
	}

	return nil
}

func encodeBody(rw http.ResponseWriter, val any) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(val); err != nil {
		return fmt.Errorf("error decoding body to buffer: %w", err)
	}

	if _, err := buf.WriteTo(rw); err != nil {
		return fmt.Errorf("error writing rebody: %w", err)
	}

	return nil
}

func respond(rw http.ResponseWriter, req *http.Request, val interface{}, statusCode int, logger ent.Logger) {
	if val == nil && statusCode == http.StatusNoContent {
		logger.Errorf("Invalid data, expected nil")

		return
	}

	if err := encodeBody(rw, val); err != nil {
		logger.Errorw("Failed encoding to JSON.",
			"object", val,
			"status code", statusCode,
			"error", err,
		)
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(statusCode)

	if err := decodeBody(req, val); err != nil {
		logger.Errorw("Failed writing response.",
			"object", val,
			"status code", statusCode,
			"error", err,
		)
	}
}
