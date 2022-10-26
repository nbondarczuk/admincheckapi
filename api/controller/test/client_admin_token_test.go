package controller_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"admincheckapi/api/controller"
	"admincheckapi/api/resource"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func routerForCheckClientAdminToken() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/client/{client}/admin/token", controller.CheckClientAdminToken)
	return r
}

func TestCheckClientAdminTokenMissingGroup(t *testing.T) {
	prolog(t)
	
	t.Run("check token for missing group", func(t *testing.T) {
		resetData(t)

		var payload resource.ClientTokenRequestResource
		payload.Token = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IjJaUXBKM1VwYmpBWVhZR2FYRUpsOGxWMFRPSSIsImtpZCI6IjJaUXBKM1VwYmpBWVhZR2FYRUpsOGxWMFRPSSJ9.eyJhdWQiOiJhcGk6Ly8wNDJkODA3ZC1hMThiLTQ5NjUtYTgyNC1lZmY3Mzg1NjA3ZTYiLCJpc3MiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC8yM2Q0MjUwZS0zNzExLTQ4ZTEtOGRmZS1kMjdmMjg1MTU0YTEvIiwiaWF0IjoxNjYyOTg1ODA4LCJuYmYiOjE2NjI5ODU4MDgsImV4cCI6MTY2Mjk4OTcwOCwiYWlvIjoiRTJaZ1lPajlmU3lCTlNpeTRldEYyOHhma1RKdkFBPT0iLCJhcHBpZCI6IjUzNTllNmEzLWRjMWUtNGFjMS04MDczLWQwZTY1OGEwMDJjNCIsImFwcGlkYWNyIjoiMSIsImdyb3VwcyI6WyIyZTEyZDFiYi1jNDhlLTQyY2QtYjRjZC01NzMzODVlNGE5ZGMiXSwiaWRwIjoiaHR0cHM6Ly9zdHMud2luZG93cy5uZXQvMjNkNDI1MGUtMzcxMS00OGUxLThkZmUtZDI3ZjI4NTE1NGExLyIsIm9pZCI6ImM2MDNhZmIwLWNjODUtNDYzMy1hYWM3LWI0YWJjODQzZjIwMyIsInJoIjoiMC5BWUlBRGlYVUl4RTM0VWlOX3RKX0tGRlVvWDJBTFFTTG9XVkpxQ1R2OXpoV0ItYVZBQUEuIiwicm9sZXMiOlsiYXJnb25BZG1pblJvbGUxIiwidGVzdFJvbGUxIl0sInN1YiI6ImM2MDNhZmIwLWNjODUtNDYzMy1hYWM3LWI0YWJjODQzZjIwMyIsInRpZCI6IjIzZDQyNTBlLTM3MTEtNDhlMS04ZGZlLWQyN2YyODUxNTRhMSIsInV0aSI6InU0X1ludC12SjAtVnk2Tk1kWFF4QUEiLCJ2ZXIiOiIxLjAifQ.ZfixXtPMrCEEUrVoW8rUQe-Wqx08nztZ1omqSRmQfWt8dimYndyJ5f4jUkAIfzKvyPtgvpXJo4dvWVpAeNUiPoLR3dcFjrMz9b_EkfXFI3NS0hkSqwFytazHa3v_o_V7TpIi5XBJRgfba5pYJlDeqnEPhEEScOD_jhSTjBDVJWp2j6iEsxmOJ5KKCzyG4FX0laJO16lcGcBsUzkVvIIH2n5FZXyrqMcL2Evgko6d64VAtx4wA5Kxanvs6igS3bIV7MUjpeP67aTcAafVCR20hgWT_IpZ6qaBidr68H6dzyeHTMYfeE-2NvO5o8mrNg65jndIDGA84zyFkxcaqh4ghQ"
		var body []byte
		if jstr, err := json.Marshal(payload); err != nil {
			t.Errorf("Error from request: %v", err)
		} else {
			body = []byte(jstr)
		}
		req := httptest.NewRequest(http.MethodPost, "/api/client/ARGON/admin/token", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		routerForCheckClientAdminToken().ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Error from request: %v", err)
		}
		var reply resource.ClientGroupAdminReplyResource
		err = json.Unmarshal(data, &reply)
		if err != nil {
			t.Errorf("Error unmarshalling response from request: %s - %s", err, data)
		}

		assert.Equal(t, true, reply.Status)
		assert.Equal(t, false, reply.Data.Admin)
	})

	epilog(t)
}

