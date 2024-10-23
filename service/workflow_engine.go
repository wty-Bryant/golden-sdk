package service

import "context"

type WorkflowEngine struct{}

func (we *WorkflowEngine) RunWorkflow(ctx context.Context, input RunWorkflowInput) error {
	return nil
}

type RunWorkflowInput struct {
	ID    string                 `json:"workflowId"`
	Input map[string]interface{} `json:"input"`
}

type WorkflowTrigger struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Type       WorkflowTriggerType    `json:"workflow_trigger_type"`
	Config     map[string]interface{} `json:"trigger_conf"`
	WorkflowID string                 `json:"workflow_id"`
	Input      string                 `json:"input"`
	Status     Status                 `json:"status"`
}
