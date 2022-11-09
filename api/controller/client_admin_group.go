package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"admincheckapi/api/config"
	"admincheckapi/api/repository"
	"admincheckapi/api/resource"
	"admincheckapi/api/stat"	
)

//
// CheckClientGroupAdmin reads a list of mappings between client
// and an admin group from local db
//
func CheckClientGroupAdmin(w http.ResponseWriter, r *http.Request) {
	log.Traceln("Begin: CheckClientGroupAdmin")

	w.Header().Set("X-Request-Id", stat.RequestId())
	
	log.Debugf("Handling request [%s] %s %s %s",
		r.Method,
		r.Host,
		r.URL.Path,
		r.URL.RawQuery)

	//
	// Parse path variables: client, group
	//
	
	client, err := pathVariableStr(r, "client", true)
	if err != nil {
		displayAppError(w, UrlPathError,
			"Missing mandatory url path variable client",
			http.StatusInternalServerError)
		return
	}
	log.Debugln("Got path variable client = " + client)

	group, err := pathVariableStr(r, "group", true)
	if err != nil {
		displayAppError(w, UrlPathError,
			"Missing mandatory url path variable group",
			http.StatusInternalServerError)
		return
	}
	log.Debugln("Got path variable group = " + group)

	//
	// Hit the backend storage
	//
	
	repo, err := repository.NewClientAdminGroupRepository(config.Setup.UsedBackend)
	if err != nil {
		displayAppError(w, RepositoryNewError,
			"Error while creating repository - "+err.Error(),
			http.StatusInternalServerError)
		return
	}
	defer repo.Close()

	log.Debugln("Client: " + client)
	log.Debugln("Group: " + group)
	count, err := repo.CountClientGroups(client, group)
	if err != nil {
		displayAppError(w, RepositoryRunError,
			"Error in repository read - "+err.Error(),
			http.StatusInternalServerError)
		return
	}
	log.Debugln("Found client groups count = " + fmt.Sprintf("%d", count))

	//
	// Found any groups?
	//
	
	var reply = resource.ClientGroupAdminReplyResource{
		Status: true,
		Data: resource.ClientGroupAdmin{
			Admin: count > 0,
		},
	}

	jstr, err := json.Marshal(&reply)
	if err != nil {
		displayAppError(w, EncoderJsonError,
			"An error while marshalling data - "+err.Error(),
			http.StatusInternalServerError)
		return
	}
	
	log.Debugln("Reply: " + string(jstr))
	writeResponseWithJson(w, http.StatusOK, jstr)

	log.Traceln("End: CheckClientGroupAdmin")
}

