package controller

import (
	"encoding/json"
	"net/http"

	"admincheckapi/api/resource"
	"admincheckapi/api/stat"
	"admincheckapi/api/version"
)

//
// ReadSystemHealth responds with system health feedback  for K8S
//
func ReadSystemHealth(w http.ResponseWriter, r *http.Request) {
	writeResponseWithJson(w, http.StatusOK, nil)
}

//
// ReadSystemAlive responds with alivenes feedback for K8S
//
func ReadSystemAlive(w http.ResponseWriter, r *http.Request) {
	writeResponseWithJson(w, http.StatusOK, nil)
}

//
// ReadSystemStat responds with status info about processing status info
//
func ReadSystemStat(w http.ResponseWriter, r *http.Request) {
	dataReplyResource := resource.StatResource{
		Status: true,
		Data: stat.Info(),
	}

	if j, err := json.Marshal(dataReplyResource); err != nil {
		displayAppError(w, err,
			"Error json encoding version info",
			http.StatusInternalServerError)
		return
	} else {
		writeResponseWithJson(w, http.StatusOK, j)
	}	
}

//
// ReadSystemVersion responds with version info
//
func ReadSystemVersion(w http.ResponseWriter, r *http.Request) {
	dataReplyResource := resource.VersionResource{
		Status: true,
		Data: version.Level(),
	}

	if j, err := json.Marshal(dataReplyResource); err != nil {
		displayAppError(w, err,
			"Error json encoding version info",
			http.StatusInternalServerError)
		return
	} else {
		writeResponseWithJson(w, http.StatusOK, j)
	}
}
