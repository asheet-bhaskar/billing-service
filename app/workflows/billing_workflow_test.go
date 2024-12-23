package workflows

import (
	"testing"
	"time"

	"github.com/asheet-bhaskar/billing-service/app/models"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type BillingWorkflowTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *BillingWorkflowTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *BillingWorkflowTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *BillingWorkflowTestSuite) Test_RemoveLineItem() {
	lineItemSignal := LineItemSignal{
		BillID: "bill-id-01",
		ItemID: "item-id-01",
	}

	bill := models.Bill{}

	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow("REMOVE_BILL_ITEM_CHANNEL", lineItemSignal)
	}, time.Millisecond*2)

	s.env.ExecuteWorkflow(BillingWorkflow, &bill)

	s.True(s.env.IsWorkflowCompleted())
}

func (s *BillingWorkflowTestSuite) Test_AddLineItem() {
	lineItemSignal := LineItemSignal{
		BillID: "bill-id-01",
		ItemID: "item-id-01",
	}

	bill := models.Bill{}

	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow("Add_BILL_ITEM_CHANNEL", lineItemSignal)
	}, time.Millisecond*2)

	s.env.ExecuteWorkflow(BillingWorkflow, &bill)

	s.True(s.env.IsWorkflowCompleted())
}

func (s *BillingWorkflowTestSuite) Test_AddAndRemoveLineItems() {
	lineItemSignal := LineItemSignal{
		BillID: "bill-id-01",
		ItemID: "item-id-01",
	}

	bill := models.Bill{}

	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow("REMOVE_BILL_ITEM_CHANNEL", lineItemSignal)
	}, time.Millisecond*2)

	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow("REMOVE_BILL_ITEM_CHANNEL", lineItemSignal)
	}, time.Millisecond*2)

	s.env.ExecuteWorkflow(BillingWorkflow, &bill)

	s.True(s.env.IsWorkflowCompleted())
}

func TestBillingWorkflowTestSuite(t *testing.T) {
	suite.Run(t, new(BillingWorkflowTestSuite))
}
