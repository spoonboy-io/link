package internal

import "time"

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
