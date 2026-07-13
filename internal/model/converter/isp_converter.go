package converter

import (
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
)

func PackageToResponse(pkg *entity.InternetPackage) model.PackageResponse {
	return model.PackageResponse{
		ID:              pkg.ID,
		Name:            pkg.Name,
		SpeedMbps:       pkg.SpeedMbps,
		Price:           pkg.Price,
		InstallationFee: pkg.InstallationFee,
		TaxRate:         pkg.TaxRate,
		IsActive:        pkg.IsActive,
		CreatedAt:       pkg.CreatedAt,
		UpdatedAt:       pkg.UpdatedAt,
	}
}

func RegistrationToResponse(reg *entity.Registration) model.RegistrationResponse {
	return model.RegistrationResponse{
		ID:                  reg.ID,
		FullName:            reg.FullName,
		NIK:                 reg.NIK,
		BirthPlace:          reg.BirthPlace,
		BirthDate:           reg.BirthDate,
		Gender:              reg.Gender,
		Email:               reg.Email,
		Phone:               reg.Phone,
		InstallationAddress: reg.InstallationAddress,
		BillingAddress:      reg.BillingAddress,
		PackageID:           reg.PackageID,
		Latitude:            reg.Latitude,
		Longitude:           reg.Longitude,
		Notes:               reg.Notes,
		Status:              reg.Status,
		KtpPath:             reg.KtpPath,
		SelfiePath:          reg.SelfiePath,
		HousePath:           reg.HousePath,
		InstallationPath:    reg.InstallationPath,
		SupportingDocPath:   reg.SupportingDocPath,
		OdpNumber:           reg.OdpNumber,
		Province:            reg.Province,
		City:                reg.City,
		District:            reg.District,
		Village:             reg.Village,
		CreatedAt:           reg.CreatedAt,
		UpdatedAt:           reg.UpdatedAt,
		Package:             PackageToResponse(&reg.Package),
	}
}

func RouterToResponse(router *entity.Router) *model.RouterResponse {
	if router == nil {
		return nil
	}
	return &model.RouterResponse{
		ID:        router.ID,
		Name:      router.Name,
		Host:      router.Host,
		Port:      router.Port,
		Status:    router.Status,
		CreatedAt: router.CreatedAt,
		UpdatedAt: router.UpdatedAt,
	}
}

func CustomerToResponse(customer *entity.Customer) model.CustomerResponse {
	var routerResp *model.RouterResponse
	if customer.Router != nil {
		routerResp = RouterToResponse(customer.Router)
	}
	var regResp *model.RegistrationResponse
	if customer.Registration != nil {
		resp := RegistrationToResponse(customer.Registration)
		regResp = &resp
	}
	return model.CustomerResponse{
		ID:             customer.ID,
		UserID:         customer.UserID,
		Status:         customer.Status,
		PackageID:      customer.PackageID,
		RouterID:       customer.RouterID,
		PppUsername:    customer.PppUsername,
		RadiusUsername: customer.RadiusUsername,
		DueDateDay:     customer.DueDateDay,
		OdpNumber:      customer.OdpNumber,
		CreatedAt:      customer.CreatedAt,
		UpdatedAt:      customer.UpdatedAt,
		User:           *UserToResponse(&customer.User),
		Package:        PackageToResponse(&customer.Package),
		Router:         routerResp,
		Registration:   regResp,
	}
}

func InvoiceToResponse(invoice *entity.Invoice) model.InvoiceResponse {
	var paidAtVal *int64
	if invoice.PaidAt != nil {
		paidAtVal = invoice.PaidAt
	}
	return model.InvoiceResponse{
		ID:              invoice.ID,
		CustomerID:      invoice.CustomerID,
		DueDate:         invoice.DueDate,
		PeriodMonth:     invoice.PeriodMonth,
		PeriodYear:      invoice.PeriodYear,
		Amount:          invoice.Amount,
		TaxAmount:       invoice.TaxAmount,
		InstallationFee: invoice.InstallationFee,
		TotalAmount:     invoice.TotalAmount,
		Status:          invoice.Status,
		SnapToken:       invoice.SnapToken,
		PaidAt:          paidAtVal,
		CreatedAt:       invoice.CreatedAt,
		UpdatedAt:       invoice.UpdatedAt,
		Customer:        CustomerToResponse(&invoice.Customer),
	}
}

func CustomerHistoryToResponse(history *entity.CustomerHistory) model.CustomerHistoryResponse {
	return model.CustomerHistoryResponse{
		ID:         history.ID,
		CustomerID: history.CustomerID,
		Action:     history.Action,
		Notes:      history.Notes,
		CreatedBy:  history.CreatedBy,
		CreatedAt:  history.CreatedAt,
		User:       *UserToResponse(&history.User),
	}
}
