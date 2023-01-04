package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/delveper/mylib/app/ent"
)

type Response struct {
	Message string `json:"message"`
	Details string `json:"details"`
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
		logger.Errorw("Failed writing response from buffer.",
			"object", val,
			"status code", statusCode,
			"error", err,
		)
	}
}
