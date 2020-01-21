package pam4sdk

import "time"

// SDKConnector is the information that using for requestor
type SDKConnector struct {
	baseURL        string
	appID          string
	appSecret      string
	requestTimeout time.Duration
}

// Tracker is the information that can send to PAM
type Tracker struct {
	Event       string                 `json:"event"`
	PageTitle   string                 `json:"page_title"`
	PageURL     string                 `json:"page_url"`
	Tags        string                 `json:"tags"`
	UserAgent   string                 `json:"useragent"`
	QueryString string                 `json:"querystring"`
	UTMCampaign string                 `json:"utm_campaign"`
	UTMTerm     string                 `json:"utm_term"`
	UTMContent  string                 `json:"utm_content"`
	UTMMedium   string                 `json:"utm_medium"`
	UTMSource   string                 `json:"utm_source"`
	IPAddress   string                 `json:"ip_address"`
	FormFields  map[string]interface{} `json:"form_fields"`
}
