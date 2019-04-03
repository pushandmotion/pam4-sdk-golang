package pam4sdk

import (
	"encoding/json"
	"net/http"
	"time"
)

// ISdk is interface for PAM client
type ISdk interface {
	SendEvent(contactID string, campaignID string, tracker *Tracker) (string, error)
}

// Sdk is struct for PAM client
type Sdk struct {
	rq     IRequester
	logger ILogger
}

// NewSdk create client using default requester
func NewSdk(baseURL string, appID string, appSecret string) *Sdk {
	timeout := 3 * time.Second
	config := NewCustomRequesterConfig(baseURL, "x-app-id", "x-secret", appID, appSecret, timeout)
	logger := NewLoggerSimple()
	r := NewRequester(config, logger)
	return &Sdk{
		rq:     r,
		logger: logger,
	}
}

// NewSdkR create new client with requester
func NewSdkR(rq IRequester, logger ILogger) *Sdk {
	return &Sdk{
		rq:     rq,
		logger: logger,
	}
}

// SendEvent post tracker event to PAM
func (sdk *Sdk) SendEvent(contactID string, campaignID string, tracker *Tracker) (string, error) {

	js, _ := json.Marshal(tracker)
	p := map[string]interface{}{}
	json.Unmarshal([]byte(js), &p)
	if len(campaignID) > 0 {
		p["_campaign_id"] = campaignID
	}

	c := []*http.Cookie{
		&http.Cookie{
			Name:  "contact_id",
			Value: contactID,
		},
	}
	_, body, err := sdk.rq.PostJSONRHC("/trackers/events", p, nil, c)

	if err != nil {
		return "", NewErrorE(sdk.logger, err)
	}
	return body, nil
}