func TestCheckClientAdminTokenWrongGroup(t *testing.T) {
	prolog(t)

	t.Run("check token for wrong group", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/client/ARGON/admin/group/whatever", nil)
		w := httptest.NewRecorder()
		routerForCreateClientAdminGroup().ServeHTTP(w, req)

		var payload resource.ClientTokenRequestResource
		payload.Token = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IjJaUXBKM1VwYmpBWVhZR2FYRUpsOGxWMFRPSSIsImtpZCI6IjJaUXBKM1VwYmpBWVhZR2FYRUpsOGxWMFRPSSJ9.eyJhdWQiOiJhcGk6Ly8wNDJkODA3ZC1hMThiLTQ5NjUtYTgyNC1lZmY3Mzg1NjA3ZTYiLCJpc3MiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC8yM2Q0MjUwZS0zNzExLTQ4ZTEtOGRmZS1kMjdmMjg1MTU0YTEvIiwiaWF0IjoxNjYyOTg1ODA4LCJuYmYiOjE2NjI5ODU4MDgsImV4cCI6MTY2Mjk4OTcwOCwiYWlvIjoiRTJaZ1lPajlmU3lCTlNpeTRldEYyOHhma1RKdkFBPT0iLCJhcHBpZCI6IjUzNTllNmEzLWRjMWUtNGFjMS04MDczLWQwZTY1OGEwMDJjNCIsImFwcGlkYWNyIjoiMSIsImdyb3VwcyI6WyIyZTEyZDFiYi1jNDhlLTQyY2QtYjRjZC01NzMzODVlNGE5ZGMiXSwiaWRwIjoiaHR0cHM6Ly9zdHMud2luZG93cy5uZXQvMjNkNDI1MGUtMzcxMS00OGUxLThkZmUtZDI3ZjI4NTE1NGExLyIsIm9pZCI6ImM2MDNhZmIwLWNjODUtNDYzMy1hYWM3LWI0YWJjODQzZjIwMyIsInJoIjoiMC5BWUlBRGlYVUl4RTM0VWlOX3RKX0tGRlVvWDJBTFFTTG9XVkpxQ1R2OXpoV0ItYVZBQUEuIiwicm9sZXMiOlsiYXJnb25BZG1pblJvbGUxIiwidGVzdFJvbGUxIl0sInN1YiI6ImM2MDNhZmIwLWNjODUtNDYzMy1hYWM3LWI0YWJjODQzZjIwMyIsInRpZCI6IjIzZDQyNTBlLTM3MTEtNDhlMS04ZGZlLWQyN2YyODUxNTRhMSIsInV0aSI6InU0X1ludC12SjAtVnk2Tk1kWFF4QUEiLCJ2ZXIiOiIxLjAifQ.ZfixXtPMrCEEUrVoW8rUQe-Wqx08nztZ1omqSRmQfWt8dimYndyJ5f4jUkAIfzKvyPtgvpXJo4dvWVpAeNUiPoLR3dcFjrMz9b_EkfXFI3NS0hkSqwFytazHa3v_o_V7TpIi5XBJRgfba5pYJlDeqnEPhEEScOD_jhSTjBDVJWp2j6iEsxmOJ5KKCzyG4FX0laJO16lcGcBsUzkVvIIH2n5FZXyrqMcL2Evgko6d64VAtx4wA5Kxanvs6igS3bIV7MUjpeP67aTcAafVCR20hgWT_IpZ6qaBidr68H6dzyeHTMYfeE-2NvO5o8mrNg65jndIDGA84zyFkxcaqh4ghQ"
		var body []byte
		if jstr, err := json.Marshal(payload); err != nil {
			t.Errorf("Error from request: %v", err)
		} else {
			body = []byte(jstr)
		}
		req = httptest.NewRequest(http.MethodPost, "/api/client/ARGON/admin/token", bytes.NewBuffer(body))
		w = httptest.NewRecorder()

		routerForCheckClientAdminToken().ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Error from request: %v", err)
		}
		var reply resource.ClientGroupAdminReplyResource
		err = json.Unmarshal(data, &reply)
		if err != nil {
			t.Errorf("Error unmarshalling response from request: %s - %s", err, data)
		}

		assert.Equal(t, true, reply.Status)
		assert.Equal(t, false, reply.Data.Admin)
	})

	epilog(t)
}

