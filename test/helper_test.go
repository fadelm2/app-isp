package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func ClearISPTables() {
	db.Exec("DELETE FROM payments")
	db.Exec("DELETE FROM invoices")
	db.Exec("DELETE FROM customer_histories")
	db.Exec("DELETE FROM customers")
	db.Exec("DELETE FROM registrations")
	db.Exec("DELETE FROM internet_packages")
	db.Exec("DELETE FROM routers")
}

func ClearAll() {
	ClearISPTables()
	ClearAddresses()
	ClearContact()
	ClearUsers()
}

func ClearUsers() {
	err := db.Where("id is not null").Delete(&entity.User{}).Error
	if err != nil {
		log.Fatalf("Failed clear user data : %+v", err)
	}
}

func ClearContact() {
	err := db.Where("id is not null").Delete(&entity.Contact{}).Error
	if err != nil {
		log.Fatalf("Failed clear user data : %+v", err)
	}
}

func ClearAddresses() {
	err := db.Where("id is not null").Delete(&entity.Address{}).Error
	if err != nil {
		log.Fatalf("Failed clear address data : %+v", err)
	}
}

func CreateContacts(user *entity.User, total int) {
	for i := 0; i < total; i++ {
		contact := entity.Contact{
			ID:        uuid.NewString(),
			FirstName: "Contact",
			LastName:  strconv.Itoa(i),
			Email:     "contact" + strconv.Itoa(i) + "@example.com",
			Phone:     "080000000" + strconv.Itoa(i),
			UserId:    user.ID,
		}
		err := db.Create(&contact).Error
		if err != nil {
			log.Fatalf("Failed create contact data :%+v", err)
		}
	}
}

func CreateAddresses(t *testing.T, contact *entity.Contact, total int) {
	for i := 0; i < total; i++ {
		address := &entity.Address{
			ID:         uuid.NewString(),
			ContactId:  contact.ID,
			Street:     "Jalan Udin belum jadi",
			City:       "Kajarta",
			Province:   "DKI Jakarta",
			PostalCode: "21321412",
			Country:    "Indonesia",
		}
		err := db.Create(address).Error
		assert.Nil(t, err)
	}
}

func GetFirstUser(t *testing.T) *entity.User {
	user := new(entity.User)
	err := db.First(user).Error
	assert.Nil(t, err)
	return user
}

func GetFirstContact(t *testing.T, user *entity.User) *entity.Contact {
	contact := new(entity.Contact)
	err := db.Where("user_id = ?", user.ID).First(contact).Error
	assert.Nil(t, err)
	return contact
}

func GetFirstAddress(t *testing.T, contact *entity.Contact) *entity.Address {
	address := new(entity.Address)
	err := db.Where("contact_id = ?", contact.ID).First(address).Error
	assert.Nil(t, err)
	return address
}

func getAdminToken(t *testing.T) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("fadel123"), bcrypt.DefaultCost)
	db.Exec("DELETE FROM users WHERE id = ?", "admin1")
	db.Exec("INSERT INTO users (id, role_id, username, email, password, company_name, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		"admin1", 1, "admin_user", "admin@greenet.id", string(hashedPassword), "GREENET", 1783899129160, 1783899129160)

	loginReq := model.LoginUserRequest{
		ID:       "admin1",
		Password: "fadel123",
	}
	bodyJson, err := json.Marshal(loginReq)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/users/_Login", bytes.NewReader(bodyJson))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	var responseBody model.WebResponse[model.UserResponse]
	err = json.Unmarshal(bodyBytes, &responseBody)
	assert.Nil(t, err)

	return responseBody.Data.Token
}
