package morpheus

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spoonboy-io/link/internal"
)

// CheckNewApprovals obtains a list of approvals from the Morpheus API, we use 'offset' to
// get the new approvals since last checked
func CheckNewApprovals(ctx context.Context, app *internal.App) error {
	// form the request
	requestURI := fmt.Sprintf("%s/api/approvals?max=50&offset=%d", app.Config.MorpheusHost, app.State.LastPollId)
	req, err := http.NewRequest("GET", requestURI, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	// add the bearer token
	bearerToken := fmt.Sprintf("BEARER %s", app.Config.MorpheusToken)
	req.Header.Add("Authorization", bearerToken)

	// going to ignore TLS errors, which we may get from a morpheus appliance running
	// self-signed certificates
	tConf := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tConf,
	}

	// make the API call
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Bad response received from API (%d)", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// TODO
	// issue in that the approvals list, nor the approval tells us much about the scope
	// nor which approval policy generated the approval - if we had that we could inspect the policy for the above
	// so we will have interogate the instances or app which are subject to the approval
	// inspect, acquire further data, (enough to match on the approval routing policies

	// update state

	// debug
	fmt.Println(string(body))

	return nil
}
