package util

import (
	"bytes"
	"fmt"
	"strings"
	"time"
	"golang-clean-architecture/internal/entity"

	"github.com/jung-kurt/gofpdf"
)

// GenerateInvoicePDF creates a styled PDF invoice and returns the byte slice
func GenerateInvoicePDF(invoice *entity.Invoice, ispName string) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()

	// Title / ISP Header Background Band
	pdf.SetFillColor(13, 138, 188) // #0D8ABC (Deep Sky Blue)
	pdf.Rect(0, 0, 210, 45, "F")

	// Header Text
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 24)
	pdf.Text(15, 25, strings.ToUpper(ispName))

	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(230, 245, 255)
	pdf.Text(15, 33, "Internet Service Provider Connection Invoice")

	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(255, 255, 255)
	pdf.Text(145, 25, "INVOICE BILL")

	// Divider space
	pdf.Ln(35)

	// Invoice Information Block
	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(13, 138, 188) // Primary Color
	pdf.CellFormat(90, 6, "BILL TO:", "", 0, "L", false, 0, "")
	pdf.CellFormat(90, 6, "INVOICE DETAILS:", "", 1, "R", false, 0, "")

	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(50, 50, 50)
	pdf.CellFormat(90, 6, invoice.Customer.User.Name, "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(90, 6, fmt.Sprintf("Invoice ID: %s", invoice.ID), "", 1, "R", false, 0, "")

	pdf.CellFormat(90, 6, fmt.Sprintf("Email: %s", invoice.Customer.User.Email), "", 0, "L", false, 0, "")
	dueDateStr := time.UnixMilli(invoice.DueDate).Format("02-01-2006")
	pdf.CellFormat(90, 6, fmt.Sprintf("Due Date: %s", dueDateStr), "", 1, "R", false, 0, "")

	// Customer Installation Address
	address := "N/A"
	if invoice.Customer.Registration != nil {
		address = invoice.Customer.Registration.InstallationAddress
	} else if len(invoice.Customer.User.Contacts) > 0 && len(invoice.Customer.User.Contacts[0].Addresses) > 0 {
		addrEntity := invoice.Customer.User.Contacts[0].Addresses[0]
		address = fmt.Sprintf("%s, %s, %s", addrEntity.Street, addrEntity.City, addrEntity.PostalCode)
	}

	pdf.CellFormat(90, 6, fmt.Sprintf("Address: %s", address), "", 0, "L", false, 0, "")
	
	statusText := strings.ToUpper(invoice.Status)
	pdf.SetFont("Arial", "B", 10)
	if invoice.Status == "paid" {
		pdf.SetTextColor(40, 167, 69) // Green
	} else if invoice.Status == "owed" {
		pdf.SetTextColor(220, 53, 69) // Red
	} else {
		pdf.SetTextColor(255, 193, 7) // Yellow/Warning
	}
	pdf.CellFormat(90, 6, fmt.Sprintf("Status: %s", statusText), "", 1, "R", false, 0, "")

	pdf.Ln(8)

	// Itemized Table Header
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(13, 138, 188) // Primary Color
	pdf.CellFormat(110, 8, "Item Description", "1", 0, "L", true, 0, "")
	pdf.CellFormat(70, 8, "Amount (IDR)", "1", 1, "R", true, 0, "")

	pdf.SetTextColor(50, 50, 50)
	pdf.SetFont("Arial", "", 10)
	
	// Line Item 1: Package
	pkgName := "N/A"
	if invoice.Customer.Package.Name != "" {
		pkgName = invoice.Customer.Package.Name
	}
	pdf.CellFormat(110, 8, fmt.Sprintf("Subscription Fee - %s", pkgName), "1", 0, "L", false, 0, "")
	pdf.CellFormat(70, 8, fmt.Sprintf("Rp %s", formatPrice(invoice.Amount)), "1", 1, "R", false, 0, "")

	// Line Item 2: Installation Fee (if any)
	if invoice.InstallationFee > 0 {
		pdf.CellFormat(110, 8, "Installation Setup Fee", "1", 0, "L", false, 0, "")
		pdf.CellFormat(70, 8, fmt.Sprintf("Rp %s", formatPrice(invoice.InstallationFee)), "1", 1, "R", false, 0, "")
	}

	// Line Item 3: Tax (VAT)
	if invoice.TaxAmount > 0 {
		pdf.CellFormat(110, 8, "Tax (VAT 11%)", "1", 0, "L", false, 0, "")
		pdf.CellFormat(70, 8, fmt.Sprintf("Rp %s", formatPrice(invoice.TaxAmount)), "1", 1, "R", false, 0, "")
	}

	// Total Row
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(245, 245, 245)
	pdf.CellFormat(110, 8, "Total Amount Due", "1", 0, "L", true, 0, "")
	pdf.CellFormat(70, 8, fmt.Sprintf("Rp %s", formatPrice(invoice.TotalAmount)), "1", 1, "R", true, 0, "")

	pdf.Ln(20)

	// Payment Notice Footer
	pdf.SetFont("Arial", "I", 9)
	pdf.SetTextColor(120, 120, 120)
	pdf.CellFormat(180, 5, fmt.Sprintf("Please complete the payment prior to the due date (%s).", dueDateStr), "", 1, "C", false, 0, "")
	pdf.CellFormat(180, 5, "Thank you for subscribing to our services!", "", 1, "C", false, 0, "")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func formatPrice(val float64) string {
	str := fmt.Sprintf("%.0f", val)
	var result []string
	length := len(str)
	for i, ch := range str {
		result = append(result, string(ch))
		if (length-i-1)%3 == 0 && i != length-1 {
			result = append(result, ".")
		}
	}
	return strings.Join(result, "")
}
