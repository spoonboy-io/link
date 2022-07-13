package approval

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spoonboy-io/link/internal"
)

var testYamlFile = "test_approvals.yaml"

func writeTestYamlFile(t *testing.T) {
	data := `---
- approval:
    description: test approval config 1
    onProvision: true
    onDelete: true
    onReconfigure: false
    linkedApproval: true
    recipientList:
        - ollie@test.io
        - test@test.io
    scope:
        group: All Clouds

- approval:
    description: test approval config 2
    template: my_custom_template.html
    onReconfigure: true
    recipientList:
        - ollie@test.io`

	if err := os.WriteFile(testYamlFile, []byte(data), 0644); err != nil {
		t.Fatalf("could not write test yaml file %+v", err)
	}
}

func removeTestYamlFile(t *testing.T) {
	if err := os.Remove(testYamlFile); err != nil {
		t.Fatal("Could not remove test yaml file")
	}
}

func makeTestTemplate(t *testing.T) {
	templatesPath := filepath.Join(".", internal.TEMPLATE_FOLDER)
	if err := os.MkdirAll(templatesPath, os.ModePerm); err != nil {
		t.Fatal("Problem checking/creating test template folder", err)
	}

	testTemplate := fmt.Sprintf("%s/test.html", internal.TEMPLATE_FOLDER)
	if _, err := os.Stat(testTemplate); errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile(testTemplate, []byte(internal.DefaultTemplate), 0644); err != nil {
			t.Fatal("Problem creating the test email template", err)
		}
	}
}

func removeTestTemplate(t *testing.T) {
	templatesPath := filepath.Join(".", internal.TEMPLATE_FOLDER)
	if err := os.RemoveAll(templatesPath); err != nil {
		t.Fatal("Problem removing template folder", err)
	}
}

func TestReadAndParseConfig(t *testing.T) {
	// write test yaml config
	writeTestYamlFile(t)

	wantConfig := ApprovalsConfig{
		{
			ApprovalConfig{
				Description:    "test approval config 1",
				OnProvision:    true,
				OnDelete:       true,
				OnReconfigure:  false,
				LinkedApproval: true,
				RecipientList:  []string{"ollie@test.io", "test@test.io"},
				Scope: Scope{
					Group: "All Clouds",
				},
			},
		},
		{
			ApprovalConfig{
				Description:   "test approval config 2",
				TemplateFile:  "my_custom_template.html",
				OnReconfigure: true,
				RecipientList: []string{"ollie@test.io"},
			},
		},
	}

	if err := ReadAndParseConfig(testYamlFile); err != nil {
		t.Fatalf("could not read test yaml file %+v", err)
	}

	gotConfig := config

	if !reflect.DeepEqual(gotConfig, wantConfig) {
		t.Errorf("\n\nWanted\n%v\n, \n\ngot \n%v\n", wantConfig, gotConfig)
	}

	// clean up
	removeTestYamlFile(t)
}

func TestValidateConfig(t *testing.T) {

	makeTestTemplate(t)

	testCases := []struct {
		name    string
		config  ApprovalsConfig
		wantErr error
	}{
		{
			name: "all good, should pass",
			config: ApprovalsConfig{
				{
					ApprovalConfig{
						Description:   "test approval config 1",
						OnProvision:   true,
						RecipientList: []string{"test@test.com"},
						TemplateFile:  "test.html",
					},
				},
			},
			wantErr: nil,
		},
	}

	/*
		OnDelete       bool     `yaml:"onDelete"`
		OnReconfigure  bool     `yaml:"onReconfigure"`
		LinkedApproval bool     `yaml:"linkedApproval"`
		RecipientList  []string `yaml:"recipientList"`
		Scope          Scope    `yaml:"scope"`
	*/

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			// set the package config
			config = tc.config
			gotErr := ValidateConfig()
			if gotErr != tc.wantErr {
				t.Errorf("wanted %v got %v", tc.wantErr, gotErr)
			}
		})
	}

	removeTestTemplate(t)
}
