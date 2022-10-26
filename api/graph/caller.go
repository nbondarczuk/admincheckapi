package graph

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	MSGraphURL = "https://graph.microsoft.com/v1.0"
	ErrorHeader = "Error while getting group name:"
)

type Caller struct {
	Token string
	URL   string
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

//
// GroupName gets the name of the group by the group id
//
func (caller *Caller) GroupName(groupId string) (string, error) {
	URL := caller.URL + fmt.Sprintf("/groups/{%s}?$select=displayName", groupId)

	req, err := http.NewRequest("GET", URL, nil)
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
		var errResp ErrorResponse
		err = json.Unmarshal(body, &errResp)
		if err != nil {
			return "", fmt.Errorf("%s %s", ErrorHeader, err)
		}
		return "", fmt.Errorf("%s %s", ErrorHeader, errResp.Error.Code)
	}
	
	var response GroupNameResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("%s %s", ErrorHeader, err)
	}
	
	return response.DisplayName, nil
}

