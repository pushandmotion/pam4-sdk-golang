package pam4sdk

import (
	"encoding/json"
	"time"
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
	Type       string        `json:"type"`
	Conditions []interface{} `json:"conditions"`
}

// Contact struct for information contact
type Contact struct {
	Attrs map[string]interface{} `json:"attrs"`
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
	return segment
}

// NewContact for parsing to segment struct
func NewContact(intf interface{}) *Contact {
	contact := &Contact{}
	bytes, err := json.Marshal(intf)
	if err != nil {
		return contact
	}
	err = json.Unmarshal(bytes, contact)
	if err != nil {
		return contact
	}
	return contact
}

// CampaignPostBody is information of body
type CampaignPostBody struct {
	IsEnabled        bool              `json:"is_enabled"`
	Alias            string            `json:"alias"`
	State            string            `json:"state"`
	Name             string            `json:"name"`
	NonExpired       bool              `json:"non_expired"`
	DatePushRanges   []*DatePushRanges `json:"date_push_ranges"`
	DateWorkingRange []string          `json:"date_working_range"`
}

// CampaignUpdateBody is information of update body
type CampaignUpdateBody struct {
	Alias              string            `json:"alias"`
	Name               string            `json:"name"`
	State              string            `json:"state"`
	IsEnabled          bool              `json:"is_enabled"`
	CampaignCategoryID string            `json:"campaign_category_id"`
	NonExpired         bool              `json:"non_expired"`
	EngageTagsAdd      []interface{}     `json:"engage_tags_add"`
	EngageTagsRemove   []interface{}     `json:"engage_tags_remove"`
	Tags               []interface{}     `json:"tags"`
	DatePushRanges     []*DatePushRanges `json:"date_push_ranges"`
	DateWorkingRange   []string          `json:"date_working_range"`
}

type DatePushRanges struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type CampaignTriger struct {
	IsCustom       bool     `json:"is_custom"`
	TriggerSavedID string   `json:"trigger_saved_id"`
	Triggers       Triggers `json:"triggers"`
}

type Triggers struct {
	DelayAmount     string        `json:"delay_amount"`
	DelayUnit       string        `json:"delay_unit"`
	Triggers        []interface{} `json:"triggers"`
	TriggerExcludes []string      `json:"trigger_excludes"`
	IsEnabled       bool          `json:"is_enabled"`
}