func TestCheckClientAdminToken(t *testing.T) {
	prolog(t)

	t.Run("check token for existing group", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/client/ARGON/admin/group/2e12d1bb-c48e-42cd-b4cd-573385e4a9dc", nil)
		w := httptest.NewRecorder()
		routerForCreateClientAdminGroup().ServeHTTP(w, req)

		var payload resource.ClientTokenRequestResource
		payload.Token = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IjJaUXBKM1VwYmpBWVhZR2FYRUpsOGxWMFRPSSIsImtpZCI6IjJaUXBKM1VwYmpBWVhZR2FYRUpsOGxWMFRPSSJ9.eyJhdWQiOiJhcGk6Ly8wNDJkODA3ZC1hMThiLTQ5NjUtYTgyNC1lZmY3Mzg1NjA3ZTYiLCJpc3MiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC8yM2Q0MjUwZS0zNzExLTQ4ZTEtOGRmZS1kMjdmMjg1MTU0YTEvIiwiaWF0IjoxNjYyOTg1ODA4LCJuYmYiOjE2NjI5ODU4MDgsImV4cCI6MTY2Mjk4OTcwOCwiYWlvIjoiRTJaZ1lPajlmU3lCTlNpeTRldEYyOHhma1RKdkFBPT0iLCJhcHBpZCI6IjUzNTllNmEzLWRjMWUtNGFjMS04MDczLWQwZTY1OGEwMDJjNCIsImFwcGlkYWNyIjoiMSIsImdyb3VwcyI6WyIyZTEyZDFiYi1jNDhlLTQyY2QtYjRjZC01NzMzODVlNGE5ZGMiXSwiaWRwIjoiaHR0cHM6Ly9zdHMud2luZG93cy5uZXQvMjNkNDI1MGUtMzcxMS00OGUxLThkZmUtZDI3ZjI4NTE1NGExLyIsIm9pZCI6ImM2MDNhZmIwLWNjODUtNDYzMy1hYWM3LWI0YWJjODQzZjIwMyIsInJoIjoiMC5BWUlBRGlYVUl4RTM0VWlOX3RKX0tGRlVvWDJBTFFTTG9XVkpxQ1R2OXpoV0ItYVZBQUEuIiwicm9sZXMiOlsiYXJnb25BZG1pblJvbGUxIiwidGVzdFJvbGUxIl0sInN1YiI6ImM2MDNhZmIwLWNjODUtNDYzMy1hYWM3LWI0YWJjODQzZjIwMyIsInRpZCI6IjIzZDQyNTBlLTM3MTEtNDhlMS04ZGZlLWQyN2YyODUxNTRhMSIsInV0aSI6InU0X1ludC12SjAtVnk2Tk1kWFF4QUEiLCJ2ZXIiOiIxLjAifQ.ZfixXtPMrCEEUrVoW8rUQe-Wqx08nztZ1omqSRmQfWt8dimYndyJ5f4jUkAIfzKvyPtgvpXJo4dvWVpAeNUiPoLR3dcFjrMz9b_EkfXFI3NS0hkSqwFytazHa3v_o_V7TpIi5XBJRgfba5pYJlDeqnEPhEEScOD_jhSTjBDVJWp2j6iEsxmOJ5KKCzyG4FX0laJO16lcGcBsUzkVvIIH2n5FZXyrqMcL2Evgko6d64VAtx4wA5Kxanvs6igS3bIV7MUjpeP67aTcAafVCR20hgWT_IpZ6qaBidr68H6dzyeHTMYfeE-2NvO5o8mrNg65jndIDGA84zyFkxcaqh4ghQ"
		var body []byte
		if jstr, err := json.Marshal(payload); err != nil {
			t.Errorf("Error from request: %v", err)
		} else {
			body = []byte(jstr)
		}
		req = httptest.NewRequest(http.MethodPost, "/api/client/ARGON/admin/token", bytes.NewBuffer(body))
		w = httptest.NewRecorder()

		routerForCheckClientAdminToken().ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Error from request: %v", err)
		}
		var reply resource.ClientGroupAdminReplyResource
		err = json.Unmarshal(data, &reply)
		if err != nil {
			t.Errorf("Error unmarshalling response from request: %s - %s", err, data)
		}

		assert.Equal(t, true, reply.Status)
		assert.Equal(t, true, reply.Data.Admin)
	})

	epilog(t)
}
