package test

import (
	"fmt"
	"golang-clean-architecture/internal/config"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
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
	viperConfig = viper.New()
	viperConfig.SetConfigName("config_test")
	viperConfig.SetConfigType("json")
	viperConfig.AddConfigPath("./../")
	viperConfig.AddConfigPath("./")
	err := viperConfig.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error test config file: %w \n", err))
	}

	log = config.NewLogger(viperConfig)

	// Pre-create the database if it doesn't exist
	username := viperConfig.GetString("database.username")
	password := viperConfig.GetString("database.password")
	host := viperConfig.GetString("database.host")
	port := viperConfig.GetInt("database.port")

	rootDsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port)
	tempDb, err := gorm.Open(mysql.Open(rootDsn), &gorm.Config{})
	if err == nil {
		tempDb.Exec("CREATE DATABASE IF NOT EXISTS app_isptest;")
		sqlDb, _ := tempDb.DB()
		sqlDb.Close()
	}

	validate = config.NewValidator(viperConfig)
	app = config.NewFiber(viperConfig)
	db = config.NewDatabase(viperConfig, log)

	runMigrations(db)

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

func runMigrations(db *gorm.DB) {
	// Drop existing tables
	db.Exec("DROP TABLE IF EXISTS radacct")
	db.Exec("DROP TABLE IF EXISTS radreply")
	db.Exec("DROP TABLE IF EXISTS radcheck")
	db.Exec("DROP TABLE IF EXISTS customer_histories")
	db.Exec("DROP TABLE IF EXISTS payments")
	db.Exec("DROP TABLE IF EXISTS invoices")
	db.Exec("DROP TABLE IF EXISTS customers")
	db.Exec("DROP TABLE IF EXISTS routers")
	db.Exec("DROP TABLE IF EXISTS registrations")
	db.Exec("DROP TABLE IF EXISTS internet_packages")
	db.Exec("DROP TABLE IF EXISTS addresses")
	db.Exec("DROP TABLE IF EXISTS contacts")
	db.Exec("DROP TABLE IF EXISTS users")
	db.Exec("DROP TABLE IF EXISTS roles")

	// Read migration directory
	migrationDir := "./db/migrations"
	var files []os.DirEntry
	var err error
	
	// Read directory entries
	files, err = os.ReadDir(migrationDir)
	if err != nil {
		migrationDir = "./../db/migrations"
		files, err = os.ReadDir(migrationDir)
		if err != nil {
			panic("Failed to find db/migrations directory")
		}
	}

	var upFiles []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".up.sql") {
			upFiles = append(upFiles, f.Name())
		}
	}
	sort.Strings(upFiles)

	for _, filename := range upFiles {
		filePath := filepath.Join(migrationDir, filename)
		contentBytes, err := os.ReadFile(filePath)
		if err != nil {
			panic(fmt.Sprintf("Failed to read migration file %s: %v", filename, err))
		}

		content := string(contentBytes)
		queries := splitQueries(content)
		for _, q := range queries {
			q = strings.TrimSpace(q)
			if q == "" {
				continue
			}
			if err := db.Exec(q).Error; err != nil {
				panic(fmt.Sprintf("Failed to execute query in %s: %v\nQuery: %s", filename, err, q))
			}
		}
	}
}

func splitQueries(content string) []string {
	// Remove SQL comments
	re := regexp.MustCompile(`(?m)^--.*$`)
	content = re.ReplaceAllString(content, "")
	
	reBlock := regexp.MustCompile(`(?s)/\*.*?\*/`)
	content = reBlock.ReplaceAllString(content, "")

	return strings.Split(content, ";")
}
