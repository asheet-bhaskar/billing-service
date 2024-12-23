package temporal

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/client"
)

type MockTemporalClient struct {
	mock.Mock
}

func (m *MockTemporalClient) ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error) {
	m.Called(ctx, options, workflow, args)
	return nil, nil
}

func (m *MockTemporalClient) SignalWorkflow(ctx context.Context, workflowID, runID, signalName string, arg interface{}) error {
	argsList := m.Called(ctx, workflowID, runID, signalName, arg)
	return argsList.Error(0)
}
