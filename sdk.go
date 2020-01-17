package pam4sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ISdk is interface for PAM client
type ISdk interface {
	SendEvent(contactID string, campaignID string, tracker *Tracker) (string, error)
	ProductTrends(limit int) (string, error)
	ProductRecommends(aiID string, contactID string, productID int) (string, error)
	AppNotifications(contactID string, mediaAlias string, mediaValue string) (string, error)

	GetTriggersCount() (string, error)
	GetTriggers(q string, page int, limit int) (string, error)
	GetTriggersStats(triggerID []string) (string, error)

	CreateTrigger(body map[string]interface{}) (string, error)
}

// Sdk is struct for PAM client
type Sdk struct {
	rq     IRequester
	logger ILogger
}

// NewSdk create client using default requester and 10 seconds timeout
func NewSdk(baseURL string, appID string, appSecret string) *Sdk {
	return NewSdkT(baseURL, appID, appSecret, 10*time.Second)
}

// NewSdkT create client using default requester with specify timeout
func NewSdkT(baseURL string, appID string, appSecret string, requestTimeout time.Duration) *Sdk {
	config := NewCustomRequesterConfig(baseURL, "x-app-id", "x-secret", appID, appSecret, requestTimeout)
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

// SendEventTransaction post tracker event to PAM
func (sdk *Sdk) SendEventTransaction(contactID string, campaignID string, transactionID string, tracker *Tracker) (string, error) {
	return sdk.sendEvent(contactID, campaignID, transactionID, tracker)
}

// SendEvent post tracker event to PAM
func (sdk *Sdk) SendEvent(contactID string, campaignID string, tracker *Tracker) (string, error) {
	return sdk.sendEvent(contactID, campaignID, "", tracker)
}

// SendEvent post tracker event to PAM
func (sdk *Sdk) sendEvent(contactID string, campaignID string, transactionID string, tracker *Tracker) (string, error) {

	if tracker.FormFields == nil {
		tracker.FormFields = make(map[string]interface{})
	}
	if len(campaignID) > 0 {
		tracker.FormFields["_campaign"] = campaignID
	}
	if len(transactionID) > 0 {
		tracker.FormFields["_transaction_id"] = transactionID
	}

	js, _ := json.Marshal(tracker)
	p := map[string]interface{}{}
	json.Unmarshal([]byte(js), &p)

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

// ProductTrends return product trendings
func (sdk *Sdk) ProductTrends(limit int) (string, error) {
	p := map[string]string{}
	if limit > 0 {
		p["limit"] = fmt.Sprintf("%v", limit)
	}

	return sdk.rq.Get("/api/products/trends", p)
}

// ProductRecommends return product recommends
func (sdk *Sdk) ProductRecommends(aiID string, contactID string, productID int) (string, error) {
	p := map[string]string{}
	if len(contactID) > 0 {
		p["contact_id"] = fmt.Sprintf("%v", contactID)
	}

	if productID > 0 {
		p["productId"] = fmt.Sprintf("%v", productID)
	}

	productRecommendsPath := fmt.Sprintf("/api/ai/%s", aiID)

	return sdk.rq.Get(productRecommendsPath, p)
}

// AppNotifications return app notifications for given contactID, mediaAlias and mediaValue
func (sdk *Sdk) AppNotifications(contactID string, mediaAlias string, mediaValue string) (string, error) {
	p := map[string]string{}
	p["contact_id"] = contactID
	p["media_alias"] = mediaAlias
	p["media_value"] = mediaValue

	notificationPath := fmt.Sprintf("/api/app-notifications")

	return sdk.rq.Get(notificationPath, p)
}

// GetTriggersCount return number of triggers amount
func (sdk *Sdk) GetTriggersCount() (string, error) {

	countTriggers := fmt.Sprintf("/api/triggers/count")

	return sdk.rq.Get(countTriggers, nil)
}

// GetTriggers return list of triggers
func (sdk *Sdk) GetTriggers(q string, page int, limit int) (string, error) {
	p := map[string]string{}
	if len(q) > 0 {
		p["q"] = q
	}

	if page > 0 {
		p["page"] = fmt.Sprintf("%d", page)
	}

	if limit > 0 {
		p["limit"] = fmt.Sprintf("%d", limit)
	}

	triggers := fmt.Sprintf("/api/triggers")

	return sdk.rq.Get(triggers, p)
}

// GetTriggersStats return number of customer in triggers amount
func (sdk *Sdk) GetTriggersStats(triggerIDs []string) (string, error) {
	p := map[string]string{}
	if len(triggerIDs) > 0 {
		p["id"] = strings.Join(triggerIDs, ",")
	}

	triggersStat := fmt.Sprintf("/api/triggers/stat")

	return sdk.rq.Get(triggersStat, p)
}

// CreateTrigger create trigger
func (sdk *Sdk) CreateTrigger(body interface{}) (string, error) {
	createTrigger := fmt.Sprintf("/api/triggers")

	return sdk.rq.PostJSON(createTrigger, body)
}

// {
// name: "Testing"
// alias: "testing"x
// triggers: [{type: "IN_OR_NOT_IN_TRIGGERS",…}]
// trigger_excludes: []
// is_enabled: true
// delay_amount: ""
// delay_unit: ""
// }

// {"id":"f355d839-9fae-4352-a074-473ee8473c08","is_enabled":true,"created_at":"2020-01-17 04:07:36","updated_at":"2020-01-17 04:07:36","type":"SAVED","name":"Testing","alias":"testing","description":"","triggers":[{"type":"IN_OR_NOT_IN_TRIGGERS","conditions":[{"operator":"in","trigger":"2f7185ea-f6dd-4106-97ba-f15015afde55"}]}],"trigger_excludes":null,"is_custom":false,"delay_amount":"","delay_unit":"","excluders":null}
