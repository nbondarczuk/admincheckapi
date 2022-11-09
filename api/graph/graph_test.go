package graph_test

import (
	"admincheckapi/api/graph"
	"admincheckapi/api/model"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGroupName(t *testing.T) {
	testResponse := `{
		"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#groups(displayName)/$entity",
		"displayName": "MyGroup1"
	  }`
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, testResponse)
	}))
	defer svr.Close()
	caller := graph.Caller{
		Token: "xyz",
		URL:   svr.URL,
	}
	name, err := caller.GroupName("id-example")
	require.Empty(t, err)
	require.Equal(t, name, "MyGroup1")
}

func TestGetGroups(t *testing.T) {
	testResponse := `{
		"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#groups(id,displayName)",
		"value": [
		  {
			"id": "1",
			"displayName": "MyGroup1"
		  },
		  {
			"id": "2",
			"displayName": "MyGroup2"
		  }
		]
	  }`
	expected := []model.ClientAdminGroup{
		{
			Client:     "1",
			AdminGroup: "MyGroup1",
		},
		{
			Client:     "2",
			AdminGroup: "MyGroup2",
		},
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, testResponse)
	}))
	defer svr.Close()
	caller := graph.Caller{
		Token: "xyz",
		URL:   svr.URL,
	}
	groups, err := caller.GetGroups()
	require.Empty(t, err)
	require.Equal(t, len(expected), len(groups))
	for i, group := range groups {
		require.Equal(t, group.Client, expected[i].Client)
		require.Equal(t, group.AdminGroup, expected[i].AdminGroup)
	}

}
