package entity

type Payment struct {
	ID            string  `gorm:"column:id;primaryKey"`
	InvoiceID     string  `gorm:"column:invoice_id"`
	TransactionID string  `gorm:"column:transaction_id"`
	PaymentType   string  `gorm:"column:payment_type"`
	PaidAmount    float64 `gorm:"column:paid_amount"`
	Status        string  `gorm:"column:status"`
	PaidAt        int64   `gorm:"column:paid_at"`
	RawResponse   string  `gorm:"column:raw_response"`
	Invoice       Invoice `gorm:"foreignKey:invoice_id;references:id"`
}

func (p *Payment) TableName() string {
	return "payments"
}
