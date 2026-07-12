package notification

import (
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
	c.Log.Infof("[Notification EMAIL] Send to: %s, Subject: %s", to, subject)
	return nil
}

func (c *NotificationClient) SendWhatsApp(phone string, message string) error {
	c.Log.Infof("[Notification WHATSAPP] Send to: %s, Msg: %s", phone, message)
	return nil
}

func (c *NotificationClient) SendTelegram(chatID string, message string) error {
	c.Log.Infof("[Notification TELEGRAM] Send to: %s, Msg: %s", chatID, message)
	return nil
}
