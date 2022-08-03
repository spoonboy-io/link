package main

/*

broadly this is what we need to build:-

configuration file & consumer, for connecting to API (path & key) and configuration for SMTP server - DONE
configuration  file (YAML) and consumer, for the approval routing logic - DONE
SMTP email routine with template for limited branding/customisation, including actions for view/approve/deny
HTTP server to handle the requests from the emails, (approve/deny/request more info) - with TLS. Providing detail and confirmation
client to poll the morpheus api for approvals, and make approve POST requests when constraints satisfied
state - we need to manage the approval state, while it is out for approval

*/

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/spoonboy-io/link/internal/morpheus"
	"github.com/spoonboy-io/link/internal/state"

	"github.com/spoonboy-io/link/internal/routes"

	"github.com/gorilla/mux"

	"github.com/spoonboy-io/koan"
	"github.com/spoonboy-io/link/internal"
	"github.com/spoonboy-io/link/internal/approval"
	"github.com/spoonboy-io/link/internal/certificate"
	"github.com/spoonboy-io/reprise"
)

var (
	version   = "Development build"
	goversion = "Unknown"
)

var logger *koan.Logger
var st *state.State
var app *internal.App

func init() {
	logger = &koan.Logger{}
	st = &state.State{}
	app = &internal.App{
		Logger: logger,
		State:  st,
	}

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
	if err := app.LoadConfig(internal.APP_CONFIG); err != nil {
		logger.FatalError("Failed to load application", err)
	}

	// validate it here rather than later
	if err := app.ValidateConfig(); err != nil {
		logger.FatalError("Application configuration is not sufficient", err)
	}

	// load approval YAML and validate we can use
	if err := approval.ReadAndParseConfig(internal.APPROVAL_CONFIG); err != nil {
		logger.FatalError("Failed to read approval configuration file", err)
	}

	if err := approval.ValidateConfig(); err != nil {
		logger.FatalError("Failed to validate approval configuration", err)
	}
}

// Shutdown runs on SIGINT and panic
func Shutdown(cancel context.CancelFunc) {
	fmt.Println("") // break after ^C
	logger.Warn("Application terminated")

	// cancel the context so we can stop our http client and current requests
	logger.Info("Cancelling HTTP client requests")
	cancel()

	logger.Info("Saving application state")
	if err := st.CreateAndWrite(); err != nil {
		logger.Error("Failed to save application state", err)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	app.Ctx = ctx
	defer Shutdown(cancel)

	// write console banner
	reprise.WriteSimple(&reprise.Banner{
		Name:         "Link",
		Description:  "Multi-person approval notifications for Morpheus",
		Version:      version,
		GoVersion:    goversion,
		WebsiteURL:   "https://spoonboy.io",
		VcsURL:       "https://github.com/spoonboy-io/link",
		VcsName:      "Github",
		EmailAddress: "hello@spoonboy.io",
	})

	// api poller which initiates most of the work
	go func() {
		pollInterval := time.NewTicker(time.Duration(app.Config.PollInterval) * time.Second)
		for range pollInterval.C {
			newApprovals, err := morpheus.CheckNewApprovals(ctx, app)
			if err != nil {
				logger.Error("Morpheus API request error", err)
			}
			lastCheckMsg := fmt.Sprintf("Checking for new Morpheus Approvals at %s", time.Now())
			logger.Info(lastCheckMsg)

			// match against the configuration
			_ = newApprovals // TODO
		}
	}()

	// handlers
	mux := mux.NewRouter()
	handler := &routes.Routes{
		App: app,
	}

	//mux.HandleFunc(`/`, handler.Ping).Methods("GET")
	mux.HandleFunc(`/ping`, handler.Ping).Methods("GET")

	// start HTTPS server
	go func() {
		hostPort := net.JoinHostPort(internal.SRV_HOST, internal.SRV_PORT)
		srvTLS := &http.Server{
			Addr:         hostPort,
			Handler:      mux,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 5 * time.Second,
		}

		logger.Info(fmt.Sprintf("Starting HTTPS server on %s", hostPort))
		if err := srvTLS.ListenAndServeTLS(fmt.Sprintf("%s/cert.pem", internal.TLS_FOLDER), fmt.Sprintf("%s/key.pem", internal.TLS_FOLDER)); err != nil {
			logger.FatalError("Failed to start HTTPS server", err)
		}
	}()

	// shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}
