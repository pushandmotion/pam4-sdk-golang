package pam4sdk

import "github.com/3dsinteractive/testify/mock"

// MockSdk is mock for PAM sdk
type MockSdk struct {
	mock.Mock
}

// NewMockSdk return new mock
func NewMockSdk() *MockSdk {
	return &MockSdk{}
}

// SendEvent is mock
func (sdk *MockSdk) SendEvent(contactID string, campaignID string, tracker *Tracker) (string, error) {
	args := sdk.Called(contactID, campaignID, tracker)
	return args.String(0), args.Error(1)
}

// ProductTrends is mock
func (sdk *MockSdk) ProductTrends(limit int) (string, error) {
	args := sdk.Called(limit)
	return args.String(0), args.Error(1)
}

// ProductRecommends is mock
func (sdk *MockSdk) ProductRecommends(aiID string, contactID string, productID int) (string, error) {
	args := sdk.Called(aiID, contactID, productID)
	return args.String(0), args.Error(1)
}

// AppNotifications is mock
func (sdk *MockSdk) AppNotifications(contactID string, mediaAlias string, mediaValue string) (string, error) {
	args := sdk.Called(contactID, mediaAlias, mediaValue)
	return args.String(0), args.Error(1)
}
