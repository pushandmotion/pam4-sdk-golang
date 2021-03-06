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

	// Segments
	GetSegmentsCount() (string, error)
	GetSegments(q string, page int, limit int) (string, error)
	GetSegmentsStats(segmentIDs []string) (string, error)
	GetSegmentByID(segmentID string) (string, error)

	CreateSegment(body *Segment) (string, error)
	UpdateSegment(segmentID string, body *Segment) (string, error)

	DeleteSegment(segmentID string) (string, error)

	// Campaigns
	CreateCampaign(body *CampaignPostBody) (string, error)
	UpdateCampaign(id string, body *CampaignUpdateBody) (string, error)
	GetCampaigns(q, aliases string, ids []string, page, limit string) (string, error)
	UpdateCampaignTrigger(id string, body *CampaignTriger) (string, error)
	GetCampaignsStats(campaignIDs []string) (string, error)
	GetCampaignDetail(campaignID string) (string, error)
	GetCampaignDetailByAlias(alias string) (string, error)
	GetCampaignReport(campaignID string) (string, error)
	DeleteCampaign(campaignID string) (string, error)
	GetMedia(isAll, isExcludeDisabled, MediaType string) (string, error)
	UpdateMessageSMS(campaignID string, body *UpdateMessageSMS) (*SMSMessageResponse, string, error)
	UpdateMessagePushNotification(campaignID string, body *UpdateMessagePushNotification) (*PushNotificationMessageResponse, string, error)

	// Contact
	CreateContact(file string, fieldMatch string, tags string) (string, error)
	CreateContactWithBody(body string) (string, error)
	UpdateContactAttr(contactID string, body *Contact) (string, error)
	GetContacts(q string, field string, page, limit string) (string, error)
	DeleteTagsByContacts(body *ContactsTags) (string, error)
	AddTagsByContacts(body *ContactsTags) (string, error)
	GetContactsTags(tags string, searchKeyword string, page, limit string) (string, error)
}

// Sdk is struct for PAM client
type Sdk struct {
	connect *RequestLogger
	cms     *RequestLogger
}

// RequestLogger is struct for request and logger
type RequestLogger struct {
	rq     IRequester
	logger ILogger
}

// NewSdk create client using default requester and 10 seconds timeout
func NewSdk(baseURL string, appID string, appSecret string) *Sdk {
	sdk := &SDKConnector{baseURL, appID, appSecret, 10 * time.Second}
	return NewSdkT(sdk, nil)
}

// NewSdkT create client using default requester with specify timeout
func NewSdkT(connectSDK, cmsSDK *SDKConnector) *Sdk {
	connectConfig := NewCustomRequesterConfig(
		connectSDK.BaseURL,
		"x-app-id",
		"x-secret",
		connectSDK.AppID,
		connectSDK.AppSecret,
		connectSDK.RequestTimeout)
	cmsConfig := NewCustomRequesterConfig(
		cmsSDK.BaseURL,
		"x-app-id",
		"x-secret",
		cmsSDK.AppID,
		cmsSDK.AppSecret,
		cmsSDK.RequestTimeout)
	logger := NewLoggerSimple()
	conr := NewRequester(connectConfig, logger)
	cmsr := NewRequester(cmsConfig, logger)
	return &Sdk{
		connect: &RequestLogger{rq: conr, logger: logger},
		cms:     &RequestLogger{rq: cmsr, logger: logger},
	}
}

// NewSdkR create new client with requester
func NewSdkR(conRL, cmsRL *RequestLogger) *Sdk {
	return &Sdk{conRL, cmsRL}
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

	sdkC := sdk.connect

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
	_, body, err := sdkC.rq.PostJSONRHC("/trackers/events", p, nil, c)

	if err != nil {
		return "", NewErrorE(sdkC.logger, err)
	}
	return body, nil
}

// ProductTrends return product trendings
func (sdk *Sdk) ProductTrends(limit int) (string, error) {
	sdkC := sdk.connect
	p := map[string]string{}
	if limit > 0 {
		p["limit"] = fmt.Sprintf("%v", limit)
	}

	return sdkC.rq.Get("/api/products/trends", p)
}

// ProductRecommends return product recommends
func (sdk *Sdk) ProductRecommends(aiID string, contactID string, productID int) (string, error) {
	sdkC := sdk.connect
	p := map[string]string{}
	if len(contactID) > 0 {
		p["contact_id"] = fmt.Sprintf("%v", contactID)
	}

	if productID > 0 {
		p["productId"] = fmt.Sprintf("%v", productID)
	}

	productRecommendsPath := fmt.Sprintf("/api/ai/%s", aiID)

	return sdkC.rq.Get(productRecommendsPath, p)
}

