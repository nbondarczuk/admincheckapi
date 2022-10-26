package controller

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"admincheckapi/api/backend"
	"admincheckapi/api/config"
	"admincheckapi/api/repository"
	"admincheckapi/api/repository/azure"
	"admincheckapi/api/resource"
	"admincheckapi/api/secretstore"
	"admincheckapi/api/token"
)

//
// CheckClientAdminToken reads token from the payload and checks if it is
// an admin group of the client
//
func CheckClientAdminToken(w http.ResponseWriter, r *http.Request) {
	log.Traceln("Begin: CheckClientAdminToken")

	log.Debugf("Handling request [%s] %s %s %s",
		r.Method,
		r.Host,
		r.URL.Path,
		r.URL.RawQuery)

	//
	// Parse path variable: client
	//
	
	client, err := pathVariableStr(r, "client", true)
	if err != nil {
		displayAppError(w, UrlPathError,
			"Missing mandatory url path variable client",
			http.StatusInternalServerError)
		return
	}
	log.Debugln("Got path variable client: " + client)

	//
	// Read payload with JWT token to be checked
	//
	
	payload, err := readPayload(r)
	if err != nil {
		displayAppError(w, PayloadReadError,
			"Unable to read payload of the request",
			http.StatusInternalServerError)
		return
	}
	log.Debugln("Got payload: " + string(payload))

	var request resource.ClientTokenRequestResource
	err = json.Unmarshal(payload, &request)
	if err != nil {
		displayAppError(w, PayloadReadError,
			"Unable to read payload of the request",
			http.StatusInternalServerError)
		return
	}

	//
	// Validate JWT token and get all group ids from claims
	//
	
	t, err := token.NewToken([]byte(request.Token))
	if err != nil {
		displayAppError(w, PayloadReadError,
			"Unable to parse the token from the request",
			http.StatusInternalServerError)
		return
	}
	log.Debugln("Validated client token")

	ids, err := t.AdminGroups()
	if err != nil {
		displayAppError(w, PayloadReadError,
			"Unable to read group from token of the request",
			http.StatusInternalServerError)
		return
	}
	log.Debugf("Got from token group ids: (%d) %v", len(ids), ids)

	//	
	// First hit the inmem cache
	//

	var found bool	
	log.Debugf("Search inmem cache")
	
	ri, err := repository.NewClientAdminGroupRepository("inmem")
	if err != nil {
		displayAppError(w, RepositoryNewError,
			"Error while creating repository - "+err.Error(),
			http.StatusInternalServerError)
		return
	}
	defer ri.Close()

	// each id of the request token
	for _, id := range ids {
		log.Debugf("Search inmem cache with group id: %s", id)
		
		count, err := ri.CountClientGroups(client, id)
		if err != nil {
			displayAppError(w, RepositoryRunError,
				"Error in repository read - "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		if count > 0 {
			log.Debugf("Found group in in inmem cache")
			found = true
			break
		}
	}

	//
	// Next hit the DB cache if nothing found
	//
	
	var rb repository.ClientAdminGroupRepository

	if !found {
		log.Debugf("Search DB cache")
		
		var err error
		rb, err = repository.NewClientAdminGroupRepository(config.Setup.UsedBackend)
		if err != nil {
			displayAppError(w, RepositoryNewError,
				"Error while creating repository - "+err.Error(),
				http.StatusInternalServerError)
			return
		}
		defer rb.Close()

		// each id of the request token
		for _, id := range ids {
			log.Debugf("Search DB cache for group id: %s", id)
			
			count, err := rb.CountClientGroups(client, id)
			if err != nil {
				displayAppError(w, RepositoryRunError,
					"Error in repository read - "+err.Error(),
					http.StatusInternalServerError)
				return
			}

			// Stop searching if some entries found
			if count > 0 {
				log.Debugf("Found group in in DB cache")
				found = true
				break
			}
		}
	}

	//
	// Next hit the MS graph using own tenant JWT token if nothing found
	//
	
	if !found {
		log.Debugf("Search graph")
		
		apptoken, err := secretstore.TenantJWTToken()
		if err != nil {
			displayAppError(w, RepositoryNewError,
				"Error while accessing secret store for token - "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		ba, err := backend.NewBackend("azure:" + apptoken)
		if err != nil {
			displayAppError(w, RepositoryNewError,
				"Error while creating Azure backend - "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		ra, err := azure.NewClientAdminGroupRepository(ba)
		if err != nil {
			displayAppError(w, RepositoryNewError,
				"Error while creating Azure repository - "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		var name string

		// each id of the request token
		for _, id := range ids {
			log.Debugf("Search graph with group id: %s", id)
			
			name, err = ra.ClientGroupName(id)
			if err != nil {
				displayAppError(w, RepositoryRunError,
					"Error in Azure repository read - "+err.Error(),
					http.StatusInternalServerError)
				return
			}

			// Is it admin group name?
			if name == config.Setup.AdminGroupName {
				log.Debugf("Found group name in graph: %s %s", id, name)
				found = true

				// add client -> id to inmem cache firt a quick storage
				_, _, err = ri.CreateClientGroup(client, id)
				if err != nil {
					displayAppError(w, RepositoryRunError,
						"Error in repository create - "+err.Error(),
						http.StatusInternalServerError)
					return
				}
				log.Debugf("Populated inmem cache with: client: %s groupid: %s", client, id)
				
				// add client -> id to DB cache as slow storage
				_, _, err = rb.CreateClientGroup(client, id)
				if err != nil {
					displayAppError(w, RepositoryRunError,
						"Error in repository write - "+err.Error(),
						http.StatusInternalServerError)
					return
				}
				log.Debugf("Populated DB cache with: client: %s groupid: %s", client, id)

				// Stop searching if the qualified group name found
				break
			} 
		}
	}

	//
	// Found admin group in JWT token?
	//
	
	var reply = resource.ClientGroupAdminReplyResource{
		Status: true,
		Data: resource.ClientGroupAdmin{
			Admin: found,
		},
	}

	if jstr, err := json.Marshal(reply); err != nil {
		displayAppError(w, EncoderJsonError,
			"An error while marshalling data - "+err.Error(),
			http.StatusInternalServerError)
		return
	} else {
		log.Debugln("Reply: " + string(jstr))
		writeResponseWithJson(w, http.StatusOK, jstr)
	}

	log.Traceln("End: CheckClientAdminToken")
}
