package approval

import (
	_ "fmt"
	"io/ioutil"
	_ "net/url"

	_ "github.com/spoonboy-io/link/internal"

	"gopkg.in/yaml.v2"
)

var config ApprovalsConfig

// ApprovalsConfig is a representation of the parsed YAML approvals.yaml configuration file
type ApprovalsConfig []struct {
	ApprovalConfig `yaml:"approval"`
}

// ApprovalConfig represents a single approval configuration
type ApprovalConfig struct {
	Description 	string  `yaml:"description"`
	OnProvision		bool 	`yaml:"onProvision"`
	OnDelete 		bool 	`yaml:"onDelete"`
	OnReconfigure 	bool 	`yaml:"onReconfigure"`
	LinkApproval	bool	`yaml:"linkApproval"`
	Scope    		Scope 	`yaml:"scope"`
}

// Scope represents the scope configuration options which can be set in the YAML.
// Default scope is 'global' unless overridden here - by a single setting here
type Scope struct {
	Group     	string	`yaml:"group"`
	Cloud 		string 	`yaml:"cloud"`
	User    	string 	`yaml:"user"`
	Role  		string	`yaml:"role"`
	Network   	string 	`yaml:"network"`
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

// ValidateConfig will check that the config parsed can used by application
func ValidateConfig() error {
	return nil
}