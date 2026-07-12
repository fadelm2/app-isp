package test

import (
	"golang-clean-architecture/internal/config"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var app *fiber.App

var db *gorm.DB
var viperConfig *viper.Viper

var log *logrus.Logger

var validate *validator.Validate

var secretJwt = "sdadsa"

func init() {
	viperConfig = config.NewViper()
	log = config.NewLogger(viperConfig)
	validate = config.NewValidator(viperConfig)
	app = config.NewFiber(viperConfig)
	db = config.NewDatabase(viperConfig, log)
	producer := config.NewKafkaProducer(viperConfig, log)
	secretJwt := secretJwt

	config.Bootstrap(&config.BootstrapConfig{
		DB:        db,
		App:       app,
		Log:       log,
		Validate:  validate,
		Config:    viperConfig,
		Producer:  producer,
		SecretKey: secretJwt,
	})
}
