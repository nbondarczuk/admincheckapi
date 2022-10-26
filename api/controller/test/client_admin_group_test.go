package controller_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"admincheckapi/api/controller"
	"admincheckapi/api/resource"
	"admincheckapi/test/testconfig"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func resetData(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/client/purge", nil)
	w := httptest.NewRecorder()
	routerForPurgeClientAdminGroup().ServeHTTP(w, req)
}

func prolog(t *testing.T) {
	ok := os.Getenv("POSTGRES")
	if ok == "" {
		t.Skip("Postgres DB not available, skip")
	}

	ok = os.Getenv("MSAD")
	if ok == "" {
		t.Skip("MS AD not available, skip")
	}

	testconfig.Set(t)
	resetData(t)
}

func epilog(t *testing.T) {
	resetData(t)
}

func routerForCheckClientGroupAdmin() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/client/{client}/group/{group}/admin", controller.CheckClientGroupAdmin)
	return r
}

func routerForReadClientAdminGroup() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/client/{client}/admin/group", controller.ReadClientAdminGroups)
	return r
}

func routerForCreateClientAdminGroup() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/client/{client}/admin/group/{group}", controller.CreateClientAdminGroup)
	return r
}

func routerForDeleteClientAdminGroup() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/client/{client}/admin/group/{group}", controller.DeleteClientAdminGroup)
	return r
}

func routerForPurgeClientAdminGroup() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/client/purge", controller.PurgeClientAdminGroups)
	return r
}

func TestCheckClientGroupAdminNotExist(t *testing.T) {
	prolog(t)

	t.Run("api request for not existing client group", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/client/CLIENT1/group/whatever/admin", nil)
		w := httptest.NewRecorder()
		routerForCheckClientGroupAdmin().ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Error from request: %v", err)
		}

		var reply resource.ClientGroupAdminReplyResource
		err = json.Unmarshal(data, &reply)
		if err != nil {
			t.Errorf("Error unmarshalling response from request: %v", err)
		}

		assert.Equal(t, true, reply.Status)
		assert.Equal(t, false, reply.Data.Admin)
	})

	epilog(t)
}

func TestCheckClientGroupAdminExist(t *testing.T) {
	prolog(t)

	t.Run("api request for existing client group", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/client/CLIENT2/admin/group/ADMINGROUP1", nil)
		w := httptest.NewRecorder()
		routerForCreateClientAdminGroup().ServeHTTP(w, req)

		req = httptest.NewRequest(http.MethodGet, "/api/client/CLIENT2/group/ADMINGROUP1/admin", nil)
		w = httptest.NewRecorder()
		routerForCheckClientGroupAdmin().ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Error from request: %v", err)
		}

		var reply resource.ClientGroupAdminReplyResource
		err = json.Unmarshal(data, &reply)
		if err != nil {
			t.Errorf("Error unmarshalling response from request: %v", err)
		}

		assert.Equal(t, true, reply.Status)
		assert.Equal(t, true, reply.Data.Admin)
	})

	epilog(t)
}

func TestReadClientAdminGroupNotExist(t *testing.T) {
	prolog(t)

	t.Run("api request for reading missing client and group", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/client/CLIENT3/admin/group", nil)
		w := httptest.NewRecorder()
		routerForReadClientAdminGroup().ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Error from request: %v", err)
		}

		var reply resource.ClientAdminGroupReplyResource
		err = json.Unmarshal(data, &reply)
		if err != nil {
			t.Errorf("Error unmarshalling response from request: %v", err)
		}

		assert.Equal(t, true, reply.Status)
		assert.Equal(t, int64(0), reply.Data.Count)
	})

	epilog(t)
}

func TestReadClientAdminGroupOne(t *testing.T) {
	prolog(t)

	t.Run("api request for reading created client and group", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/client/CLIENT4/admin/group/ADMINGROUP", nil)
		w := httptest.NewRecorder()
		routerForCreateClientAdminGroup().ServeHTTP(w, req)

		req = httptest.NewRequest(http.MethodGet, "/api/client/CLIENT4/admin/group", nil)
		w = httptest.NewRecorder()
		routerForReadClientAdminGroup().ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Error from request: %v", err)
		}

		var reply resource.ClientAdminGroupReplyResource
		err = json.Unmarshal(data, &reply)
		if err != nil {
			t.Errorf("Error unmarshalling response from request: %v", err)
		}

		assert.Equal(t, true, reply.Status)
		assert.Equal(t, int64(1), reply.Data.Count)
		assert.Equal(t, 1, len(reply.Data.Data))
		assert.Equal(t, "CLIENT4", reply.Data.Data[0].Client)
		assert.Equal(t, "ADMINGROUP", reply.Data.Data[0].AdminGroupId)
	})

	epilog(t)
}

