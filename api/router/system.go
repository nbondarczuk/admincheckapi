package router

import (
	"github.com/gorilla/mux"

	"admincheckapi/api/controller"
)

//
// NewSystemRouter creates router for system methods
//
func NewSystemRouter(r *mux.Router) *mux.Router {
	r.HandleFunc("/system/health",
		controller.ReadSystemHealth).
		Methods("GET").
		Name("read-system-health")

	r.HandleFunc("/system/alive",
		controller.ReadSystemAlive).
		Methods("GET").
		Name("read-system-alive")

	r.HandleFunc("/system/stat",
		controller.ReadSystemStat).
		Methods("GET").
		Name("read-system-stat")

	r.HandleFunc("/system/version",
		controller.ReadSystemVersion).
		Methods("GET").
		Name("read-system-version")

	return r
}
