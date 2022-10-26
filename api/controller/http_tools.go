package controller

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

//
// writeResponseJson makes a successful response with status and payload
//
func writeResponseWithJson(w http.ResponseWriter, status int, payload []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if payload != nil {
		w.Write(payload)
	}
}

//
// pathVariableStr gets & validates existence of string parameter
//
func pathVariableStr(r *http.Request, label string, mandatory bool) (value string, err error) {
	vars := mux.Vars(r)
	value, ok := vars[label]
	if !ok {
		if mandatory {
			err = fmt.Errorf("Manadatory variable does not exist: " + label)
		}
	}

	return value, err
}

//
// readPayload reads contents of the payload from the request body
//
func readPayload(r *http.Request) (body []byte, err error) {
	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("Can't read request body")
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return
}
