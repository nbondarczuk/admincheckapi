package controller

import (
	"encoding/json"
	"net/http"
    "regexp"

	log "github.com/sirupsen/logrus"

	"admincheckapi/api/backend"
	"admincheckapi/api/config"
	"admincheckapi/api/repository"
	"admincheckapi/api/repository/azure"
	"admincheckapi/api/resource"
	"admincheckapi/api/secretstore"
	"admincheckapi/api/stat"
	"admincheckapi/api/token"
)

//
// CheckClientAdminToken reads token from the payload and checks if it is
// an admin group of the client
//
func CheckClientAdminToken(w http.ResponseWriter, r *http.Request) {
	log.Traceln("Begin: CheckClientAdminToken")

	w.Header().Set("X-Request-Id", stat.RequestId())

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
	log.Debugf("Got from client token group ids: (%d) %v", len(ids), ids)

	var (
		found bool
		cache = 0
	)

	//
	// First hit the inmem cache
	//

	var ri repository.ClientAdminGroupRepository

	if !found && len(ids) > 0 {
		cache = 1
		log.Debugf("Search inmem cache for groups: (%d) %v", len(ids), ids)

		ri, err = repository.NewClientAdminGroupRepository("inmem")
		if err != nil {
			displayAppError(w, RepositoryNewError,
				"Error while creating repository - "+err.Error(),
				http.StatusInternalServerError)
			return
		}
		defer ri.Close()

		// each id of the request token
		for i, id := range ids {
			log.Debugf("Search inmem cache with group id: %s, round: %d", id, i)

			count, err := ri.CountClientGroups(client, id)
			if err != nil {
				displayAppError(w, RepositoryRunError,
					"Error in repository read - "+err.Error(),
					http.StatusInternalServerError)
				return
			}

			if count > 0 {
				log.Debugf("Found admin group id in in the inmem cache: %s", id)
				found = true
				break
			}
		}
	}

	if found && cache == 1 && len(ids) > 0 {
		log.Debugf("Found one of group in the inmem cache: %v", ids)
	} else if !found && cache == 1 && len(ids) > 0 {
		log.Debugf("Inmem cache miss for group ids: %v", ids)
	}

	//
	// Next hit the DB cache if nothing found
	//

	var rb repository.ClientAdminGroupRepository

	if !found && len(ids) > 0 {
		cache = 2
		log.Debugf("Search DB cache for groups: (%d) %v", len(ids), ids)

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
		for i, id := range ids {
			log.Debugf("Search DB cache for group id: %s, round: %d", id, i)

			count, err := rb.CountClientGroups(client, id)
			if err != nil {
				displayAppError(w, RepositoryRunError,
					"Error in repository read - "+err.Error(),
					http.StatusInternalServerError)
				return
			}

			// Stop searching if some entries found
			if count > 0 {
				log.Debugf("Found admin group id in in DB cache: %s", id)
				found = true
				break
			}
		}
	}

	if found && cache == 2 && len(ids) > 0 {
		log.Debugf("Found one of groups in DB cache: %s", ids)
	} else if !found && cache == 2 && len(ids) > 0 {
		log.Debugf("DB cache miss for group id: %v", ids)
	}

	//
	// Next hit the MS graph using own tenant JWT token if nothing found
	//

	if !found && len(ids) > 0 {
		log.Debugf("Search MS graph for groups: (%d) %v", len(ids), ids)

		clientTenantId, err := t.TenantId()
		if err != nil {
			displayAppError(w, PayloadReadError,
				"Unable to read group from token of the request",
				http.StatusInternalServerError)
			return
		}
		log.Debugf("Got from token client tenent id: %s", clientTenantId)

		// TenantJWTToken provides jwt token in the client context for MS graph hit
		clientContextAppToken, err := secretstore.TenantJWTToken(clientTenantId)
		if err != nil {
			displayAppError(w, RepositoryNewError,
				"Error while accessing secret store for token - "+err.Error(),
				http.StatusInternalServerError)
			return
		}
		log.Debugf("Got client context token: %s", clientContextAppToken)

		//
		// The token used to access Azure is the token in the client's
		// context, which is obtained using the credentials stored in the
		// secret store.
		//
		ba, err := backend.NewBackend("azure:" + clientContextAppToken)
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
		log.Debugf("Connected to Azure with client context token")

		var adminGroupId string
		
		if !config.Setup.UseGroupNamePattern {
			log.Debugf("Accessing graph with specific group name: %s", config.Setup.AdminGroupName)
			
			//
			// Get admin id of the client and search the list. Expectation is
			// that the list is short and there is only one group defined as admin.
			//			
			
			adminGroupId, err = ra.ClientGroupId(config.Setup.AdminGroupName)
			if err != nil {
				displayAppError(w, RepositoryRunError,
					"Error in Azure repository read - "+err.Error(),
					http.StatusInternalServerError)
				return
			}

			// In the list of token group ids?
			for _, id := range ids {
				if id == adminGroupId {
					found = true
					break
				}
			}
			
		} else {
			//
			// Map every id to name: slower but more coherent with regexp match
			//

			log.Debugf("Accessing graph with group name pattern: %s", config.Setup.AdminGroupName)
			
			var name string

			// each id of the request token
			for i, id := range ids {
				log.Debugf("Search MS graph for group name with id: %s round: %d", id, i)

				name, err = ra.ClientGroupName(id)
				if err != nil {
					displayAppError(w, RepositoryRunError,
						"Error in Azure repository read - "+err.Error(),
						http.StatusInternalServerError)
					return
				}
				log.Debugf("Found in MS graph group name: %s <- id: %s", name, id)
				
				// Is it admin group name?
				match, _ := regexp.MatchString(config.Setup.AdminGroupName, name)
				log.Debugf("Check for admin group name match: %s with group name %s -> %t",
					config.Setup.AdminGroupName, name, match)
				if match {
					log.Debugf("Found admin group in MS graph: %s <- %s", name, id)
					found = true
					break
				} else {
					log.Debugf("Found not an admin group in MS graph: %s <- %s", name, id)
				}
			}
		}

		if found {
			// add client -> id to inmem cache firt a quick storage
			if ri != nil {
				_, _, err = ri.CreateClientGroup(client, adminGroupId)
				if err != nil {
					displayAppError(w, RepositoryRunError,
						"Error in repository create - "+err.Error(),
						http.StatusInternalServerError)
					return
				}
				log.Debugf("Populated inmem cache with: client: %s groupid: %s", client, adminGroupId)
			}
			
			// add client -> id to DB cache as slow storage
			if rb != nil {
						_, _, err = rb.CreateClientGroup(client, adminGroupId)
				if err != nil {
					displayAppError(w, RepositoryRunError,
						"Error in repository write - "+err.Error(),
						http.StatusInternalServerError)
					return
				}
				log.Debugf("Populated DB cache with: client: %s groupid: %s", client, adminGroupId)
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

	jstr, err := json.Marshal(reply)
	if err != nil {
		displayAppError(w, EncoderJsonError,
			"An error while marshalling data - "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	log.Debugln("Reply: " + string(jstr))
	writeResponseWithJson(w, http.StatusOK, jstr)

	log.Traceln("End: CheckClientAdminToken")
}
