package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
)

func NewFiber(config *viper.Viper) *fiber.App {
	var app = fiber.New(fiber.Config{
		AppName:      config.GetString("app.name"),
		ErrorHandler: NewErrorHandler(),
		Prefork:      config.GetBool("web.prefork"),
		ColorScheme:  fiber.DefaultColors,
	})
	var link1 = config.GetString("app.link1")
	var link2 = config.GetString("app.link2")
	var link3 = config.GetString("app.link3")
	var link4 = config.GetString("app.link4")

	corsSettings := cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins: link1 + "," +
			link2 + "," +
			link3 + "," +
			link4 + "," +
			"https://kidul.greenet.id," +
			"https://ticket.greenet.id," +
			"https://apiticket.greenet.id," +
			"https://portal.greenet.id," +
			"https://ticketapi.greenet.id," +
			"https://ticket.fadelweb.site," +
			"https://ticketapi.fadelweb.site," +
			"http://172.16.21.211:3000," +
			"http://192.168.1.70:3000," +
			"http://192.168.1.70:8080,",
		AllowMethods: "GET,POST,HEAD,OPTIONS,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept,  Accept-Encoding, X-CSRF-Token, Authorization,X-Requested-With",
		//   ExposeHeaders:    "Origin",
	})
	app.Use(corsSettings)

	return app
}
func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		return ctx.Status(code).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}
}
