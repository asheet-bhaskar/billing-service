package temporal

import (
	"context"

	"go.temporal.io/sdk/client"
)

type TemporalClient interface {
	ExecuteWorkflow(context.Context, client.StartWorkflowOptions, interface{}, ...interface{}) (client.WorkflowRun, error)
	SignalWorkflow(context.Context, string, string, string, interface{}) error
}

type temporalClient struct {
	client client.Client
}

func NewTemporalClient(c client.Client) TemporalClient {
	return &temporalClient{client: c}
}

func (t *temporalClient) ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error) {
	return t.client.ExecuteWorkflow(ctx, options, workflow, args...)
}

func (t *temporalClient) SignalWorkflow(ctx context.Context, workflowID, runID, signalName string, arg interface{}) error {
	return t.client.SignalWorkflow(ctx, workflowID, runID, signalName, arg)
}
