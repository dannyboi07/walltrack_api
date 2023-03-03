package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"walltrack/schema"
)

func JsonParseErr(err error) (int, error) {
	if err == nil {
		return 0, nil
	}

	var syntaxError *json.SyntaxError
	var unmarshallTypeError *json.UnmarshalTypeError
	switch {
	case errors.As(err, &syntaxError):
		return http.StatusBadRequest, fmt.Errorf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)

	case errors.Is(err, io.ErrUnexpectedEOF):
		return http.StatusBadRequest, fmt.Errorf("Request body contains badly formed JSON")

	case errors.As(err, &unmarshallTypeError):
		return http.StatusBadRequest, fmt.Errorf("Request body contains an invalid value for field: %q, value: %q (at position: %d)", unmarshallTypeError.Field, unmarshallTypeError.Value, unmarshallTypeError.Offset)

	case strings.HasPrefix(err.Error(), "json: unknown field "):
		unknownFieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		return http.StatusBadRequest, fmt.Errorf("Request body contains unknown field: %s", unknownFieldName)

	case errors.Is(err, io.EOF):
		return http.StatusBadRequest, errors.New("Request body cannot be empty")

	case err.Error() == "http: request body too large":
		return http.StatusRequestEntityTooLarge, errors.New("Request body cannot be larger than 1MB")
	default:
		return http.StatusBadRequest, err
	}
}

// Passing `statusCode` as 0, sets `statusCode` as 200
func WriteApiMessage(w http.ResponseWriter, statusCode int, message string) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(schema.ApiMessage{
		Status:  "success",
		Message: message,
	})
}

func WriteApiErrMessage(w http.ResponseWriter, statusCode int, message string) {
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}
	if message == "" {
		message = "Sorry, something went wrong"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(schema.ApiMessage{
		Status:  "failed",
		Message: message,
	})
}

func WriteDataToResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schema.ApiData{
		Status: "success",
		Data:   data,
	})
}
