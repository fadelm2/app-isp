package middleware

import (
	"golang-clean-architecture/internal/usecase"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func NewRequestLogger(userUserCase *usecase.UserUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userUserCase.Log.Info("RequestLogger HIT") // 🔥 test log

		start := time.Now()
		requestID := uuid.New().String()

		c.Locals("request_id", requestID)

		err := c.Next()

		status := c.Response().StatusCode()

		entry := userUserCase.Log.WithFields(logrus.Fields{
			"request_id": requestID,
			"method":     c.Method(),
			"path":       c.OriginalURL(),
			"status":     status,
			"latency_ms": time.Since(start).Milliseconds(),
			"ip":         c.IP(),
		})

		// ❌ jangan log detail untuk endpoint auth
		if strings.Contains(c.Path(), "/login") ||
			strings.Contains(c.Path(), "/register") {
			entry.Info("request finished (auth endpoint)")
			return err
		}
		userUserCase.Log.Info("After Next()") // 🔥 test log

		entry.Info("request finished")
		return err
	}
}
