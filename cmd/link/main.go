package main

import (
	"errors"
	"fmt"
	"github.com/spoonboy-io/link/internal"
	"github.com/spoonboy-io/link/internal/approval"
	"github.com/spoonboy-io/link/internal/certificate"
	"os"
	"path/filepath"
	"github.com/joho/godotenv"
	"github.com/spoonboy-io/koan"
)

var (
	version   = "Development build"
	goversion =	"Unknown"
)

var logger *koan.Logger

func init() {
	logger = &koan.Logger{}

	// check/create data folder
	templatesPath := filepath.Join(".", internal.TEMPLATE_FOLDER)
	if err := os.MkdirAll(templatesPath, os.ModePerm); err != nil {
		logger.FatalError("Problem checking/creating templates folder", err)
	}

	// create starter default email template if not exist
	defaultTemplate := fmt.Sprintf("%s/default.html", internal.TEMPLATE_FOLDER)
	if _, err := os.Stat(defaultTemplate); errors.Is(err, os.ErrNotExist) {
		logger.Info("Creating default email template")
		if err := os.WriteFile(defaultTemplate, []byte(internal.DefaultTemplate), 0644); err != nil {
			logger.FatalError("Problem creating the default email template", err)
		}
	}

	// check/create certificates folder
	tlsPath := filepath.Join(".", internal.TLS_FOLDER)
	if err := os.MkdirAll(tlsPath, os.ModePerm); err != nil {
		logger.FatalError("Problem checking/creating 'certificates' folder", err)
	}

	// add self-signed certificate only if folder empty, if the cert expires it
	// it can be deleted so the code here creates a new cert.pem and key.pem file
	cert := fmt.Sprintf("%s/cert.pem", internal.TLS_FOLDER)
	if _, err := os.Stat(cert); errors.Is(err, os.ErrNotExist) {
		logger.Info("Creating self-signed TLS certificate for the server")
		if err := certificate.Make(logger); err != nil {
			logger.FatalError("Problem creating the certificate/key", err)
		}
	}

	// load application config
	err := godotenv.Load(internal.APP_CONFIG)
	if err != nil {
		logger.FatalError("Failed to read application configuration file", err)
	}

	// load approval YAML and validate we can use
	err = approval.ReadAndParseConfig(internal.APPROVAL_CONFIG)
	if err != nil {
		logger.FatalError("Failed to read approval configuration file", err)
	}

	err = approval.ValidateConfig()
	if err != nil {
		logger.FatalError("Failed to validate approval configuration", err)
	}

}

func main(){
	fmt.Println("Hello World!")
}


/*

broadly this is what we need to build:-

configuration file & consumer, for connecting to API (path & key) and configuration for SMTP server
configuration  file (YAML) and consumer, for the approval routing logic
SMTP email routine with template for limited branding/customisation, including actions for view/approve/deny
HTTP server to handle the requests from the emails, (approve/deny/request more info) - with TLS. Providing detail and confirmation
client to poll the morpheus api for approvals, and make approve POST requests when constraints satisfied
state - we need to manage the approval state, while it is out for approval

 */
