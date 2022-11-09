package graph

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	MSGraphURL         = "https://graph.microsoft.com/v1.0"
	LoginURL           = "https://login.microsoftonline.com/"
	ErrorHeader string = "Error while calling graph-api:"
)

type Caller struct {
	Token string
	URL   string
}

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    string `json:"expires_in"`
	ExtExpiresIn string `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}

type TokenErrorResponse struct {
	Error            string   `json:"error"`
	ErrorDescription string   `json:"error_description"`
	ErrorCode        []string `json:"error_codes"`
	Timestamp        string   `json:"timestamp"`
	TraceId          string   `json:"trace_id"`
	CorrelationId    string   `json:"correlation_id"`
}

type ErrorResponse struct {
	Error struct {
		Code       string `json:"code"`
		Message    string `json:"message"`
		InnerError struct {
			Date            string `json:"date"`
			RequestId       string `json:"request-id"`
			ClientRequestId string `json:"client-request-id"`
		}
	}
}

type GroupNameResponse struct {
	DataContext string `json:"@odata.context"`
	DisplayName string `json:"displayName"`
}

type GroupIdsResponse struct {
	DataContext string       `json:"@odata.context"`
	Value       []GroupValue `json:"value"`
}

type GroupValue struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
}

//
// GroupName maps the group id to a name
//
func (caller *Caller) GroupName(groupId string) (string, error) {
	log.Traceln("Begin: GroupName")
	defer log.Traceln("End: GroupName")

	URL := fmt.Sprintf("%s/groups/{%s}?$select=displayName", caller.URL, groupId)

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return "", fmt.Errorf("%s %s", ErrorHeader, err.Error())
	}
	req.Header.Add("Authorization", "Bearer "+caller.Token)

	client := &http.Client{}
	log.Debugf("MS graph request GET:%s", URL)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%s %s", ErrorHeader, err)
	}
	defer resp.Body.Close()
	log.Debugf("MS graph response: %v", resp)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%s %s", ErrorHeader, err)
	}
	if resp.StatusCode != 200 {
		var errResp ErrorResponse
		err = json.Unmarshal(body, &errResp)
		if err != nil {
			return "", fmt.Errorf("%s %s", ErrorHeader, err)
		}
		return "", fmt.Errorf("%s %v", ErrorHeader, errResp.Error)
	}

	var response GroupNameResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("Error while unmarshalling %s %s", ErrorHeader, err)
	}
	
	return response.DisplayName, nil
}

//
// GroupId maps the group name to an id, only one id may be returned due
// equality condition in search.
//
func (caller *Caller) GroupId(groupName string) (string, error) {
	log.Traceln("Begin: GroupId")
	defer log.Traceln("End: GroupId")

	// Prepare GET request with Azure endpoint as target
	URL := "%s/groups?$select=id,displayName&$filter=displayName%%20eq%%20%c%s%c"
	URL = fmt.Sprintf(URL, caller.URL, 0x27, groupName, 0x27)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return "", fmt.Errorf("%s: %s", "Error making new request", err.Error())
	}
	req.Header.Add("Authorization", "Bearer "+caller.Token)

	// Hit the endpoint
	client := &http.Client{}
	log.Debugf("MS graph request GET:%s", URL)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%s: %s", "Error doing request", err)
	}
	defer resp.Body.Close()
	log.Debugf("Result: %+v", resp)

	// Process results of the request
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%s: %s", "Error reading response", err)
	}
	if resp.StatusCode != 200 {
		var errResp ErrorResponse
		err = json.Unmarshal(body, &errResp)
		if err != nil {
			return "", fmt.Errorf("%s: %s", "Error parsing error response", err)
		}
		return "", fmt.Errorf("%s: %s", "Error from request", errResp.Error)
	}
  
	var response GroupIdsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("%s: %s", "Error decoding response body", err)
	}
	log.Debugf("Response: %+v", response)

	// Only one not empty entry may exist
	var n = len(response.Value)
	if n != 1 {
		return "", fmt.Errorf("%s %d", "Only one group expected, got groups no: ", n)
	}
	if response.Value[0].Id == "" {
		return "", fmt.Errorf("%s %d", "Empty group id value received")
	}

	return response.Value[0].Id, nil
}

func (caller *Caller) ClientToken(tenantId, clientId, clientSecret, scope string) (Token string, err error) {
	URL := LoginURL + fmt.Sprintf("/%s/oauth2/v2.0/token", tenantId)
	data := url.Values{}
	data.Add("grant_type", "client_credentials")
	data.Add("client_id", clientId)
	data.Add("client_secret", clientSecret)
	data.Add("scope", scope)
	req, err := http.NewRequest("POST", URL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("%s %s", ErrorHeader, err.Error())
	}
	req.Header.Add("Authorization", "Bearer "+caller.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%s %s", ErrorHeader, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%s %s", ErrorHeader, err)
	}
	if resp.StatusCode != 200 {
		var errResp TokenErrorResponse
		err = json.Unmarshal(body, &errResp)
		if err != nil {
			return "", fmt.Errorf("%s %s", ErrorHeader, err)
		}
		return "", fmt.Errorf("%s error: %s error description: %s", ErrorHeader, errResp.Error, errResp.ErrorDescription)
	}

	// Getting an array of values
	var response TokenResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("%s: %s", "Error parsing response", err)
	}

	return response.AccessToken, nil
}

func (caller *Caller) UserGroups(principal bool, oid string) ([]GroupValue, error) {
	urlTemplate := "https://graph.microsoft.com/v1.0/%s/%s/transitiveMemberOf?$select=id,displayName"
	if principal {
		urlTemplate = fmt.Sprintf(urlTemplate, "servicePrincipals", oid)
	} else {
		urlTemplate = fmt.Sprintf(urlTemplate, "users", oid)
	}
	req, err := http.NewRequest("GET", urlTemplate, nil)
	if err != nil {
		return []GroupValue{}, fmt.Errorf("%s: %s", "Error making new request", err.Error())
	}
	req.Header.Add("Authorization", "Bearer "+caller.Token)

	// Hit the endpoint
	client := &http.Client{}
	log.Debugf("MS graph request GET:%s", urlTemplate)
	resp, err := client.Do(req)
	if err != nil {
		return []GroupValue{}, fmt.Errorf("%s: %s", "Error doing request", err)
	}
	defer resp.Body.Close()
	log.Debugf("Result: %+v", resp)

	// Process results of the request
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []GroupValue{}, fmt.Errorf("%s: %s", "Error reading response", err)
	}
	if resp.StatusCode != 200 {
		var errResp ErrorResponse
		err = json.Unmarshal(body, &errResp)
		if err != nil {
			return []GroupValue{}, fmt.Errorf("%s: %s", "Error parsing error response", err)
		}
		return []GroupValue{}, fmt.Errorf("%s: %d", "Invalid status code", errResp.Error.Code)
	}

	// Getting an array of values
	var response GroupIdsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return []GroupValue{}, fmt.Errorf("%s: %s", "Error parsing response", err)
	}
	log.Debugf("Response: %+v", response)

	return response.Value, nil
}