// AppNotifications return app notifications for given contactID, mediaAlias and mediaValue
func (sdk *Sdk) AppNotifications(contactID string, mediaAlias string, mediaValue string) (string, error) {
	sdkC := sdk.connect
	p := map[string]string{}
	p["contact_id"] = contactID
	p["media_alias"] = mediaAlias
	p["media_value"] = mediaValue

	notificationPath := fmt.Sprintf("/api/app-notifications")

	return sdkC.rq.Get(notificationPath, p)
}

// GetSegmentsCount return number of segments amount
func (sdk *Sdk) GetSegmentsCount() (string, error) {
	sdkC := sdk.cms
	countSegments := fmt.Sprintf("/triggers/count")

	return sdkC.rq.Get(countSegments, nil)
}

// GetSegments return list of segments
func (sdk *Sdk) GetSegments(q string, page int, limit int) (string, error) {
	sdkC := sdk.cms
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

	segments := fmt.Sprintf("/triggers")

	return sdkC.rq.Get(segments, p)
}

// GetSegmentsStats return number of customer in segments amount
func (sdk *Sdk) GetSegmentsStats(segmentIDs []string) (string, error) {
	sdkC := sdk.connect
	p := map[string]string{}
	if len(segmentIDs) > 0 {
		p["id"] = strings.Join(segmentIDs, ",")
	}

	segmentStat := fmt.Sprintf("/api/triggers/stat")

	return sdkC.rq.Get(segmentStat, p)
}

// GetSegmentByID return segment info by segment ID
func (sdk *Sdk) GetSegmentByID(segmentID string) (string, error) {
	sdkC := sdk.cms
	segmentByID := fmt.Sprintf("/triggers/%s", segmentID)

	return sdkC.rq.Get(segmentByID, nil)
}

// CreateSegment create segment
func (sdk *Sdk) CreateSegment(body *Segment) (string, error) {
	sdkC := sdk.cms
	createSegment := fmt.Sprintf("/triggers")

	return sdkC.rq.PostJSON(createSegment, body)
}

// UpdateSegment update segment by id
func (sdk *Sdk) UpdateSegment(segmentID string, body *Segment) (string, error) {
	sdkC := sdk.cms
	updateSegment := fmt.Sprintf("/triggers/%s", segmentID)

	return sdkC.rq.PutJSON(updateSegment, body)
}

// DeleteSegment delete segment by id
func (sdk *Sdk) DeleteSegment(segmentID string) (string, error) {
	sdkC := sdk.cms
	deleteSegment := fmt.Sprintf("/triggers/%s", segmentID)

	return sdkC.rq.Delete(deleteSegment, nil)
}

// CreateCampaign create campaign
func (sdk *Sdk) CreateCampaign(body *CampaignPostBody) (string, error) {
	sdkC := sdk.cms

	return sdkC.rq.PostJSON("/campaigns", body)
}

// UpdateCampaign update campaign by id
func (sdk *Sdk) UpdateCampaign(id string, body *CampaignUpdateBody) (string, error) {
	sdkC := sdk.cms
	endpoint := fmt.Sprintf("/campaigns/%s", id)

	return sdkC.rq.PutJSON(endpoint, body)
}

// GetCampaigns return list of campaigns
func (sdk *Sdk) GetCampaigns(q, aliases string, ids []string, page, limit string) (string, error) {
	sdkC := sdk.cms
	p := map[string]string{}
	if len(q) > 0 {
		p["q"] = q
	}

	if len(aliases) > 0 {
		p["aliases"] = aliases
	}

	if page != "" {
		p["page"] = page
	}

	if limit != "" {
		p["limit"] = limit
	}

	if len(ids) > 0 {
		p["ids"] = strings.Join(ids, ",")
	}

	campaigns := fmt.Sprintf("/campaigns")

	return sdkC.rq.Get(campaigns, p)
}

// UpdateCampaignTrigger update segment in campaign
func (sdk *Sdk) UpdateCampaignTrigger(id string, body *CampaignTriger) (string, error) {
	sdkC := sdk.cms
	endpoint := fmt.Sprintf("/campaigns/%s/triggers", id)

	return sdkC.rq.PutJSON(endpoint, body)
}

// GetCampaignsStats return number of campaign in campaigns amount
func (sdk *Sdk) GetCampaignsStats(campaignIDs []string) (string, error) {
	sdkC := sdk.connect
	p := map[string]string{}
	if len(campaignIDs) > 0 {
		p["id"] = strings.Join(campaignIDs, ",")
	}

	campaignStat := fmt.Sprintf("/api/campaigns/stat")

	return sdkC.rq.Get(campaignStat, p)
}

// GetCampaignDetail return detail of Campaign
func (sdk *Sdk) GetCampaignDetail(campaignID string) (string, error) {
	sdkC := sdk.cms
	campaigns := fmt.Sprintf("/campaigns/%s", campaignID)

	return sdkC.rq.Get(campaigns, nil)
}

