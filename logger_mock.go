package pam4sdk

import "github.com/3dsinteractive/testify/mock"

// MockLogger is mock for logger
type MockLogger struct {
	mock.Mock
}

// NewMockLogger return MockLogger
func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

// Info ILogger implementation of Info
func (logger *MockLogger) Info(message string) {
	logger.Called(message)
}

// InfoFL ILogger implementation of Info
func (logger *MockLogger) InfoFL(message string, fn string, line int) {
	logger.Called(message)
}

// Warn ILogger implementation of Warn
func (logger *MockLogger) Warn(message string) {
	logger.Called(message)
}

// WarnFL ILogger implementation of Warn
func (logger *MockLogger) WarnFL(message string, fn string, line int) {
	logger.Called(message)
}

// Debug ILogger implementation of Debug
func (logger *MockLogger) Debug(message string) {
	// logger.Called(message)
}

// DebugFL ILogger implementation of Debug
func (logger *MockLogger) DebugFL(message string, fn string, line int) {
	// logger.Called(message)
}

// Error ILogger implementation of Error
func (logger *MockLogger) Error(message string) {
	logger.Called(message)
}

// ErrorFL ILogger implementation of Error
func (logger *MockLogger) ErrorFL(message string, fn string, line int) {
	logger.Called(message)
}

// Print ILogger implementation of Print
func (logger *MockLogger) Print(obj interface{}) {
	logger.Called(obj)
}
