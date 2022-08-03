package morpheus

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spoonboy-io/link/internal"
	"github.com/spoonboy-io/link/internal/approval"
)

// ApprovalsResponse stores the limited data we need to hold about the approvals at this stage
// the remainder we will find by making further requests for each approval id
type ApprovalsResponse struct {
	Approvals []approval.Approval `json:"approvals"`
}

type ApprovalResponse struct {
	Approval approval.Approval `json:"approval"`
}

const (
	STATUS_REQUESTED = "1 requested"
)

// CheckNewApprovals obtains a list of approvals from the Morpheus API, we use 'offset' to
// get the new approvals since last checked
func CheckNewApprovals(ctx context.Context, app *internal.App) ([]approval.Approval, error) {

	// approvalsRequested contains only the approvals request since last call which
	// we need to further inspect and match against approval policy logic
	var approvalsRequested []approval.Approval

	// form the request
	// the api does not support an index in this call, so we have to set a high max to get everything
	// and then process it
	requestURI := fmt.Sprintf("%s/api/approvals?max=10000", app.Config.MorpheusHost)
	req, err := http.NewRequest("GET", requestURI, nil)
	if err != nil {
		return approvalsRequested, err
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
		return approvalsRequested, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return approvalsRequested, fmt.Errorf("Bad response received from API (%d)", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return approvalsRequested, fmt.Errorf("Could not read response body", err)
	}

	// capture to struct
	approvalsRes := ApprovalsResponse{}
	if err := json.Unmarshal(body, &approvalsRes); err != nil {
		return approvalsRequested, fmt.Errorf("Could not unmarshal response body", err)
	}

	// TODO IN-PROGRESS
	// issue in that the approvals list, nor the approval tells us much about the scope
	// nor which approval policy generated the approval - if we had that we could inspect the policy for the above
	// so we will have interogate the instances or app which are subject to the approval
	// inspect, acquire further data, (enough to match on the approval routing policies

	for i := range approvalsRes.Approvals {
		// only process if > last poll id
		if approvalsRes.Approvals[i].Id > app.State.LastPollId {
			app.State.LastPollId = approvalsRes.Approvals[i].Id

			// only retrieve & process if in "requested" state
			if approvalsRes.Approvals[i].Status == STATUS_REQUESTED {
				//make a request for the approval/id endpoint
				approval, err := GetApproval(ctx, &approvalsRes.Approvals[i], app)
				if err != nil {
					return approvalsRequested, err
				}

				// append to slice of data to keep for post-processing
				approvalsRequested = append(approvalsRequested, approval)
			}
		}

	}

	// at this point we need to interogate the instance or app to determine scope
	for i := range approvalsRequested {
		//TODO
		fmt.Printf("%+v", approvalsRequested[i])

	}

	return approvalsRequested, nil
}

// GetApproval obtains information from the Morpheus API about the approval
// we will update the pointer so we are not returning approval as value
func GetApproval(ctx context.Context, approval *approval.Approval, app *internal.App) (approval.Approval, error) {

	approvalRes := ApprovalResponse{}

	// form the request
	requestURI := fmt.Sprintf("%s/api/approvals/%d", app.Config.MorpheusHost, approval.Id)
	req, err := http.NewRequest("GET", requestURI, nil)
	if err != nil {
		return approvalRes.Approval, err
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
		return approvalRes.Approval, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return approvalRes.Approval, fmt.Errorf("Bad response received from API (%d)", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return approvalRes.Approval, err
	}

	// capture to struct
	if err := json.Unmarshal(body, &approvalRes); err != nil {
		return approvalRes.Approval, err
	}

	return approvalRes.Approval, nil
}
