package pam4sdk

import (
	"fmt"
	"testing"

	"github.com/3dsinteractive/testify/assert"
	"github.com/3dsinteractive/testify/suite"
)

type ErrorTestSuite struct {
	suite.Suite
}

func (ts *ErrorTestSuite) TestNewErrorE_GivenError_MustReturnNewError() {
	logger := NewMockLogger()
	logger.On("ErrorFL", "Error 1")

	err := fmt.Errorf("Error 1")
	newErr := NewErrorE(logger, err)
	logger.AssertExpectations(ts.T())

	assert.Equal(ts.T(), err.Error(), newErr.Error())
	assert.Equal(ts.T(), err, newErr.err)
}

func (ts *ErrorTestSuite) makeNewErrorE(logger ILogger, err error) error {
	return NewErrorE(logger, err)
}

func (ts *ErrorTestSuite) TestNewErrorE_GivenSameErrorTwice_MustReturnFirstError() {
	logger := NewMockLogger()
	logger.On("ErrorFL", "Error 1")

	err := fmt.Errorf("Error 1")
	newErr := ts.makeNewErrorE(logger, err)
	newErr2 := NewErrorE(logger, newErr)

	// newErr and newErr2 must be same object
	assert.Equal(ts.T(), newErr, newErr2)
	// logger Error must be call only once
	logger.AssertNumberOfCalls(ts.T(), "ErrorFL", 1)
	logger.AssertExpectations(ts.T())
}

func (ts *ErrorTestSuite) TestNewErrorM_GiveMessage_MustReturnNewError() {
	logger := NewMockLogger()
	logger.On("ErrorFL", "Error 1")

	newErr := NewErrorM(logger, "Error 1")
	assert.Equal(ts.T(), "Error 1", newErr.Error())
	assert.Nil(ts.T(), newErr.err)
	logger.AssertExpectations(ts.T())
}

func TestErrorTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorTestSuite))
}
