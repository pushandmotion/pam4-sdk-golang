// +build e2e

package pam4sdk

import (
	"fmt"
	"testing"

	"github.com/3dsinteractive/testify/assert"
	"github.com/3dsinteractive/testify/suite"
)

type SdkE2ETestSute struct {
	suite.Suite
}

func TestSdkE2ETestSute(t *testing.T) {
	suite.Run(t, new(SdkE2ETestSute))
}

func (ts *SdkE2ETestSute) TestSendEvent_GivenParams_ExpectSendToPAM() {
	sdk := NewSdk("https://connect.pushandmotion.com", "key", "secret")
	tracker := &Tracker{
		Event: "event-sample",
		FormFields: map[string]interface{}{
			"media-alias": "chaiyapong@3dsinteractive.com",
		},
	}
	res, err := sdk.SendEvent("contact_id", "campaign_id", tracker)
	is := assert.New(ts.T())

	if is.NoError(err) {
		fmt.Println(res)
	}
}

func (ts *SdkE2ETestSute) TestProductsTrends_GivenLimit_ExpectReturnProducts() {
	sdk := NewSdk("https://connect.pushandmotion.com", "key", "secret")
	res, err := sdk.ProductTrends(500)
	is := assert.New(ts.T())

	if is.NoError(err) {
		fmt.Println(res)
	}
}
