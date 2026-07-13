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
	"golang.org/x/crypto/bcrypt"
)

func TestCreateInvoiceSuccess(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	db.Exec("INSERT INTO internet_packages (id, name, speed_mbps, price, installation_fee, tax_rate, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"pkg-1", "Test Package", 10, 100000, 50000, 0.11, true, 1783899129160, 1783899129160)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	db.Exec("INSERT INTO users (id, role_id, username, email, password, company_name, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		"cust-u-1", 2, "Agus", "agus@example.com", string(hashedPassword), "GREENET", 1783899129160, 1783899129160)
	db.Exec("INSERT INTO customers (id, user_id, status, package_id, due_date_day, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		"CUST-1", "cust-u-1", "active", "pkg-1", 5, 1783899129160, 1783899129160)

	requestBody := model.CreateInvoiceRequest{
		CustomerID:  "CUST-1",
		PeriodMonth: 7,
		PeriodYear:  2026,
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/admin/invoices", bytes.NewReader(bodyJson))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.InvoiceResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, requestBody.CustomerID, responseBody.Data.CustomerID)
	assert.Equal(t, float64(100000), responseBody.Data.Amount)
}

func TestCreateInvoiceFailedValidation(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	// Month out of range (failed validation)
	requestBody := model.CreateInvoiceRequest{
		CustomerID:  "CUST-1",
		PeriodMonth: 13,
		PeriodYear:  2026,
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/admin/invoices", bytes.NewReader(bodyJson))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestListInvoicesSuccess(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	db.Exec("INSERT INTO internet_packages (id, name, speed_mbps, price, installation_fee, tax_rate, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"pkg-1", "Test Package", 10, 100000, 50000, 0.11, true, 1783899129160, 1783899129160)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	db.Exec("INSERT INTO users (id, role_id, username, email, password, company_name, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		"cust-u-1", 2, "Agus", "agus@example.com", string(hashedPassword), "GREENET", 1783899129160, 1783899129160)
	db.Exec("INSERT INTO customers (id, user_id, status, package_id, due_date_day, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		"CUST-1", "cust-u-1", "active", "pkg-1", 5, 1783899129160, 1783899129160)
	db.Exec("INSERT INTO invoices (id, customer_id, due_date, period_month, period_year, amount, tax_amount, installation_fee, total_amount, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"INV-1", "CUST-1", 1783899129160, 7, 2026, 100000, 11000, 0, 111000, "pending", 1783899129160, 1783899129160)

	request := httptest.NewRequest(http.MethodGet, "/api/admin/invoices", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[[]model.InvoiceResponse])
	json.Unmarshal(bytes, responseBody)
	assert.True(t, len(responseBody.Data) > 0)
}

func TestGetInvoiceSuccess(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	db.Exec("INSERT INTO internet_packages (id, name, speed_mbps, price, installation_fee, tax_rate, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"pkg-1", "Test Package", 10, 100000, 50000, 0.11, true, 1783899129160, 1783899129160)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	db.Exec("INSERT INTO users (id, role_id, username, email, password, company_name, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		"cust-u-1", 2, "Agus", "agus@example.com", string(hashedPassword), "GREENET", 1783899129160, 1783899129160)
	db.Exec("INSERT INTO customers (id, user_id, status, package_id, due_date_day, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		"CUST-1", "cust-u-1", "active", "pkg-1", 5, 1783899129160, 1783899129160)
	db.Exec("INSERT INTO invoices (id, customer_id, due_date, period_month, period_year, amount, tax_amount, installation_fee, total_amount, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"INV-1", "CUST-1", 1783899129160, 7, 2026, 100000, 11000, 0, 111000, "pending", 1783899129160, 1783899129160)

	request := httptest.NewRequest(http.MethodGet, "/api/admin/invoices/INV-1", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[model.InvoiceResponse])
	json.Unmarshal(bytes, responseBody)
	assert.Equal(t, "INV-1", responseBody.Data.ID)
}

func TestGetInvoiceFailedNotFound(t *testing.T) {
	ClearAll()
	token := getAdminToken(t)

	request := httptest.NewRequest(http.MethodGet, "/api/admin/invoices/INV-NOT-EXIST", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestListPublicCustomerInvoicesSuccess(t *testing.T) {
	ClearAll()

	db.Exec("INSERT INTO internet_packages (id, name, speed_mbps, price, installation_fee, tax_rate, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"pkg-1", "Test Package", 10, 100000, 50000, 0.11, true, 1783899129160, 1783899129160)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	db.Exec("INSERT INTO users (id, role_id, username, email, password, company_name, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		"cust-u-1", 2, "Agus", "agus@example.com", string(hashedPassword), "GREENET", 1783899129160, 1783899129160)
	db.Exec("INSERT INTO customers (id, user_id, status, package_id, due_date_day, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		"CUST-1", "cust-u-1", "active", "pkg-1", 5, 1783899129160, 1783899129160)
	db.Exec("INSERT INTO invoices (id, customer_id, due_date, period_month, period_year, amount, tax_amount, installation_fee, total_amount, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"INV-1", "CUST-1", 1783899129160, 7, 2026, 100000, 11000, 0, 111000, "pending", 1783899129160, 1783899129160)

	request := httptest.NewRequest(http.MethodGet, "/api/public/customers/CUST-1/invoices", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[[]model.InvoiceResponse])
	json.Unmarshal(bytes, responseBody)
	assert.Equal(t, 1, len(responseBody.Data))
	assert.Equal(t, "INV-1", responseBody.Data[0].ID)
}

func TestListPublicCustomerInvoicesFailedNotFound(t *testing.T) {
	ClearAll()

	request := httptest.NewRequest(http.MethodGet, "/api/public/customers/CUST-NOT-EXIST/invoices", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestGetPublicSnapTokenFailedNotFound(t *testing.T) {
	ClearAll()

	request := httptest.NewRequest(http.MethodGet, "/api/public/invoices/INV-NOT-EXIST/pay", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}
