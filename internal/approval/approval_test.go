package approval

import (
	"os"
	"reflect"
	"testing"
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
