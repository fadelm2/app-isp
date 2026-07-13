package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang-clean-architecture/internal/model"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestListCustomersSuccess(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	// Seed package
	db.Exec("INSERT INTO internet_packages (id, name, speed_mbps, price, installation_fee, tax_rate, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"pkg-1", "Test Package", 10, 100000, 50000, 0.11, true, 1783899129160, 1783899129160)

	// Seed user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	db.Exec("INSERT INTO users (id, role_id, username, email, password, company_name, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		"cust-u-1", 2, "Agus", "agus@example.com", string(hashedPassword), "GREENET", 1783899129160, 1783899129160)

	// Seed customer
	db.Exec("INSERT INTO customers (id, user_id, status, package_id, due_date_day, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		"CUST-1", "cust-u-1", "active", "pkg-1", 5, 1783899129160, 1783899129160)

	request := httptest.NewRequest(http.MethodGet, "/api/admin/customers", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[[]model.CustomerResponse])
	json.Unmarshal(bytes, responseBody)
	assert.True(t, len(responseBody.Data) > 0)
}

func TestGetCustomerSuccess(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	db.Exec("INSERT INTO internet_packages (id, name, speed_mbps, price, installation_fee, tax_rate, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"pkg-1", "Test Package", 10, 100000, 50000, 0.11, true, 1783899129160, 1783899129160)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	db.Exec("INSERT INTO users (id, role_id, username, email, password, company_name, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		"cust-u-1", 2, "Agus", "agus@example.com", string(hashedPassword), "GREENET", 1783899129160, 1783899129160)
	db.Exec("INSERT INTO customers (id, user_id, status, package_id, due_date_day, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		"CUST-1", "cust-u-1", "active", "pkg-1", 5, 1783899129160, 1783899129160)

	request := httptest.NewRequest(http.MethodGet, "/api/admin/customers/CUST-1", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[model.CustomerResponse])
	json.Unmarshal(bytes, responseBody)
	assert.Equal(t, "CUST-1", responseBody.Data.ID)
}

func TestGetCustomerFailedNotFound(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	request := httptest.NewRequest(http.MethodGet, "/api/admin/customers/CUST-NOT-EXIST", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestSuspendCustomerSuccess(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	db.Exec("INSERT INTO internet_packages (id, name, speed_mbps, price, installation_fee, tax_rate, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"pkg-1", "Test Package", 10, 100000, 50000, 0.11, true, 1783899129160, 1783899129160)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	db.Exec("INSERT INTO users (id, role_id, username, email, password, company_name, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		"cust-u-1", 2, "Agus", "agus@example.com", string(hashedPassword), "GREENET", 1783899129160, 1783899129160)
	db.Exec("INSERT INTO customers (id, user_id, status, package_id, due_date_day, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		"CUST-1", "cust-u-1", "active", "pkg-1", 5, 1783899129160, 1783899129160)

	bodyJson := `{"notes":"Violation of terms"}`
	request := httptest.NewRequest(http.MethodPost, "/api/admin/customers/CUST-1/_suspend", strings.NewReader(bodyJson))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[model.CustomerResponse])
	json.Unmarshal(bytes, responseBody)
	assert.Equal(t, "suspended", responseBody.Data.Status)
}

func TestUnsuspendCustomerSuccess(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	db.Exec("INSERT INTO internet_packages (id, name, speed_mbps, price, installation_fee, tax_rate, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"pkg-1", "Test Package", 10, 100000, 50000, 0.11, true, 1783899129160, 1783899129160)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	db.Exec("INSERT INTO users (id, role_id, username, email, password, company_name, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		"cust-u-1", 2, "Agus", "agus@example.com", string(hashedPassword), "GREENET", 1783899129160, 1783899129160)
	db.Exec("INSERT INTO customers (id, user_id, status, package_id, due_date_day, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		"CUST-1", "cust-u-1", "suspended", "pkg-1", 5, 1783899129160, 1783899129160)

	bodyJson := `{"notes":"Paid outstanding bills"}`
	request := httptest.NewRequest(http.MethodPost, "/api/admin/customers/CUST-1/_unsuspend", strings.NewReader(bodyJson))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[model.CustomerResponse])
	json.Unmarshal(bytes, responseBody)
	assert.Equal(t, "active", responseBody.Data.Status)
}

func TestTerminateCustomerSuccess(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	db.Exec("INSERT INTO internet_packages (id, name, speed_mbps, price, installation_fee, tax_rate, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"pkg-1", "Test Package", 10, 100000, 50000, 0.11, true, 1783899129160, 1783899129160)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	db.Exec("INSERT INTO users (id, role_id, username, email, password, company_name, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		"cust-u-1", 2, "Agus", "agus@example.com", string(hashedPassword), "GREENET", 1783899129160, 1783899129160)
	db.Exec("INSERT INTO customers (id, user_id, status, package_id, due_date_day, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		"CUST-1", "cust-u-1", "active", "pkg-1", 5, 1783899129160, 1783899129160)

	bodyJson := `{"notes":"Customer request"}`
	request := httptest.NewRequest(http.MethodPost, "/api/admin/customers/CUST-1/_terminate", strings.NewReader(bodyJson))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[model.CustomerResponse])
	json.Unmarshal(bytes, responseBody)
	assert.Equal(t, "terminated", responseBody.Data.Status)
}
