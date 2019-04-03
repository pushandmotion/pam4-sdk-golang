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
