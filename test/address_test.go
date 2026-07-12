package test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"golang-clean-architecture/internal/model"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateAddress(t *testing.T) {
	TestCreateContact(t)

	user := GetFirstUser(t)
	contact := GetFirstContact(t, user)

	requestBody := model.CreateAddressRequest{
		Street:     "Jalan Kemuning",
		City:       "Bandung",
		Province:   "Jawa Barat",
		PostalCode: "3422131",
		Country:    "Indonesia",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/contacts/"+contact.ID+"/addresses", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, requestBody.Street, responseBody.Data.Street)
	assert.Equal(t, requestBody.City, responseBody.Data.City)
	assert.Equal(t, requestBody.Province, responseBody.Data.Province)
	assert.Equal(t, requestBody.Country, responseBody.Data.Country)
	assert.Equal(t, requestBody.PostalCode, responseBody.Data.PostalCode)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
	assert.NotNil(t, responseBody.Data.ID)
}

func TestCreateAddressFailed(t *testing.T) {
	TestCreateContact(t)

	user := GetFirstUser(t)
	contact := GetFirstContact(t, user)

	requestBody := model.CreateAddressRequest{
		Street:     "Jalan Belum jadi",
		City:       "Jakarta",
		Province:   "DKI Jakarta",
		PostalCode: "2313211345345345353543535443",
		Country:    "Indonesia",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/contacts/"+contact.ID+"/addresses", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

}

func TestListAddresses(t *testing.T) {
	TestCreateContact(t)

	user := GetFirstUser(t)
	contact := GetFirstContact(t, user)

	CreateAddresses(t, contact, 5)

	request := httptest.NewRequest(http.MethodGet, "/api/contacts/"+contact.ID+"/addresses", nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 5, len(responseBody.Data))
}

func TestListAddressesFailed(t *testing.T) {
	TestCreateContact(t)

	user := GetFirstUser(t)
	contact := GetFirstContact(t, user)

	CreateAddresses(t, contact, 5)

	request := httptest.NewRequest(http.MethodGet, "/api/contacts/"+"wrong"+"/addresses", nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestGetAddressFailed(t *testing.T) {
	TestCreateAddress(t)

	user := GetFirstUser(t)
	contact := GetFirstContact(t, user)

	request := httptest.NewRequest(http.MethodGet, "/api/contacts/"+contact.ID+"/addresses/"+"wrong", nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestUpdateAddress(t *testing.T) {
	TestCreateAddress(t)

	user := GetFirstUser(t)
	contact := GetFirstContact(t, user)
	address := GetFirstAddress(t, contact)

	requestBody := model.CreateAddressRequest{
		Street:     "Jalan Lagi DIjieun",
		City:       "Bandung",
		Province:   "Jawa Barat",
		PostalCode: "1232131",
		Country:    "Indonesia",
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPut, "/api/contacts/"+contact.ID+"/addresses/"+address.ID, strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, requestBody.Street, responseBody.Data.Street)
	assert.Equal(t, requestBody.City, responseBody.Data.City)
	assert.Equal(t, requestBody.Province, responseBody.Data.Province)
	assert.Equal(t, requestBody.Country, responseBody.Data.Country)
	assert.Equal(t, requestBody.PostalCode, responseBody.Data.PostalCode)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
	assert.NotNil(t, responseBody.Data.ID)
}

func TestUpdateAddressFailed(t *testing.T) {
	TestCreateAddress(t)

	user := GetFirstUser(t)
	contact := GetFirstContact(t, user)
	address := GetFirstAddress(t, contact)

	requestBody := model.CreateAddressRequest{
		Street:     "Jalan Lagi DIjieun",
		City:       "Bandung",
		Province:   "Jawa Barat",
		PostalCode: "123213324324324234242424234234242424243242234321",
		Country:    "Indonesia",
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPut, "/api/contacts/"+contact.ID+"/addresses/"+address.ID, strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

}

func TestDeleteAddress(t *testing.T) {
	TestCreateAddress(t)

	user := GetFirstUser(t)
	contact := GetFirstContact(t, user)
	address := GetFirstAddress(t, contact)

	request := httptest.NewRequest(http.MethodDelete, "/api/contacts/"+contact.ID+"/addresses/"+address.ID, nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, true, responseBody.Data)
}

func TestDeleteAddressFailed(t *testing.T) {
	TestCreateAddress(t)

	user := GetFirstUser(t)
	contact := GetFirstContact(t, user)

	request := httptest.NewRequest(http.MethodDelete, "/api/contacts/"+contact.ID+"/addressess/"+"wrong", nil)

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}
