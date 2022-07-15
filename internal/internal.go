package internal

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/spoonboy-io/koan"
)

// App is provides a wrapper for passing the Config, Context and Logger around as dependencies
type App struct {
	Logger *koan.Logger
	Ctx    context.Context
	Config struct {
		MorpheusHost  string
		MorpheusToken string
		PollInterval  int
		SmtpServer    string
		SmtpPort      int
		SmtpUser      string
		SmtpPassword  string
	}
}

const (
	// config
	APP_CONFIG      = "config.env"
	APPROVAL_CONFIG = "approvals.yaml"

	// server
	SRV_HOST = ""
	SRV_PORT = "18652"

	// email templates
	TEMPLATE_FOLDER = "templates"

	// tls configuration
	TLS_FOLDER    = "certs"
	TLS_ORG       = "Spoon Boy"
	TLS_VALID_FOR = 365 * 24 * time.Hour
)

var DefaultTemplate string = `Template here TODO`
var (
	ERR_FAILED_READ_CONFIG    = errors.New("Failed to read application configuration file")
	ERR_NO_API_HOST           = errors.New("No Morpheus API Host found")
	ERR_NO_API_TOKEN          = errors.New("No Morpheus API Token found")
	ERR_POLL_INTERVAL_NOT_INT = errors.New("Poll interval not integer")
	ERR_NO_SMTP_SERVER        = errors.New("No SMTP server found")
	ERR_NO_SMTP_PORT          = errors.New("No SMTP Port found")
	ERR_NO_SMTP_USER          = errors.New("No SMTP User found")
	ERR_NO_SMTP_PASSWORD      = errors.New("No SMTP Password found")
)

// LoadConfig loads the application configuration file
func (*App) LoadConfig(configFile string) error {
	err := godotenv.Load(configFile)
	if err != nil {
		return ERR_FAILED_READ_CONFIG
	}
	return nil
}

// ValidateConfig checks that we have configuration we can use in the application
func (a *App) ValidateConfig() error {
	// host
	if os.Getenv("MORPHEUS_API_HOST") == "" {
		return ERR_NO_API_HOST
	}
	a.Config.MorpheusHost = os.Getenv("MORPHEUS_API_HOST")

	// token
	if os.Getenv("MORPHEUS_API_BEARER_TOKEN") == "" {
		return ERR_NO_API_TOKEN
	}
	a.Config.MorpheusToken = os.Getenv("MORPHEUS_API_BEARER_TOKEN")

	// poll interval
	pollInt, err := strconv.Atoi(os.Getenv("POLL_INTERVAL"))
	if err != nil {
		return ERR_POLL_INTERVAL_NOT_INT
	}
	if pollInt == 0 {
		pollInt = 30
	}
	a.Config.PollInterval = pollInt

	// smtp server
	if os.Getenv("SMTP_SERVER") == "" {
		return ERR_NO_SMTP_SERVER
	}
	a.Config.SmtpServer = os.Getenv("SMTP_SERVER")

	// smtp port
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return ERR_NO_SMTP_PORT
	}
	if port == 0 {
		return ERR_NO_SMTP_PORT
	}
	a.Config.SmtpPort = port

	// smtp user
	if os.Getenv("SMTP_USER") == "" {
		return ERR_NO_SMTP_USER
	}
	a.Config.SmtpUser = os.Getenv("SMTP_USER")

	// smtp password
	if os.Getenv("SMTP_PASSWORD") == "" {
		return ERR_NO_SMTP_PASSWORD
	}
	a.Config.SmtpPassword = os.Getenv("SMTP_PASSWORD")

	return nil
}
