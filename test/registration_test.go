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

func TestCreateRegistrationSuccess(t *testing.T) {
	ClearAll()

	// Seed package first
	db.Exec("INSERT INTO internet_packages (id, name, speed_mbps, price, installation_fee, tax_rate, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"pkg-1", "Test Package", 10, 100000, 50000, 0.11, true, 1783899129160, 1783899129160)

	requestBody := model.CreateRegistrationRequest{
		FullName:            "Agus Salim",
		NIK:                 "1234567890123456",
		BirthPlace:          "Bandung",
		BirthDate:           "1995-05-15",
		Gender:              "Laki-laki",
		Email:               "agus@example.com",
		Phone:               "08123456789",
		InstallationAddress: "Jl. Sukabumi No. 10",
		BillingAddress:      "Jl. Sukabumi No. 10",
		PackageID:           "pkg-1",
		Latitude:            -6.2,
		Longitude:           106.81,
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/registrations", bytes.NewReader(bodyJson))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.RegistrationResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, requestBody.FullName, responseBody.Data.FullName)
	assert.Equal(t, "pending", responseBody.Data.Status)
}

func TestCreateRegistrationFailedValidation(t *testing.T) {
	ClearAll()

	// Invalid email and missing required NIK/FullName
	requestBody := model.CreateRegistrationRequest{
		FullName:  "",
		NIK:       "",
		Email:     "invalid-email",
		PackageID: "pkg-1",
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/registrations", bytes.NewReader(bodyJson))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestListRegistrationsSuccess(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	db.Exec("INSERT INTO internet_packages (id, name, speed_mbps, price, installation_fee, tax_rate, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"pkg-1", "Test Package", 10, 100000, 50000, 0.11, true, 1783899129160, 1783899129160)
	db.Exec("INSERT INTO registrations (id, full_name, nik, birth_place, birth_date, gender, email, phone, installation_address, billing_address, package_id, latitude, longitude, status, ktp_path, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"reg-1", "Agus Salim", "1234567890123456", "Bandung", "1995-05-15", "Laki-laki", "agus@example.com", "08123456789", "Jl. Sukabumi", "Jl. Sukabumi", "pkg-1", -6.2, 106.81, "pending", "/storage/ktp.jpg", 1783899129160, 1783899129160)

	request := httptest.NewRequest(http.MethodGet, "/api/admin/registrations", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[[]model.RegistrationResponse])
	json.Unmarshal(bytes, responseBody)
	assert.True(t, len(responseBody.Data) > 0)
}

func TestListRegistrationsFailedUnauthorized(t *testing.T) {
	ClearAll()

	request := httptest.NewRequest(http.MethodGet, "/api/admin/registrations", nil)
	request.Header.Set("Authorization", "Bearer invalid-token")

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestUpdateRegistrationStatusApproved(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	db.Exec("INSERT INTO internet_packages (id, name, speed_mbps, price, installation_fee, tax_rate, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"pkg-1", "Test Package", 10, 100000, 50000, 0.11, true, 1783899129160, 1783899129160)
	db.Exec("INSERT INTO registrations (id, full_name, nik, birth_place, birth_date, gender, email, phone, installation_address, billing_address, package_id, latitude, longitude, status, ktp_path, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"reg-1", "Agus Salim", "1234567890123456", "Bandung", "1995-05-15", "Laki-laki", "agus@example.com", "08123456789", "Jl. Sukabumi", "Jl. Sukabumi", "pkg-1", -6.2, 106.81, "pending", "/storage/ktp.jpg", 1783899129160, 1783899129160)

	updateReq := model.UpdateRegistrationStatusRequest{
		ID:     "reg-1",
		Status: "approved",
	}
	bodyJson, _ := json.Marshal(updateReq)

	request := httptest.NewRequest(http.MethodPatch, "/api/admin/registrations/reg-1/status", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[model.RegistrationResponse])
	json.Unmarshal(bytes, responseBody)
	assert.Equal(t, "approved", responseBody.Data.Status)
}

func TestUpdateRegistrationStatusFailedNotFound(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	updateReq := model.UpdateRegistrationStatusRequest{
		ID:     "non-existent",
		Status: "approved",
	}
	bodyJson, _ := json.Marshal(updateReq)

	request := httptest.NewRequest(http.MethodPatch, "/api/admin/registrations/non-existent/status", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestUpdateRegistrationStatusFailedValidation(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	// Status is empty
	updateReq := model.UpdateRegistrationStatusRequest{
		ID:     "reg-1",
		Status: "",
	}
	bodyJson, _ := json.Marshal(updateReq)

	request := httptest.NewRequest(http.MethodPatch, "/api/admin/registrations/reg-1/status", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}
