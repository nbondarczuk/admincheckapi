package router

import (
	"github.com/gorilla/mux"

	"admincheckapi/api/controller"
)

//
// NewClientAdminRouter creates the router for client API
//
func NewClientAdminRouter(r *mux.Router) *mux.Router {
	r.HandleFunc("/api/client/{client:[A-Za-z0-9]+}/admin/token",
		controller.CheckClientAdminToken).
		Methods("POST").
		Name("CheckClientAdminAuthToken")

	r.HandleFunc("/api/client/{client:[A-Za-z0-9]+}/group/{group}/admin",
		controller.CheckClientGroupAdmin).
		Methods("GET").
		Name("CheckClientGroupAdmin")

	r.HandleFunc("/api/client/{client:[A-Za-z0-9]+}/admin/group",
		controller.ReadClientAdminGroups).
		Methods("GET").
		Name("ReadClientAdminGroup")

	r.HandleFunc("/api/client/{client:[A-Za-z0-9]+}/admin/group/{group}",
		controller.CreateClientAdminGroup).
		Methods("POST").
		Name("CreateClientAdminGroup")

	r.HandleFunc("/api/client/{client:[A-Za-z0-9]+}/admin/group/{group}",
		controller.DeleteClientAdminGroup).
		Methods("DELETE").
		Name("DeleteClientAdminGroup")

	r.HandleFunc("/api/client/purge",
		controller.PurgeClientAdminGroups).
		Methods("POST").
		Name("PurgeClientAdminGroups")

	r.HandleFunc("/api/client/{client:[A-Za-z0-9]+}/admin/auth/{method}",
		controller.CheckClientAdminAuth).
		Methods("POST").
		Name("CheckClientAdminAuthWithCode")

	return r
}
