package approval

import (
	"errors"
	"fmt"
	"net/mail"
	"os"
	"time"

	"github.com/spoonboy-io/link/internal"

	"gopkg.in/yaml.v2"
)

var config ApprovalsConfig

var (
	ERR_NO_DESCRIPTION      = errors.New("No description is set")
	ERR_NO_ACTION           = errors.New("Approval is not configured for 'provision', 'delete' nor 'reconfigure'")
	ERR_NO_RECIPIENTS       = errors.New("No recipients are configured")
	ERR_BAD_RECIPIENT_EMAIL = errors.New("Recipient email address appears to be invalid")
	ERR_TEMPLATE_NOT_EXIST  = errors.New("Configured message template cannot be found")
	ERR_MULTIPLE_SCOPES     = errors.New("Multiple Scopes found")
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

// hold an approval
type Approval struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	RequestType string    `json:"requestType"`
	Status      string    `json:"status"`
	DateCreated time.Time `json:"datedCreated"`
	RequestBy   string    `json:"requestBy"`
	Items       []Item    `json:"approvalItems"`
	Scope       Scope     `json:"scope"`
}

type Item struct {
	Id int `json:"id"`
}

// ReadAndParseConfig reads the contents of the YAML approvals config filer
// and parses it to a map of Approval structs
func ReadAndParseConfig(cfgFile string) error {
	yamlConfig, err := os.ReadFile(cfgFile)
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
		if !config[i].OnProvision && !config[i].OnDelete && !config[i].OnReconfigure {
			return ERR_NO_ACTION
		}

		// if template configured check it exists
		if config[i].TemplateFile != "" {
			tmplFile := fmt.Sprintf("%s/%s", internal.TEMPLATE_FOLDER, config[i].TemplateFile)
			if _, err := os.Stat(tmplFile); errors.Is(err, os.ErrNotExist) {
				return ERR_TEMPLATE_NOT_EXIST
			}
		}

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

		// if scope is set we need to further validate that
		scope := config[i].Scope
		if scope != (Scope{}) {
			var set bool
			if scope.Group != "" {
				set = true
			}

			if scope.Cloud != "" {
				if set {
					return ERR_MULTIPLE_SCOPES
				}
				set = true
			}

			if scope.User != "" {
				if set {
					return ERR_MULTIPLE_SCOPES
				}
				set = true
			}

			if scope.Role != "" {
				if set {
					return ERR_MULTIPLE_SCOPES
				}
				set = true
			}

			if scope.Network != "" {
				if set {
					return ERR_MULTIPLE_SCOPES
				}
				set = true
			}
		}
	}

	return nil
}
