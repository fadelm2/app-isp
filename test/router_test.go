package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang-clean-architecture/internal/model"

	"github.com/stretchr/testify/assert"
)

func TestCreateRouterSuccess(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	requestBody := model.CreateRouterRequest{
		Name:     "Core Router",
		Host:     "192.168.1.1",
		Port:     8728,
		Username: "admin",
		Password: "password",
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/admin/routers", bytes.NewReader(bodyJson))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.RouterResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
	assert.Equal(t, requestBody.Host, responseBody.Data.Host)
}

func TestCreateRouterFailedValidation(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	// Missing Name and Host is invalid
	requestBody := model.CreateRouterRequest{
		Name:     "",
		Host:     "",
		Port:     -10,
		Username: "",
		Password: "",
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/admin/routers", bytes.NewReader(bodyJson))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestListRoutersSuccess(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	db.Exec("INSERT INTO routers (id, name, host, port, username, password, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"router-1", "Mikrotik 1", "10.0.0.1", 8728, "admin", "pass", "offline", 1783899129160, 1783899129160)

	request := httptest.NewRequest(http.MethodGet, "/api/admin/routers", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[[]model.RouterResponse])
	json.Unmarshal(bytes, responseBody)
	assert.True(t, len(responseBody.Data) > 0)
}

func TestGetRouterSuccess(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	db.Exec("INSERT INTO routers (id, name, host, port, username, password, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"router-1", "Mikrotik 1", "10.0.0.1", 8728, "admin", "pass", "offline", 1783899129160, 1783899129160)

	request := httptest.NewRequest(http.MethodGet, "/api/admin/routers/router-1", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[model.RouterResponse])
	json.Unmarshal(bytes, responseBody)
	assert.Equal(t, "router-1", responseBody.Data.ID)
}

func TestGetRouterFailedNotFound(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	request := httptest.NewRequest(http.MethodGet, "/api/admin/routers/router-not-exist", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestDeleteRouterSuccess(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	db.Exec("INSERT INTO routers (id, name, host, port, username, password, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"router-1", "Mikrotik 1", "10.0.0.1", 8728, "admin", "pass", "offline", 1783899129160, 1783899129160)

	request := httptest.NewRequest(http.MethodDelete, "/api/admin/routers/router-1", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestDeleteRouterFailedNotFound(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	request := httptest.NewRequest(http.MethodDelete, "/api/admin/routers/router-not-exist", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}
