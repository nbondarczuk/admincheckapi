package controller

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"admincheckapi/api/auth"
	"admincheckapi/api/resource"
	"admincheckapi/api/stat"
)

//
// CheckClientAdminAuth gets a token, try to refresh the cache and check
//
func CheckClientAdminAuth(w http.ResponseWriter, r *http.Request) {
	log.Traceln("Begin: CheckClientAdminAuth")

	w.Header().Set("X-Request-Id", stat.RequestId())
	
	log.Debugf("Handling request [%s] %s %s %s",
		r.Method,
		r.Host,
		r.URL.Path,
		r.URL.RawQuery)

	//
	// Parse request args: client, method
	//
	
	client, err := pathVariableStr(r, "client", true)
	if err != nil {
		displayAppError(w, UrlPathError,
			"Missing mandatory url path variable client",
			http.StatusInternalServerError)
		return
	}
	log.Debugln("Got path variable client = " + client)

	method, err := pathVariableStr(r, "method", true)
	if err != nil {
		displayAppError(w, UrlPathError,
			"Missing mandatory url path variable method",
			http.StatusInternalServerError)
		return
	}
	log.Debugln("Got path variable method = " + method)

	if !isValidMethod(method) {
		displayAppError(w, UrlPathError,
			"Invalid method: " + method,
			http.StatusInternalServerError)
		return
	}

	//
	// Get payload with credentials
	//
	
	payload, err := readPayload(r)
	if err != nil {
		displayAppError(w, PayloadReadError,
			"Unable to read payload of the request",
			http.StatusInternalServerError)
		return
	}
	log.Debugln("Got payload = " + string(payload))

	var request resource.ClientAdminAuthRequestResource
	err = json.Unmarshal(payload, &request)
	if err != nil {
		displayAppError(w, PayloadReadError,
			"Unable to decode json payload of the request",
			http.StatusInternalServerError)
		return
	}
	log.Debugf("Got request: %+v", request)
	
	//
	// Hit auth provider with auth request in order to obtain JWT token
	//
	
	// structural equivalence of external type and internal one: same fields
	var claim auth.Claim
	claim.ClientID = request.ClientID
	claim.Authority = request.Authority
	claim.Scopes = request.Scopes
	claim.ClientSecret = request.ClientSecret
	log.Debugf("Using claim: %+v", claim)
	
	am, err := auth.NewAuthMethod(method, claim)
	if err != nil {
		displayAppError(w, AuthError,
			"Unable to authorise",
			http.StatusInternalServerError)
		return
	}

	//
	// Give feedback with result: JWT token
	//
	
	var reply = resource.ClientAdminAuthReplyResource{
		Status: true,
		Token: am.Token(),
	}

	jstr, err := json.Marshal(reply)
	if err != nil {
		displayAppError(w, EncoderJsonError,
			"An error while marshalling data - "+err.Error(),
			http.StatusInternalServerError)
		return
	}
	
	log.Debugln("Reply: " + string(jstr))
	writeResponseWithJson(w, http.StatusOK, jstr)

	log.Traceln("End: CheckClientAdminAuth")
}

// isValidMethod checks the method value
func isValidMethod(method string) bool {
	switch method {
	case "secret":
		return true
	case "certificate":
		return true
	case "userpassword":
		return true
	case "code":
		return true
	}

	return false
}
