package notification

import (
	"os"

	"github.com/sirupsen/logrus"
)

type NotificationClient struct {
	Log *logrus.Logger
}

func NewNotificationClient(log *logrus.Logger) *NotificationClient {
	return &NotificationClient{
		Log: log,
	}
}

func (c *NotificationClient) SendEmail(to string, subject string, body string) error {
	sender := os.Getenv("SMTP_SENDER_EMAIL")
	if sender == "" {
		sender = "billing@greenet.id"
	}
	c.Log.Infof("[Notification EMAIL] Sender: %s, Send to: %s, Subject: %s", sender, to, subject)
	return nil
}

func (c *NotificationClient) SendWhatsApp(phone string, message string) error {
	sender := os.Getenv("WHATSAPP_SENDER_NUMBER")
	if sender == "" {
		sender = "628123456789"
	}
	c.Log.Infof("[Notification WHATSAPP] Sender (Owner): %s, Send to: %s, Msg: %s", sender, phone, message)
	return nil
}

func (c *NotificationClient) SendTelegram(chatID string, message string) error {
	c.Log.Infof("[Notification TELEGRAM] Send to: %s, Msg: %s", chatID, message)
	return nil
}
