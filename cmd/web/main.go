package main

import (
	"fmt"
	"golang-clean-architecture/internal/config"

	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load()

	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)
	secretkey := config.SecretKey(viperConfig)
	//producer := config.NewKafkaProducer(viperConfig, log)

	config.Bootstrap(&config.BootstrapConfig{
		DB:        db,
		App:       app,
		Log:       log,
		Validate:  validate,
		Config:    viperConfig,
		SecretKey: secretkey,
		//Producer: producer,
	})

	webPort := viperConfig.GetInt("web.port")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
