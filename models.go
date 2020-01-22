package pam4sdk

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// SDKConnector is the information that using for requestor
type SDKConnector struct {
	BaseURL        string
	AppID          string
	AppSecret      string
	RequestTimeout time.Duration
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

// Segment struct for information segment
type Segment struct {
	ID          string     `json:"id"`
	IsEnabled   bool       `json:"is_enabled"`
	CreateAt    *time.Time `json:"created_at"`
	UpdateAt    *time.Time `json:"updated_at"`
	Type        string     `json:"type"`
	Name        string     `json:"name"`
	Alias       string     `json:"alias"`
	Description string     `json:"description"`
}

// NewSegment for parsing to segment struct
func NewSegment(intf interface{}) *Segment {
	segment := &Segment{}
	bytes, err := json.Marshal(intf)
	if err != nil {
		return segment
	}
	err = json.Unmarshal(bytes, segment)
	if err != nil {
		return segment
	}
	segment.ID = uuid.New().String()
	return segment
}
