package router

import (
	"admincheckapi/api/config"

	"github.com/MadAppGang/httplog"
	"github.com/gorilla/mux"
)

//
// NewRouter creates a new router for API collecting all sub-routers. They handle
// route groups like client admin or system routes.
//
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r = NewSystemRouter(r)
	r = NewClientAdminRouter(r)
	
	if config.Setup.LogHTTP {
		r.Use(httplog.Logger)
	} else {
		r.Use(RequestLoggerMiddleware(r))
	}

	return r
}

