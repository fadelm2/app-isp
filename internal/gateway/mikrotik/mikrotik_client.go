package mikrotik

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type MikrotikClient struct {
	Log *logrus.Logger
}

func NewMikrotikClient(log *logrus.Logger) *MikrotikClient {
	return &MikrotikClient{
		Log: log,
	}
}

type PppSecret struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Service  string `json:"service"`
	Profile  string `json:"profile"`
	Disabled string `json:"disabled"`
}

func (c *MikrotikClient) getHttpClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func (c *MikrotikClient) executeREST(host string, port int, username string, password string, method string, path string, body interface{}) ([]byte, error) {
	if host == "" {
		// Mock Mode
		c.Log.Infof("[MikroTik MOCK] Simulating REST %s to %s", method, path)
		return []byte(`[]`), nil
	}

	url := fmt.Sprintf("https://%s:%d/rest%s", host, port, path)
	if port == 0 {
		url = fmt.Sprintf("https://%s/rest%s", host, path)
	}

	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	client := c.getHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return respBytes, fmt.Errorf("mikrotik rest error: status %d, body: %s", resp.StatusCode, string(respBytes))
	}

	return respBytes, nil
}

func (c *MikrotikClient) PingRouter(host string, port int, username string, password string) (bool, error) {
	if host == "" {
		return true, nil // Mock success
	}
	_, err := c.executeREST(host, port, username, password, "GET", "/system/resource", nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *MikrotikClient) CreatePPPoESecret(host string, port int, username string, password string, pppUser string, pppPass string, speedProfile string) error {
	body := map[string]string{
		"name":     pppUser,
		"password": pppPass,
		"service":  "pppoe",
		"profile":  speedProfile,
	}
	_, err := c.executeREST(host, port, username, password, "PUT", "/ppp/secret", body)
	if err != nil {
		// Try POST if PUT fails
		_, err = c.executeREST(host, port, username, password, "POST", "/ppp/secret", body)
	}
	return err
}

func (c *MikrotikClient) EnablePPPoESecret(host string, port int, username string, password string, pppUser string) error {
	// Find ID first or directly patch using name
	path := fmt.Sprintf("/ppp/secret/%s", pppUser)
	body := map[string]string{
		"disabled": "no",
	}
	_, err := c.executeREST(host, port, username, password, "PATCH", path, body)
	return err
}

func (c *MikrotikClient) DisablePPPoESecret(host string, port int, username string, password string, pppUser string) error {
	path := fmt.Sprintf("/ppp/secret/%s", pppUser)
	body := map[string]string{
		"disabled": "yes",
	}
	_, err := c.executeREST(host, port, username, password, "PATCH", path, body)
	return err
}

func (c *MikrotikClient) DisconnectActiveSession(host string, port int, username string, password string, pppUser string) error {
	// Find active session and remove it
	path := fmt.Sprintf("/ppp/active/%s", pppUser)
	_, err := c.executeREST(host, port, username, password, "DELETE", path, nil)
	return err
}