// GetCampaignDetailByAlias return detail of Campaign
func (sdk *Sdk) GetCampaignDetailByAlias(alias string) (string, error) {
	sdkC := sdk.cms
	campaigns := fmt.Sprintf("/campaigns/aliases/%s", alias)

	return sdkC.rq.Get(campaigns, nil)
}

// GetCampaignReport return report of Campaign
func (sdk *Sdk) GetCampaignReport(campaignID string) (string, error) {
	sdkC := sdk.connect
	campaigns := fmt.Sprintf("/api/reports/campaigns/%s", campaignID)

	return sdkC.rq.Get(campaigns, nil)
}

// DeleteCampaign delete campaign by id
func (sdk *Sdk) DeleteCampaign(campaignID string) (string, error) {
	sdkC := sdk.cms
	endpoint := fmt.Sprintf("/campaigns/%s", campaignID)

	return sdkC.rq.Delete(endpoint, nil)
}

// CreateContact return nil when create success
func (sdk *Sdk) CreateContact(filePath, attrs, tags string) (string, error) {
	sdkC := sdk.connect
	extraData := fmt.Sprintf(`attrs=%s&&tags=%s`, attrs, tags)

	return sdkC.rq.PostFile("/api/contacts/upload", filePath, "file", extraData)
}

// CreateContactWithBody is create contact api
func (sdk *Sdk) CreateContactWithBody(body string) (string, error) {
	sdkC := sdk.connect

	return sdkC.rq.PostJSON("/api/contacts", body)
}

// UpdateContactAttr return contact information when update success
func (sdk *Sdk) UpdateContactAttr(contactID string, body *Contact) (string, error) {
	sdkC := sdk.connect
	updateContact := fmt.Sprintf("/api/contacts/%s", contactID)

	return sdkC.rq.PutJSON(updateContact, body)
}

// GetContacts return contact list
func (sdk *Sdk) GetContacts(searchKeyword string, field, page, limit string) (string, error) {
	sdkC := sdk.connect
	params := map[string]string{
		"q":     searchKeyword,
		"field": field,
		"page":  page,
		"limit": limit,
	}

	return sdkC.rq.Get("/api/contacts", params)
}

// AddTagsByContacts add tag in old contact
func (sdk *Sdk) AddTagsByContacts(body *ContactsTags) (string, error) {
	sdkC := sdk.connect

	return sdkC.rq.PostJSON("/api/contacts/tags", body)
}

// DeleteTagsByContacts return tags available
func (sdk *Sdk) DeleteTagsByContacts(body *ContactsTags) (string, error) {
	sdkC := sdk.connect

	return sdkC.rq.DeleteJSON("/api/contacts/tags", body)
}

// GetMedia return media list
func (sdk *Sdk) GetMedia(isAll, isExcludeDisabled, MediaType string) (string, error) {
	sdkC := sdk.cms
	params := map[string]string{
		"is_all":           isAll,
		"exclude_disabled": isExcludeDisabled,
		"type":             MediaType,
	}

	return sdkC.rq.Get("/media", params)
}

// UpdateMessageSMS update message by media type
func (sdk *Sdk) UpdateMessageSMS(campaignID string, body *UpdateMessageSMS) (*SMSMessageResponse, string, error) {
	sdkC := sdk.cms
	endpoint := fmt.Sprintf("/campaigns/%s/messages/sms", campaignID)

	resultStr, err := sdkC.rq.PutJSON(endpoint, body)

	if err != nil {
		return &SMSMessageResponse{}, "", err
	}

	res := &SMSMessageResponse{}
	err = json.Unmarshal([]byte(resultStr), res)

	if err != nil {
		return &SMSMessageResponse{}, "", err
	}

	return res, resultStr, nil
}

// UpdateMessagePushNotification message by media type
func (sdk *Sdk) UpdateMessagePushNotification(
	campaignID string,
	body *UpdateMessagePushNotification,
) (*PushNotificationMessageResponse, string, error) {
	sdkC := sdk.cms
	endpoint := fmt.Sprintf("/campaigns/%s/messages/mobile_notification", campaignID)

	resultStr, err := sdkC.rq.PutJSON(endpoint, body)

	if err != nil {
		return &PushNotificationMessageResponse{}, "", err
	}

	res := &PushNotificationMessageResponse{}
	err = json.Unmarshal([]byte(resultStr), res)

	if err != nil {
		return &PushNotificationMessageResponse{}, "", err
	}

	return res, resultStr, nil
}

// GetContactsTags return contact list
func (sdk *Sdk) GetContactsTags(tags string, searchKeyword string, page, limit string) (string, error) {
	sdkC := sdk.connect
	params := map[string]string{
		"q":     searchKeyword,
		"page":  page,
		"limit": limit,
		"tags":  tags,
	}

	return sdkC.rq.Get("/api/contacts/tag/multiple", params)
}
