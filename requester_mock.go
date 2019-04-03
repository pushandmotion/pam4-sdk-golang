package pam4sdk

import (
	"net/http"
	"time"

	"github.com/3dsinteractive/testify/mock"
)

// MockRequesterConfig is mock for producer config
type MockRequesterConfig struct {
	mock.Mock
}

// NewMockRequesterConfig return MockRequesterConfig
func NewMockRequesterConfig() *MockRequesterConfig {
	return &MockRequesterConfig{}
}

// Endpoint return endpoint
func (cfg *MockRequesterConfig) Endpoint() string {
	args := cfg.Called()
	return args.String(0)
}

// AppID return app id
func (cfg *MockRequesterConfig) AppID() string {
	args := cfg.Called()
	return args.String(0)
}

// Secret return app secret
func (cfg *MockRequesterConfig) Secret() string {
	args := cfg.Called()
	return args.String(0)
}

// AppIDHeaderKey is key for app id
func (cfg *MockRequesterConfig) AppIDHeaderKey() string {
	args := cfg.Called()
	return args.String(0)
}

// SecretHeaderKey is key for secret
func (cfg *MockRequesterConfig) SecretHeaderKey() string {
	args := cfg.Called()
	return args.String(0)
}

// Timeout is mock
func (cfg *MockRequesterConfig) Timeout() time.Duration {
	args := cfg.Called()
	return args.Get(0).(time.Duration)
}

// MockRequester is mock for producer
type MockRequester struct {
	mock.Mock
}

// NewMockRequester return mock producer
func NewMockRequester() *MockRequester {
	return &MockRequester{}
}

// Get make GET request
func (rqt *MockRequester) Get(path string, params map[string]string) (string, error) {
	args := rqt.Called(path, params)
	return args.String(0), args.Error(1)
}

// Post make POST request
func (rqt *MockRequester) Post(path string, params map[string]string) (string, error) {
	args := rqt.Called(path, params)
	return args.String(0), args.Error(1)
}

// PostJSON make POST request with json body
func (rqt *MockRequester) PostJSON(path string, body interface{}) (string, error) {
	args := rqt.Called(path, body)
	return args.String(0), args.Error(1)
}

func (rqt *MockRequester) PostJSONR(path string, body interface{}) (*http.Response, string, error) {
	args := rqt.Called(path, body)
	return args.Get(0).(*http.Response), args.String(1), args.Error(2)
}

func (rqt *MockRequester) PostJSONRH(path string, body interface{}, headers map[string]string) (*http.Response, string, error) {
	args := rqt.Called(path, body, headers)
	return args.Get(0).(*http.Response), args.String(1), args.Error(2)
}

func (rqt *MockRequester) PostJSONRHC(path string, body interface{}, headers map[string]string, cookies []*http.Cookie) (*http.Response, string, error) {
	args := rqt.Called(path, body, headers, cookies)
	return args.Get(0).(*http.Response), args.String(1), args.Error(2)
}

// PostFile is mock function
func (rqt *MockRequester) PostFile(path string, filePath string, postParam string) (string, error) {
	args := rqt.Called(path, filePath, postParam)
	return args.String(0), args.Error(1)
}

func (rqt *MockRequester) GetR(path string, params map[string]string) (*http.Response, string, error) {
	args := rqt.Called(path, params)
	return args.Get(0).(*http.Response), args.String(1), args.Error(2)
}

func (rqt *MockRequester) GetRH(path string, params map[string]string, headers map[string]string) (*http.Response, string, error) {
	args := rqt.Called(path, params, headers)
	return args.Get(0).(*http.Response), args.String(1), args.Error(2)
}

func (rqt *MockRequester) GetRHC(path string, params map[string]string, headers map[string]string, cookies []*http.Cookie) (*http.Response, string, error) {
	args := rqt.Called(path, params, headers, cookies)
	return args.Get(0).(*http.Response), args.String(1), args.Error(2)
}

func (rqt *MockRequester) PostR(path string, params map[string]string) (*http.Response, string, error) {
	args := rqt.Called(path, params)
	return args.Get(0).(*http.Response), args.String(1), args.Error(2)
}

func (rqt *MockRequester) PostRH(path string, params map[string]string, headers map[string]string) (*http.Response, string, error) {
	args := rqt.Called(path, params, headers)
	return args.Get(0).(*http.Response), args.String(1), args.Error(2)
}

func (rqt *MockRequester) PostRHC(path string, params map[string]string, headers map[string]string, cookies []*http.Cookie) (*http.Response, string, error) {
	args := rqt.Called(path, params, headers, cookies)
	return args.Get(0).(*http.Response), args.String(1), args.Error(2)
}

func (rqt *MockRequester) PostRaw(path string, data interface{}, headers map[string]string) (*http.Response, string, error) {
	args := rqt.Called(path, data, headers)
	return args.Get(0).(*http.Response), args.String(1), args.Error(2)
}
