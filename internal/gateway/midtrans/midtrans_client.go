package midtrans

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type MidtransClient struct {
	ServerKey    string
	IsProduction bool
	Log          *logrus.Logger
}

func NewMidtransClient(serverKey string, isProduction bool, log *logrus.Logger) *MidtransClient {
	return &MidtransClient{
		ServerKey:    serverKey,
		IsProduction: isProduction,
		Log:          log,
	}
}

func (c *MidtransClient) getBaseURL() string {
	if c.IsProduction {
		return "https://app.midtrans.com"
	}
	return "https://app.sandbox.midtrans.com"
}

func (c *MidtransClient) CreateSnapToken(invoiceID string, totalAmount float64, customerName string, customerEmail string) (string, error) {
	if c.ServerKey == "" {
		// Mock token
		c.Log.Infof("[Midtrans MOCK] Creating Snap Token for Invoice: %s", invoiceID)
		return fmt.Sprintf("mock-snap-token-%s", invoiceID), nil
	}

	url := c.getBaseURL() + "/snap/v1/transactions"

	body := map[string]interface{}{
		"transaction_details": map[string]interface{}{
			"order_id":     invoiceID,
			"gross_amount": int64(totalAmount),
		},
		"customer_details": map[string]interface{}{
			"first_name": customerName,
			"email":      customerEmail,
		},
		"credit_card": map[string]interface{}{
			"secure": true,
		},
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(c.ServerKey, "")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("midtrans error: status %d, body %s", resp.StatusCode, string(respBytes))
	}

	var response struct {
		Token      string `json:"token"`
		RedirectURL string `json:"redirect_url"`
	}

	err = json.Unmarshal(respBytes, &response)
	if err != nil {
		return "", err
	}

	return response.Token, nil
}

func (c *MidtransClient) VerifySignature(orderID string, statusCode string, grossAmount string, signatureKey string) bool {
	// SHA512(order_id + status_code + gross_amount + ServerKey)
	payload := orderID + statusCode + grossAmount + c.ServerKey
	hasher := sha512.New()
	hasher.Write([]byte(payload))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash == signatureKey
}
