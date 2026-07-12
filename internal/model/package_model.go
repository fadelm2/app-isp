package model

type CreatePackageRequest struct {
	Name            string  `json:"name" validate:"required"`
	SpeedMbps       int     `json:"speed_mbps" validate:"required,numeric"`
	Price           float64 `json:"price" validate:"required,numeric"`
	InstallationFee float64 `json:"installation_fee" validate:"numeric"`
	TaxRate         float64 `json:"tax_rate" validate:"numeric"`
}

type UpdatePackageRequest struct {
	ID              string   `json:"id" validate:"required"`
	Name            string   `json:"name"`
	SpeedMbps       int      `json:"speed_mbps"`
	Price           float64  `json:"price"`
	InstallationFee float64  `json:"installation_fee"`
	TaxRate         float64  `json:"tax_rate"`
	IsActive        *bool    `json:"is_active"`
}

type PackageResponse struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	SpeedMbps       int     `json:"speed_mbps"`
	Price           float64 `json:"price"`
	InstallationFee float64 `json:"installation_fee"`
	TaxRate         float64 `json:"tax_rate"`
	IsActive        bool    `json:"is_active"`
	CreatedAt       int64   `json:"created_at"`
	UpdatedAt       int64   `json:"updated_at"`
}
