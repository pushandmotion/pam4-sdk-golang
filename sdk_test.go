package pam4sdk

import (
	"net/http"
	"testing"

	"github.com/3dsinteractive/testify/assert"
	"github.com/3dsinteractive/testify/suite"
)

type SdkTestSute struct {
	suite.Suite
}

func TestSdkTestSute(t *testing.T) {
	suite.Run(t, new(SdkTestSute))
}

func (ts *SdkTestSute) TestSendEvent_GivenEventWithParams_ExpectEventSent() {
	mockRq := NewMockRequester()
	mockLogger := NewMockLogger()

	p := map[string]interface{}{
		"event":        "event_1234",
		"page_url":     "page url 1234",
		"page_title":   "page title 1234",
		"tags":         "tag1,tag2",
		"useragent":    "user agent",
		"querystring":  "querystring 1234",
		"utm_campaign": "utm campaign 1234",
		"utm_term":     "utm term 1234",
		"utm_content":  "utm content 1234",
		"utm_medium":   "utm medium 1234",
		"utm_source":   "utm source 1234",
		"ip_address":   "10.10.10.10",
		"form_fields": map[string]interface{}{
			"_campaign":     "campaign_123",
			"custom-field1": "custom-value1",
			"custom-field2": "custom-value2",
		},
	}
	c := []*http.Cookie{
		&http.Cookie{
			Name:  "contact_id",
			Value: "contact_123",
		},
	}

	response := &http.Response{}
	mockRq.On("PostJSONRHC", "/trackers/events", p, map[string]string(nil), c).Return(response, "response_1234", nil)

	sdk := NewSdkR(mockRq, mockLogger)
	tracker := &Tracker{
		Event:       "event_1234",
		PageURL:     "page url 1234",
		PageTitle:   "page title 1234",
		Tags:        "tag1,tag2",
		UserAgent:   "user agent",
		QueryString: "querystring 1234",
		UTMCampaign: "utm campaign 1234",
		UTMTerm:     "utm term 1234",
		UTMContent:  "utm content 1234",
		UTMMedium:   "utm medium 1234",
		UTMSource:   "utm source 1234",
		IPAddress:   "10.10.10.10",
		FormFields: map[string]interface{}{
			"custom-field1": "custom-value1",
			"custom-field2": "custom-value2",
		},
	}
	res, err := sdk.SendEvent("contact_123", "campaign_123", tracker)

	is := assert.New(ts.T())
	if is.NoError(err) {
		is.Equal("response_1234", res)
		mockRq.AssertExpectations(ts.T())
	}
}

func (ts *SdkTestSute) TestProductTrends_GiveLimit_ExpectRequestSentWithLimit() {
	mockRq := NewMockRequester()
	mockLogger := NewMockLogger()
	sdk := NewSdkR(mockRq, mockLogger)

	p := map[string]string{
		"limit": "400",
	}
	mockRq.On("Get", "/api/products/trends", p).Return("response_1234", nil)

	res, err := sdk.ProductTrends(400)

	is := assert.New(ts.T())
	if is.NoError(err) {
		is.Equal("response_1234", res)
		mockRq.AssertExpectations(ts.T())
	}
}
