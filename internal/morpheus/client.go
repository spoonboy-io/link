package morpheus

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/spoonboy-io/link/internal"
)

func MakeApprovalCheck(ctx context.Context, app *internal.App) error {
	var data io.Reader // TODO we will need to track approvals already seen

	req, err := http.NewRequest("POST", app.Config.MorpheusHost, data)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	// add the bearer token
	req.Header.Add("Authorization", app.Config.MorpheusToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Bad response received from API (%d)", res.StatusCode)
	}

	return nil
}
