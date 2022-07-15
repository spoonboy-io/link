package internal_test

import (
	"os"
	"testing"

	"github.com/spoonboy-io/link/internal"
)

func createTestConfigFile(file string, config []byte, t *testing.T) {
	if err := os.WriteFile(file, config, 0644); err != nil {
		t.Fatalf("could not write test config file %+v", err)
	}
}

func removeTestConfigFileAndResetEnv(file string, t *testing.T) {
	if err := os.Remove(file); err != nil {
		t.Fatal("Could not remove test config file")
	}
	os.Clearenv()
}

func TestApp_LoadConfig(t *testing.T) {
	app := &internal.App{}
	wantErr := internal.ERR_FAILED_READ_CONFIG
	gotErr := app.LoadConfig("badConfigFile.env")

	if gotErr != wantErr {
		t.Errorf("wanted %v got %v", wantErr, gotErr)
	}
}

func TestApp_ValidateConfig(t *testing.T) {

	var testCases = []struct {
		name     string
		filename string
		config   []byte
		wantErr  error
	}{
		{
			name:     "all good, should pass",
			filename: "test1.env",
			config: []byte(`## Morpheus
MORPHEUS_API_HOST=https://testhost
MORPHEUS_API_BEARER_TOKEN=xxx-testtoken-xxx
POLL_INTERVAL=30

## SMTP
SMTP_SERVER=testmailserver.net
SMTP_PORT=587
SMTP_USER=testuser
SMTP_PASSWORD=testpassword
`),
			wantErr: nil,
		},

		{
			name:     "no host, should fail",
			filename: "test2.env",
			config: []byte(`## Morpheus
MORPHEUS_API_HOST=
MORPHEUS_API_BEARER_TOKEN=xxx-testtoken-xxx
POLL_INTERVAL=30

## SMTP
SMTP_SERVER=testmailserver.net
SMTP_PORT=587
SMTP_USER=testuser
SMTP_PASSWORD=testpassword
`),
			wantErr: internal.ERR_NO_API_HOST,
		},

		{
			name:     "no token, should fail",
			filename: "test3.env",
			config: []byte(`## Morpheus
MORPHEUS_API_HOST=https://testhost
MORPHEUS_API_BEARER_TOKEN=
POLL_INTERVAL=30

## SMTP
SMTP_SERVER=testmailserver.net
SMTP_PORT=587
SMTP_USER=testuser
SMTP_PASSWORD=testpassword
`),
			wantErr: internal.ERR_NO_API_TOKEN,
		},

		{
			name:     "bad poll data, should fail",
			filename: "test4.env",
			config: []byte(`## Morpheus
MORPHEUS_API_HOST=https://testhost
MORPHEUS_API_BEARER_TOKEN=xxx-testtoken-xxx
POLL_INTERVAL=30SECS(STRING)

## SMTP
SMTP_SERVER=testmailserver.net
SMTP_PORT=587
SMTP_USER=testuser
SMTP_PASSWORD=testpassword
`),
			wantErr: internal.ERR_POLL_INTERVAL_NOT_INT,
		},

		{
			name:     "no smtp server, should fail",
			filename: "test5.env",
			config: []byte(`## Morpheus
MORPHEUS_API_HOST=https://testhost
MORPHEUS_API_BEARER_TOKEN=xxx-testtoken-xxx
POLL_INTERVAL=30

## SMTP
SMTP_SERVER=
SMTP_PORT=587
SMTP_USER=testuser
SMTP_PASSWORD=testpassword
`),
			wantErr: internal.ERR_NO_SMTP_SERVER,
		},

		{
			name:     "no smtp port, should fail",
			filename: "test6.env",
			config: []byte(`## Morpheus
MORPHEUS_API_HOST=https://testhost
MORPHEUS_API_BEARER_TOKEN=xxx-testtoken-xxx
POLL_INTERVAL=30

## SMTP
SMTP_SERVER=testmailserver.net
SMTP_PORT=
SMTP_USER=testuser
SMTP_PASSWORD=testpassword
`),
			wantErr: internal.ERR_NO_SMTP_PORT,
		},

		{
			name:     "no smtp user, should fail",
			filename: "test7.env",
			config: []byte(`## Morpheus
MORPHEUS_API_HOST=https://testhost
MORPHEUS_API_BEARER_TOKEN=xxx-testtoken-xxx
POLL_INTERVAL=30

## SMTP
SMTP_SERVER=testmailserver.net
SMTP_PORT=587
SMTP_USER=
SMTP_PASSWORD=testpassword
`),
			wantErr: internal.ERR_NO_SMTP_USER,
		},

		{
			name:     "no smtp password, should fail",
			filename: "test8.env",
			config: []byte(`## Morpheus
MORPHEUS_API_HOST=https://testhost
MORPHEUS_API_BEARER_TOKEN=xxx-testtoken-xxx
POLL_INTERVAL=30

## SMTP
SMTP_SERVER=testmailserver.net
SMTP_PORT=587
SMTP_USER=testuser
SMTP_PASSWORD=
`),
			wantErr: internal.ERR_NO_SMTP_PASSWORD,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			createTestConfigFile(tc.filename, tc.config, t)
			app := &internal.App{}
			err := app.LoadConfig(tc.filename)
			if err != nil {
				t.Fatalf("Expected error %v", err)
			}
			gotErr := app.ValidateConfig()
			if gotErr != tc.wantErr {
				t.Errorf("wanted %v got %v", tc.wantErr, gotErr)
			}
			removeTestConfigFileAndResetEnv(tc.filename, t)
		})

	}
}
