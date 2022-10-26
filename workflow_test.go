package firstfailure

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *UnitTestSuite) Test_NilWorkflow() {
	var a *Activities

	// Everything returning nil should result in a successful nil
	s.env.OnActivity(a.CheckInventory, mock.Anything, mock.Anything).Return(nil)
	s.env.OnActivity(a.Charge, mock.Anything, mock.Anything).Return(nil)
	s.env.OnActivity(a.FulfillOrder, mock.Anything, mock.Anything).Return(nil)

	s.env.ExecuteWorkflow(Workflow, OrderInfo{})

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
	s.NoError(s.env.GetWorkflowResult(nil))
}

func (s *UnitTestSuite) Test_ChargeFail_InsufficientFunds() {
	var a *Activities

	s.env.OnActivity(a.CheckInventory, mock.Anything, mock.Anything).Return(nil)
	s.env.OnActivity(a.Charge, mock.Anything, mock.Anything).Times(1).Return(&InsufficientFundsError{})
	s.env.AssertNotCalled(s.T(), "FulfillOrder", mock.Anything, mock.Anything)

	//s.env.OnActivity(a.FulfillOrder, mock.Anything, mock.Anything).Panic("mock-panic")

	s.env.ExecuteWorkflow(Workflow, OrderInfo{})

	s.True(s.env.IsWorkflowCompleted())

	err := s.env.GetWorkflowError()
	s.Error(err)
	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("InsufficientFundsError", applicationErr.Type())
}

func (s *UnitTestSuite) Test_ChargeFail_403() {
	var a *Activities

	s.env.OnActivity(a.CheckInventory, mock.Anything, mock.Anything).Return(nil)
	s.env.OnActivity(a.Charge, mock.Anything, mock.Anything).Times(1).Return(
		temporal.NewNonRetryableApplicationError(
			"",
			"ConnectionError",
			&ConnectionError{code: 403, message: "Auth Fail"},
		))
	s.env.AssertNotCalled(s.T(), "FulfillOrder", mock.Anything, mock.Anything)

	//s.env.OnActivity(a.FulfillOrder, mock.Anything, mock.Anything).Panic("mock-panic")

	s.env.ExecuteWorkflow(Workflow, OrderInfo{})

	s.True(s.env.IsWorkflowCompleted())

	err := s.env.GetWorkflowError()
	s.Error(err)
	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("ConnectionError", applicationErr.Type())
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
