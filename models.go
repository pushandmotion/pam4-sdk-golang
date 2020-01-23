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
	Name            string            `json:"name"`
	Alias           string            `json:"alias"`
	Description     string            `json:"description"`
	IsEnabled       bool              `json:"is_enabled"`
	Triggers        []*SegmentTrigger `json:"triggers"`
	TriggerExcludes []string          `json:"trigger_excludes"`
	DelayAmount     string            `json:"delay_amount"`
	DelayUnit       string            `json:"delay_unit"`
	// Operator        string            `json:"operator"`
	// Trigger         string            `json:"trigger"`
	// ID              string            `json:"id"`
	// IsCustom        bool              `json:"is_custom"`
	// CreateAt        *time.Time        `json:"created_at"`
	// UpdateAt        *time.Time        `json:"updated_at"`
	// Type            string            `json:"type"`
	// Excluders string `json:"excluders"`
}

// SegmentTrigger is segment trigger
type SegmentTrigger struct {
	Type       string             `json:"type"`
	Conditions []SegmentCondition `json:"conditions"`
}

// SegmentCondition is segment condition
type SegmentCondition struct {
	Operation string `json:"operation"`
	Trigger   string `json:"trigger"`
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
	segment.Description = ""
	for _, st := range segment.Triggers {
		for _, tc := range st.Conditions {
			tc.Trigger = uuid.New().String()
		}
	}
	return segment
}