func TestCreateClientAdminGroupOne(t *testing.T) {
	prolog(t)

	t.Run("api request for creating a new client admin group", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/client/CLIENT5/admin/group/ADMINGROUP2", nil)
		w := httptest.NewRecorder()
		routerForCreateClientAdminGroup().ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Error reading response data from request: %v", err)
		}

		var reply resource.ClientAdminGroupReplyResource
		err = json.Unmarshal(data, &reply)
		if err != nil {
			t.Errorf("Error unmarshalling response from request: %v", err)
		}

		assert.Equal(t, true, reply.Status)
		assert.Equal(t, true, reply.Status)
		assert.Equal(t, int64(1), reply.Data.Count)
		assert.Equal(t, 1, len(reply.Data.Data))
		assert.Equal(t, "CLIENT5", reply.Data.Data[0].Client)
		assert.Equal(t, "ADMINGROUP2", reply.Data.Data[0].AdminGroupId)
	})

	epilog(t)
}

func TestDeleteClientAdminGroupotExist(t *testing.T) {
	prolog(t)

	t.Run("api request for deleting not existent client admin group", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/client/CLIENT6/admin/group/whatever", nil)
		w := httptest.NewRecorder()
		routerForDeleteClientAdminGroup().ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Error from request: %v", err)
		}

		var reply resource.ClientAdminGroupReplyResource
		err = json.Unmarshal(data, &reply)
		if err != nil {
			t.Errorf("Error unmarshalling response from request: %v", err)
		}

		assert.Equal(t, true, reply.Status)
		assert.Equal(t, int64(0), reply.Data.Count)
		assert.Equal(t, 0, len(reply.Data.Data))
	})

	epilog(t)
}

func TestDeleteClientAdminGroupOne(t *testing.T) {
	prolog(t)

	t.Run("api request for deleting one created client admin group", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/client/CLIENT7/admin/group/ADMINGROUP", nil)
		w := httptest.NewRecorder()
		routerForCreateClientAdminGroup().ServeHTTP(w, req)

		req = httptest.NewRequest(http.MethodDelete, "/api/client/CLIENT7/admin/group/ADMINGROUP", nil)
		w = httptest.NewRecorder()
		routerForDeleteClientAdminGroup().ServeHTTP(w, req)

		req = httptest.NewRequest(http.MethodGet, "/api/client/CLIENT7/admin/group", nil)
		w = httptest.NewRecorder()
		routerForReadClientAdminGroup().ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Error from request: %v", err)
		}

		var reply resource.ClientAdminGroupReplyResource
		err = json.Unmarshal(data, &reply)
		if err != nil {
			t.Errorf("Error unmarshalling response from request: %v", err)
		}

		assert.Equal(t, true, reply.Status)
		assert.Equal(t, int64(0), reply.Data.Count)
		assert.Equal(t, 0, len(reply.Data.Data))
	})

	epilog(t)
}

func TestDeleteClientAdminGroupOneofTwo(t *testing.T) {
	prolog(t)

	t.Run("api request for deleting one of two created client admin groups", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/client/CLIENT8/admin/group/ADMINGROUP1", nil)
		w := httptest.NewRecorder()
		routerForCreateClientAdminGroup().ServeHTTP(w, req)

		req = httptest.NewRequest(http.MethodPost, "/api/client/CLIENT8/admin/group/ADMINGROUP2", nil)
		w = httptest.NewRecorder()
		routerForCreateClientAdminGroup().ServeHTTP(w, req)

		req = httptest.NewRequest(http.MethodDelete, "/api/client/CLIENT8/admin/group/ADMINGROUP1", nil)
		w = httptest.NewRecorder()
		routerForDeleteClientAdminGroup().ServeHTTP(w, req)

		req = httptest.NewRequest(http.MethodGet, "/api/client/CLIENT8/admin/group", nil)
		w = httptest.NewRecorder()
		routerForReadClientAdminGroup().ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Error from request: %v", err)
		}

		var reply resource.ClientAdminGroupReplyResource
		err = json.Unmarshal(data, &reply)
		if err != nil {
			t.Errorf("Error unmarshalling response from request: %v", err)
		}

		assert.Equal(t, true, reply.Status)
		assert.Equal(t, int64(1), reply.Data.Count)
		assert.Equal(t, 1, len(reply.Data.Data))
		assert.Equal(t, "CLIENT8", reply.Data.Data[0].Client)
		assert.Equal(t, "ADMINGROUP2", reply.Data.Data[0].AdminGroupId)
	})

	epilog(t)
}
