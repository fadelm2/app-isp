package model

type CreateCustomerRequest struct {
	RegistrationID string `json:"registration_id"`
	UserID         string `json:"user_id" validate:"required"`
	PackageID      string `json:"package_id" validate:"required"`
	RouterID       string `json:"router_id"`
	PppUsername    string `json:"ppp_username" validate:"required"`
	PppPassword    string `json:"ppp_password" validate:"required"`
	RadiusUsername string `json:"radius_username" validate:"required"`
	RadiusPassword string `json:"radius_password" validate:"required"`
	DueDateDay     int    `json:"due_date_day" validate:"required,numeric"`
}

type UpdateCustomerRequest struct {
	ID             string `json:"id" validate:"required"`
	Status         string `json:"status"` // active, suspended, terminated
	PackageID      string `json:"package_id"`
	RouterID       string `json:"router_id"`
	PppUsername    string `json:"ppp_username"`
	PppPassword    string `json:"ppp_password"`
	RadiusUsername string `json:"radius_username"`
	RadiusPassword string `json:"radius_password"`
	DueDateDay     int    `json:"due_date_day"`
	OdpNumber      string `json:"odp_number"`
}

type CustomerResponse struct {
	ID             string                `json:"id"`
	UserID         string                `json:"user_id"`
	Status         string                `json:"status"`
	PackageID      string                `json:"package_id"`
	RouterID       *string               `json:"router_id"`
	PppUsername    string                `json:"ppp_username"`
	RadiusUsername string                `json:"radius_username"`
	DueDateDay     int                   `json:"due_date_day"`
	OdpNumber      string                `json:"odp_number"`
	CreatedAt      int64                 `json:"created_at"`
	UpdatedAt      int64                 `json:"updated_at"`
	User           UserResponse          `json:"user"`
	Package        PackageResponse       `json:"package"`
	Router         *RouterResponse       `json:"router"`
	Registration   *RegistrationResponse `json:"registration"`
}

type CustomerHistoryResponse struct {
	ID         string       `json:"id"`
	CustomerID string       `json:"customer_id"`
	Action     string       `json:"action"`
	Notes      string       `json:"notes"`
	CreatedBy  string       `json:"created_by"`
	CreatedAt  int64        `json:"created_at"`
	User       UserResponse `json:"user"`
}

type SearchCustomerRequest struct {
	Search string `json:"search"`  // search by name, email, ppp_username, customer ID, odp
	Status string `json:"status"`  // filter: active, suspended, terminated
	Page   int    `json:"page" validate:"min=1"`
	Size   int    `json:"size" validate:"min=1,max=100"`
}
