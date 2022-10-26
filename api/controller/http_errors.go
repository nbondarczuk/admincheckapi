package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type (
	AppError struct {
		Error      string `json:"error"`
		Message    string `json:"message"`
		HttpStatus int    `json:"status"`
	}

	ErrorResource struct {
		Data AppError `json:"data"`
	}
)

var (
	UrlPathError       = errors.New("URL path decoding error")
	DecoderJsonError   = errors.New("Decoder JSON error")
	EncoderJsonError   = errors.New("Encoder JSON error")
	RepositoryNewError = errors.New("Repository creation error")
	RepositoryRunError = errors.New("Repository runtime error")
	ControllerError    = errors.New("Controller error")
	PayloadReadError   = errors.New("Payload read error")
	AuthError          = errors.New("Authorisation error")
)

//
// displayAppError showing results in response json stored in header
//
func displayAppError(w http.ResponseWriter, handlerError error, message string, code int) {
	var info string
	if handlerError != nil {
		info = handlerError.Error()
	}

	var ae AppError = AppError{
		Error:      info,
		Message:    message,
		HttpStatus: code,
	}

	log.Errorf("Error: %d: %s", code, message)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if j, err := json.Marshal(ErrorResource{Data: ae}); err == nil {
		w.Write(j)
	}
}
