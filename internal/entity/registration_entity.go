package entity

type Registration struct {
	ID                  string          `gorm:"column:id;primaryKey"`
	FullName            string          `gorm:"column:full_name"`
	NIK                 string          `gorm:"column:nik"`
	BirthPlace          string          `gorm:"column:birth_place"`
	BirthDate           string          `gorm:"column:birth_date"`
	Gender              string          `gorm:"column:gender"`
	Email               string          `gorm:"column:email"`
	Phone               string          `gorm:"column:phone"`
	InstallationAddress string          `gorm:"column:installation_address"`
	BillingAddress      string          `gorm:"column:billing_address"`
	PackageID           string          `gorm:"column:package_id"`
	Latitude            float64         `gorm:"column:latitude"`
	Longitude           float64         `gorm:"column:longitude"`
	Notes               string          `gorm:"column:notes"`
	Status              string          `gorm:"column:status"` // pending, under_review, surveying, approved, rejected
	KtpPath             string          `gorm:"column:ktp_path"`
	SelfiePath          string          `gorm:"column:selfie_path"`
	HousePath           string          `gorm:"column:house_path"`
	InstallationPath    string          `gorm:"column:installation_path"`
	SupportingDocPath   string          `gorm:"column:supporting_doc_path"`
	OdpNumber           string          `gorm:"column:odp_number"`
	Province            string          `gorm:"column:province"`
	City                string          `gorm:"column:city"`
	District            string          `gorm:"column:district"`
	Village             string          `gorm:"column:village"`
	CreatedAt           int64           `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt           int64           `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	Package             InternetPackage `gorm:"foreignKey:package_id;references:id;->"`
}

func (r *Registration) TableName() string {
	return "registrations"
}
