package approval

import (
	"errors"
	_ "fmt"
	"io/ioutil"
	"net/mail"
	_ "net/url"

	_ "github.com/spoonboy-io/link/internal"

	"gopkg.in/yaml.v2"
)

var config ApprovalsConfig

var (
	ERR_NO_DESCRIPTION      = errors.New("No description is set")
	ERR_NO_ACTION           = errors.New("Approval is not configured for 'provision', 'delete' nor 'reconfigure'")
	ERR_NO_RECIPIENTS       = errors.New("No recipients are configured")
	ERR_BAD_RECIPIENT_EMAIL = errors.New("Recipient email address appears to be invalid")
	ERR_TEMPLATE_NOT_EXIST  = errors.New("Configured mesage template cannot be found")

	/*
		ERR_BAD_METHOD                  = errors.New("method is not acceptable")
		ERR_BAD_URL                     = errors.New("url is appears to be invalid")
		ERR_NO_BODY                     = errors.New("method requires requestBody")
		ERR_NO_TRIGGER                  = errors.New("No triggers defined in the hook")
		ERR_BAD_STATUS_TRIGGER          = errors.New("Trigger set on status is not recognised")
		ERR_NO_EXECUTING_STATUS_TRIGGER = errors.New("Can not trigger on status 'executing'")
		ERR_NOT_HTTPS                   = errors.New("url is not secure (no HTTPS)")
		ERR_COULD_NOT_PARSE_BODY        = errors.New("Problem parsing request body, check included variables")

	*/
)

// ApprovalsConfig is a representation of the parsed YAML approvals.yaml configuration file
type ApprovalsConfig []struct {
	ApprovalConfig `yaml:"approval"`
}

// ApprovalConfig represents a single approval configuration
type ApprovalConfig struct {
	Description    string   `yaml:"description"`
	TemplateFile   string   `yaml:"template"`
	OnProvision    bool     `yaml:"onProvision"`
	OnDelete       bool     `yaml:"onDelete"`
	OnReconfigure  bool     `yaml:"onReconfigure"`
	LinkedApproval bool     `yaml:"linkedApproval"`
	RecipientList  []string `yaml:"recipientList"`
	Scope          Scope    `yaml:"scope"`
}

// Scope represents the scope configuration options which can be set in the YAML.
// Default scope is 'global' unless overridden here - by a single setting here
type Scope struct {
	Group   string `yaml:"group"`
	Cloud   string `yaml:"cloud"`
	User    string `yaml:"user"`
	Role    string `yaml:"role"`
	Network string `yaml:"network"`
}

// ReadAndParseConfig reads the contents of the YAML approvals config filer
// and parses it to a map of Approval structs
func ReadAndParseConfig(cfgFile string) error {
	yamlConfig, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(yamlConfig, &config); err != nil {
		return err
	}

	return nil
}

// ValidateConfig will check that the config parsed can be used by application
func ValidateConfig() error {
	for i := range config {
		// check description
		if config[i].Description == "" {
			return ERR_NO_DESCRIPTION
		}

		// check has action
		if !config[i].OnProvision && !config[i].OnDelete && config[i].OnReconfigure {
			return ERR_NO_ACTION
		}

		// TODO if template configured check it exists

		// check at least one recipient
		if len(config[i].RecipientList) == 0 {
			return ERR_NO_RECIPIENTS
		}

		// check recipient email addresses seem valid
		for _, email := range config[i].RecipientList {
			if _, err := mail.ParseAddress(email); err != nil {
				return ERR_BAD_RECIPIENT_EMAIL
			}
		}

	}

	return nil
}