//
// ReadClientAdminGroup returns mapping of client to groups
//
func ReadClientAdminGroups(w http.ResponseWriter, r *http.Request) {
	log.Traceln("Begin: ReadClientAdminGroups")

	w.Header().Set("X-Request-Id", stat.RequestId())
	
	log.Debugf("Handling request [%s] %s %s %s",
		r.Method,
		r.Host,
		r.URL.Path,
		r.URL.RawQuery)

	//
	// Parse path variables: client, group
	//
	
	client, err := pathVariableStr(r, "client", true)
	if err != nil {
		displayAppError(w, UrlPathError,
			"Missing mandatory url path variable client",
			http.StatusInternalServerError)
		return
	}
	log.Debugln("Got path variable client = " + client)

	//
	// The the backend storage
	//
	
	repo, err := repository.NewClientAdminGroupRepository(config.Setup.UsedBackend)
	if err != nil {
		displayAppError(w, RepositoryNewError,
			"Error while creating repository - "+err.Error(),
			http.StatusInternalServerError)
		return
	}
	defer repo.Close()

	log.Debugln("Client: " + client)
	entries, count, err := repo.ReadClientGroups(client)
	if err != nil {
		displayAppError(w, RepositoryRunError,
			"Error in repository read - "+err.Error(),
			http.StatusInternalServerError)
		return
	}
	log.Debugln("Found entries count = " + fmt.Sprintf("%d", count))

	//
	// Give feedback about the operation results
	//
	
	var reply = resource.ClientAdminGroupReplyResource{
		Status: true,
		Data: resource.ClientAdminGroups{
			Count: count,
			Data:  entries,
		},
	}

	jstr, err := json.Marshal(&reply)
	if err != nil {
		displayAppError(w, EncoderJsonError,
			"An error while marshalling data - "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	log.Debugln("Reply: " + string(jstr))
	writeResponseWithJson(w, http.StatusOK, jstr)

	log.Traceln("End: ReadClientAdminGroups")
}

//
// CreateClientAdminGroup creates a mapping for client to a group
// returning what was created
//
func CreateClientAdminGroup(w http.ResponseWriter, r *http.Request) {
	log.Traceln("Begin: CreateClientAdminGroup")

	w.Header().Set("X-Request-Id", stat.RequestId())
	
	log.Debugf("Handling request [%s] %s %s %s",
		r.Method,
		r.Host,
		r.URL.Path,
		r.URL.RawQuery)

	//
	// Parse path variables: client, group
	//

	client, err := pathVariableStr(r, "client", true)
	if err != nil {
		displayAppError(w, UrlPathError,
			"Missing mandatory url path variable client",
			http.StatusInternalServerError)
		return
	}
	log.Debugln("Got path variable client = " + client)

	group, err := pathVariableStr(r, "group", true)
	if err != nil {
		displayAppError(w, UrlPathError,
			"Missing mandatory url path variable group",
			http.StatusInternalServerError)
		return
	}
	log.Debugln("Got path variable group = " + group)

	//
	// Hit the backend storage
	//
	
	rb, err := repository.NewClientAdminGroupRepository(config.Setup.UsedBackend)
	if err != nil {
		displayAppError(w, RepositoryNewError,
			"Error while creating repository - "+err.Error(),
			http.StatusInternalServerError)
		return
	}
	defer rb.Close()

	log.Debugln("Client: " + client)
	log.Debugln("Group: " + group)
	entries, count, err := rb.CreateClientGroup(client, group)
	if err != nil {
		displayAppError(w, RepositoryRunError,
			"Error in repository read - "+err.Error(),
			http.StatusInternalServerError)
		return
	}
	log.Debugln("Found entries count = " + fmt.Sprintf("%d", count))

	//
	// Give feedback about the operation results
	//
		
	var reply = resource.ClientAdminGroupReplyResource{
		Status: true,
		Data: resource.ClientAdminGroups{
			Count: count,
			Data:  entries,
		},
	}

	jstr, err := json.Marshal(&reply)
	if err != nil {
		displayAppError(w, EncoderJsonError,
			"An error while marshalling data - "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	log.Debugln("Reply: " + string(jstr))
	writeResponseWithJson(w, http.StatusOK, jstr)

	log.Traceln("End: CreateClientAdminGroup")
}

//
// DeleteClientAdminGroup deletes mapping of client to a group
// and returns what was deleted
//
func DeleteClientAdminGroup(w http.ResponseWriter, r *http.Request) {
	log.Traceln("Begin: DeleteClientAdminGroup")

	w.Header().Set("X-Request-Id", stat.RequestId())
	
	log.Debugf("Handling request [%s] %s %s %s",
		r.Method,
		r.Host,
		r.URL.Path,
		r.URL.RawQuery)

	//
	// Parse path variables: client, group
	//
	
	client, err := pathVariableStr(r, "client", true)
	if err != nil {
		displayAppError(w, UrlPathError,
			"Missing mandatory url path variable client",
			http.StatusInternalServerError)
		return
	}
	log.Debugln("Got path variable client = " + client)

	group, err := pathVariableStr(r, "group", true)
	if err != nil {
		displayAppError(w, UrlPathError,
			"Missing mandatory url path variable group",
			http.StatusInternalServerError)
		return
	}
	log.Debugln("Got path variable group = " + group)

	//
	// Hit the storage via repository access
	//
	
	rb, err := repository.NewClientAdminGroupRepository(config.Setup.UsedBackend)
	if err != nil {
		displayAppError(w, RepositoryNewError,
			"Error while creating repository - "+err.Error(),
			http.StatusInternalServerError)
		return
	}
	defer rb.Close()

	log.Debugln("Client: " + client)
	log.Debugln("Group: " + group)
	entries, count, err := rb.DeleteClientGroup(client, group)
	if err != nil {
		displayAppError(w, RepositoryRunError,
			"Error in repository read - "+err.Error(),
			http.StatusInternalServerError)
		return
	}
	log.Debugln("Found entries count = " + fmt.Sprintf("%d", count))

	//
	// Give feedback about opertation
	//
	
	var reply = resource.ClientAdminGroupReplyResource{
		Status: true,
		Data: resource.ClientAdminGroups{
			Count: count,
			Data:  entries,
		},
	}

	jstr, err := json.Marshal(&reply)
	if err != nil {
		displayAppError(w, EncoderJsonError,
			"An error while marshalling data - "+err.Error(),
			http.StatusInternalServerError)
		return
	}
	
	log.Debugln("Reply: " + string(jstr))
	writeResponseWithJson(w, http.StatusOK, jstr)

	log.Traceln("End: DeleteClientAdminGroup")
}

//
// PurgeClientAdminGroup deletes mapping of client to a group and returns what was deleted
//
func PurgeClientAdminGroups(w http.ResponseWriter, r *http.Request) {
	log.Traceln("Begin: PurgeClientAdminGroups")

	w.Header().Set("X-Request-Id", stat.RequestId())
	
	log.Debugf("Handling request [%s] %s %s %s",
		r.Method,
		r.Host,
		r.URL.Path,
		r.URL.RawQuery)

	//
	// Purge whole table (testing mode!)
	//
	
	rb, err := repository.NewClientAdminGroupRepository(config.Setup.UsedBackend)
	if err != nil {
		displayAppError(w, RepositoryNewError,
			"Error while creating repository - "+err.Error(),
			http.StatusInternalServerError)
		return
	}
	defer rb.Close()

	err = rb.PurgeClientGroups()
	if err != nil {
		displayAppError(w, RepositoryRunError,
			"Error in repository read - "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	//
	// Operation status returned
	//
	
	var reply = resource.ClientAdminGroupReplyResource{
		Status: true,
	}

	jstr, err := json.Marshal(&reply)
	if err != nil {
		displayAppError(w, EncoderJsonError,
			"An error while marshalling data - "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	log.Debugln("Reply: " + string(jstr))
	writeResponseWithJson(w, http.StatusOK, jstr)

	log.Traceln("End: PurgeClientAdminGroups")
}
