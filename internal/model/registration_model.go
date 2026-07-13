package model

type CreateRegistrationRequest struct {
	FullName            string  `json:"full_name" validate:"required"`
	NIK                 string  `json:"nik" validate:"required"`
	BirthPlace          string  `json:"birth_place" validate:"required"`
	BirthDate           string  `json:"birth_date" validate:"required"`
	Gender              string  `json:"gender" validate:"required"`
	Email               string  `json:"email" validate:"required,email"`
	Phone               string  `json:"phone" validate:"required"`
	InstallationAddress string  `json:"installation_address" validate:"required"`
	BillingAddress      string  `json:"billing_address" validate:"required"`
	PackageID           string  `json:"package_id" validate:"required"`
	Latitude            float64 `json:"latitude"`
	Longitude           float64 `json:"longitude"`
	Notes               string  `json:"notes"`
	KtpPath             string  `json:"ktp_path"`
	SelfiePath          string  `json:"selfie_path"`
	HousePath           string  `json:"house_path"`
	InstallationPath    string  `json:"installation_path"`
	SupportingDocPath   string  `json:"supporting_doc_path"`
	OdpNumber           string  `json:"odp_number"`
	Province            string  `json:"province"`
	City                string  `json:"city"`
	District            string  `json:"district"`
	Village             string  `json:"village"`
}

type UpdateRegistrationStatusRequest struct {
	ID        string `json:"id" validate:"required"`
	Status    string `json:"status" validate:"required"` // pending, under_review, surveying, approved, rejected
	OdpNumber string `json:"odp_number"`
}

type RegistrationResponse struct {
	ID                  string          `json:"id"`
	FullName            string          `json:"full_name"`
	NIK                 string          `json:"nik"`
	BirthPlace          string          `json:"birth_place"`
	BirthDate           string          `json:"birth_date"`
	Gender              string          `json:"gender"`
	Email               string          `json:"email"`
	Phone               string          `json:"phone"`
	InstallationAddress string          `json:"installation_address"`
	BillingAddress      string          `json:"billing_address"`
	PackageID           string          `json:"package_id"`
	Latitude            float64         `json:"latitude"`
	Longitude           float64         `json:"longitude"`
	Notes               string          `json:"notes"`
	Status              string          `json:"status"`
	KtpPath             string          `json:"ktp_path"`
	SelfiePath          string          `json:"selfie_path"`
	HousePath           string          `json:"house_path"`
	InstallationPath    string          `json:"installation_path"`
	SupportingDocPath   string          `json:"supporting_doc_path"`
	OdpNumber           string          `json:"odp_number"`
	Province            string          `json:"province"`
	City                string          `json:"city"`
	District            string          `json:"district"`
	Village             string          `json:"village"`
	CreatedAt           int64           `json:"created_at"`
	UpdatedAt           int64           `json:"updated_at"`
	Package             PackageResponse `json:"package"`
}

type SearchRegistrationRequest struct {
	Search string `json:"search"`
	Status string `json:"status"`
	Page   int    `json:"page" validate:"min=1"`
	Size   int    `json:"size" validate:"min=1,max=100"`
}
