package controller

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
	
	"admincheckapi/api/resource"
	"admincheckapi/api/stat"
	"admincheckapi/api/version"
)

//
// ReadSystemHealth responds with system health feedback  for K8S
//
func ReadSystemHealth(w http.ResponseWriter, r *http.Request) {
	log.Traceln("Begin: ReadSystemHealth")

	log.Debugf("Handling request [%s] %s %s %s",
		r.Method,
		r.Host,
		r.URL.Path,
		r.URL.RawQuery)
	
	var status int
	if stat.IsHealthy() {
		status = http.StatusOK
	} else {
		status = http.StatusServiceUnavailable
	}

	writeResponseWithJson(w, status, nil)
	
	log.Traceln("End: ReadSystemHealth")
}

//
// ReadSystemAlive responds with alivenes feedback for K8S
//
func ReadSystemAlive(w http.ResponseWriter, r *http.Request) {
	log.Traceln("Begin: ReadSystemAlive")
	
	log.Debugf("Handling request [%s] %s %s %s",
		r.Method,
		r.Host,
		r.URL.Path,
		r.URL.RawQuery)
	
	var status int
	if stat.IsAlive() {
		status = http.StatusOK
	} else {
		status = http.StatusServiceUnavailable
	}

	writeResponseWithJson(w, status, nil)

	log.Traceln("End: ReadSystemAlive")
}

//
// ReadSystemStat responds with status info about processing status info
//
func ReadSystemStat(w http.ResponseWriter, r *http.Request) {
	log.Traceln("Begin: ReadSystemStat")
	
	log.Debugf("Handling request [%s] %s %s %s",
		r.Method,
		r.Host,
		r.URL.Path,
		r.URL.RawQuery)
	
	dataReplyResource := resource.StatResource{
		Status: true,
		Data:   stat.Info(),
	}

	jstr, err := json.Marshal(dataReplyResource)
	if err != nil {
		displayAppError(w, err,
			"Error json encoding version info",
			http.StatusInternalServerError)
		return
	} 

	writeResponseWithJson(w, http.StatusOK, jstr)

	log.Traceln("End: ReadSystemStat")
}

//
// ReadSystemVersion responds with version info
//
func ReadSystemVersion(w http.ResponseWriter, r *http.Request) {
	log.Traceln("Begin: ReadSystemVersion")
	
	log.Debugf("Handling request [%s] %s %s %s",
		r.Method,
		r.Host,
		r.URL.Path,
		r.URL.RawQuery)
	
	dataReplyResource := resource.VersionResource{
		Status: true,
		Data:   version.Level(),
	}

	jstr, err := json.Marshal(dataReplyResource)
	if err != nil {
		displayAppError(w, err,
			"Error json encoding version info",
			http.StatusInternalServerError)
		return
	}
	
	writeResponseWithJson(w, http.StatusOK, jstr)

	log.Traceln("End: ReadSystemVersion")
}
