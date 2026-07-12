package model

type CreateInvoiceRequest struct {
	CustomerID  string `json:"customer_id" validate:"required"`
	PeriodMonth int    `json:"period_month" validate:"required,numeric"`
	PeriodYear  int    `json:"period_year" validate:"required,numeric"`
}

type InvoiceResponse struct {
	ID              string           `json:"id"`
	CustomerID      string           `json:"customer_id"`
	DueDate         int64            `json:"due_date"`
	PeriodMonth     int              `json:"period_month"`
	PeriodYear      int              `json:"period_year"`
	Amount          float64          `json:"amount"`
	TaxAmount       float64          `json:"tax_amount"`
	InstallationFee float64          `json:"installation_fee"`
	TotalAmount     float64          `json:"total_amount"`
	Status          string           `json:"status"` // pending, paid, owed, expired, cancelled
	SnapToken       string           `json:"snap_token,omitempty"`
	PaidAt          *int64           `json:"paid_at,omitempty"`
	CreatedAt       int64            `json:"created_at"`
	UpdatedAt       int64            `json:"updated_at"`
	Customer        CustomerResponse `json:"customer"`
}
