package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang-clean-architecture/internal/model"

	"github.com/stretchr/testify/assert"
)

func TestCreatePackageSuccess(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	requestBody := model.CreatePackageRequest{
		Name:            "Family Super 30M",
		SpeedMbps:       30,
		Price:           200000,
		InstallationFee: 150000,
		TaxRate:         0.11,
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/admin/packages", bytes.NewReader(bodyJson))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.PackageResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
	assert.Equal(t, requestBody.SpeedMbps, responseBody.Data.SpeedMbps)
	assert.Equal(t, requestBody.Price, responseBody.Data.Price)
}

func TestCreatePackageFailedValidation(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	// Name is empty (failed validation)
	requestBody := model.CreatePackageRequest{
		Name:      "",
		SpeedMbps: 0,
		Price:     -100,
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/admin/packages", bytes.NewReader(bodyJson))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestCreatePackageFailedUnauthorized(t *testing.T) {
	ClearAll()

	requestBody := model.CreatePackageRequest{
		Name:      "Unauthorized Package",
		SpeedMbps: 10,
		Price:     100000,
	}
	bodyJson, _ := json.Marshal(requestBody)
	request := httptest.NewRequest(http.MethodPost, "/api/admin/packages", bytes.NewReader(bodyJson))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer wrong_token")

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestListPackagesSuccess(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	// Create a package first
	db.Exec("INSERT INTO internet_packages (id, name, speed_mbps, price, installation_fee, tax_rate, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"pkg-1", "Test Package", 10, 100000, 50000, 0.11, true, 1783899129160, 1783899129160)

	// Admin list packages
	request := httptest.NewRequest(http.MethodGet, "/api/admin/packages", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[[]model.PackageResponse])
	json.Unmarshal(bytes, responseBody)
	assert.True(t, len(responseBody.Data) > 0)
}

func TestGetPackageSuccess(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	db.Exec("INSERT INTO internet_packages (id, name, speed_mbps, price, installation_fee, tax_rate, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"pkg-1", "Test Package", 10, 100000, 50000, 0.11, true, 1783899129160, 1783899129160)

	request := httptest.NewRequest(http.MethodGet, "/api/admin/packages/pkg-1", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[model.PackageResponse])
	json.Unmarshal(bytes, responseBody)
	assert.Equal(t, "pkg-1", responseBody.Data.ID)
	assert.Equal(t, "Test Package", responseBody.Data.Name)
}

func TestGetPackageFailedNotFound(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	request := httptest.NewRequest(http.MethodGet, "/api/admin/packages/non-existent-id", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestUpdatePackageSuccess(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	db.Exec("INSERT INTO internet_packages (id, name, speed_mbps, price, installation_fee, tax_rate, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"pkg-1", "Test Package", 10, 100000, 50000, 0.11, true, 1783899129160, 1783899129160)

	isActiveVal := false
	updateReq := model.UpdatePackageRequest{
		ID:       "pkg-1",
		Name:     "Test Package Updated",
		Price:    120000,
		IsActive: &isActiveVal,
	}
	bodyJson, _ := json.Marshal(updateReq)
	request := httptest.NewRequest(http.MethodPatch, "/api/admin/packages/pkg-1", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[model.PackageResponse])
	json.Unmarshal(bytes, responseBody)
	assert.Equal(t, "Test Package Updated", responseBody.Data.Name)
	assert.Equal(t, float64(120000), responseBody.Data.Price)
}

func TestUpdatePackageFailedNotFound(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	updateReq := model.UpdatePackageRequest{
		ID:    "non-existent",
		Name:  "Test Package Updated",
		Price: 120000,
	}
	bodyJson, _ := json.Marshal(updateReq)
	request := httptest.NewRequest(http.MethodPatch, "/api/admin/packages/non-existent", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestDeletePackageSuccess(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	db.Exec("INSERT INTO internet_packages (id, name, speed_mbps, price, installation_fee, tax_rate, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"pkg-1", "Test Package", 10, 100000, 50000, 0.11, true, 1783899129160, 1783899129160)

	request := httptest.NewRequest(http.MethodDelete, "/api/admin/packages/pkg-1", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestDeletePackageFailedNotFound(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	request := httptest.NewRequest(http.MethodDelete, "/api/admin/packages/non-existent", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}
