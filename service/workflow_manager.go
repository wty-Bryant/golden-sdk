package service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type WorkflowManager struct {
	Workflows map[string]*Workflow
}

func (wm *WorkflowManager) CreateWorkflow(ctx context.Context, input *CreateWorkflowInput) (CreateWorkflowOutput, error) {
	var out CreateWorkflowOutput
	return out, nil
}

type CreateWorkflowInput struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Input      map[string]interface{} `json:"input"`
	Variables  map[string]interface{} `json:"variables"`
	Components []ComponentInfo        `json:"steps"`
	Status     Status                 `json:"status"`
	Output     map[string]interface{} `json:"output"`
}

type CreateWorkflowOutput struct {
	Metadata map[string]ComponentMetadata `json:"component-metadata"`
}

type ComponentInfo struct {
	ID     string              `json:"id"`
	Type   string              `json:"type"`
	Inputs []map[string]string `json:"parameters"`
	Next   string              `json:"next"` // next component id
}

func (wm *WorkflowManager) createWorkflowComponents(ctx context.Context, workflowID string) (map[string]Component, error) {
	componentsInfo := wm.Workflows[workflowID].Components
	components := make(map[string]Component, 0)
	for _, ci := range componentsInfo {
		switch ci.Type {
		case "S3:PutObject":
			components[ci.ID] = &ComponentPutObject{
				next: ci.Next,
			}
		case "ReadFile":
			components[ci.ID] = &ComponentReadFile{
				next: ci.Next,
			}
		case "ZipFile":
			components[ci.ID] = &ComponentZipFile{
				next: ci.Next,
			}
		case "HandleError":
			components[ci.ID] = &ComponentHandleError{
				next: ci.Next,
			}
		default:
		}
	}
	return components, nil
}

func (wm *WorkflowManager) ListWorkflows() (ListWorkflowsOutput, error) {
	var out ListWorkflowsOutput
	return out, nil
}

type ListWorkflowsOutput struct {
	workflows []Workflow `json:"workflows"`
}

func (wm *WorkflowManager) CreateWorkflowTrigger(ctx context.Context, input *CreateWorkflowTriggerInput) error {
	return nil
}

type CreateWorkflowTriggerInput struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Type       string                 `json:"workflow_trigger_type"`
	Config     map[string]interface{} `json:"trigger_conf"`
	WorkflowID string                 `json:"workflow_id"`
	Input      string                 `json:"input"`
	Status     Status                 `json:"status"`
}

type Trigger struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"workflow_trigger_type"`
	Status Status `json:"status"`
}

func (wm *WorkflowManager) ListWorkflowTriggers() (ListWorkflowTriggersOutput, error) {
	var out ListWorkflowTriggersOutput
	return out, nil
}

type ListWorkflowTriggersOutput struct {
	triggers []Trigger `json:"triggers"`
}

type Workflow struct {
	ID         string                   `json:"id"`
	Name       string                   `json:"name"`
	Endpoint   string                   `json:"endpoint"`
	Status     Status                   `json:"status"`
	Components map[string]ComponentInfo // ID to Component
}

type ComponentMetadata struct {
	ID   string        `json:"id"`
	Type ComponentType `json:"component-type"`
	Next string        `json:"next-id"`
}

type Component interface {
	ID() string
	Do(ctx context.Context, input interface{}) (output interface{}, err error)
}

// file system component
type ComponentReadFile struct {
	id        string
	directory string
	next      string
}

type ReadFileInput struct {
	directory string
}

type ReadFileOutput struct{}

func (c *ComponentReadFile) ID() string { return c.id }

func (c *ComponentReadFile) Do(ctx context.Context, input interface{}) (output interface{}, err error) {
	in, ok := input.(ReadFileInput)
	if !ok {
		return nil, fmt.Errorf("failed to cast ReadFileInput")
	}
	return c.do(ctx, in)
}

func (c *ComponentReadFile) do(ctx context.Context, input ReadFileInput) (output ReadFileOutput, err error) {
	return ReadFileOutput{}, nil
}

// file zipper component
type ComponentZipFile struct {
	id    string
	files []string
	next  string
}

type ZipFileInput struct {
	files []string
}

type ZipFileOutput struct{}

func (c *ComponentZipFile) ID() string { return c.id }

func (c *ComponentZipFile) Do(ctx context.Context, input interface{}) (output interface{}, err error) {
	in, ok := input.(ZipFileInput)
	if !ok {
		return nil, fmt.Errorf("failed to cast ZipFileInput")
	}
	return c.do(ctx, in)
}

func (c *ComponentZipFile) do(ctx context.Context, input ZipFileInput) (output ZipFileOutput, err error) {
	return ZipFileOutput{}, nil
}

// s3 PutObject component
type ComponentPutObject struct {
	id     string
	client s3.Client
	next   string
}

type PutObjectInput struct {
	bucket string
	files  []string
}

type PutObjectOutput struct{}

func (c *ComponentPutObject) ID() string { return c.id }

func (c *ComponentPutObject) Do(ctx context.Context, input interface{}) (output interface{}, err error) {
	in, ok := input.(PutObjectInput)
	if !ok {
		return nil, fmt.Errorf("failed to cast PutObjectInput")
	}
	return c.do(ctx, in)
}

func (c *ComponentPutObject) do(ctx context.Context, input PutObjectInput) (output PutObjectOutput, err error) {
	return PutObjectOutput{}, nil
}

// error handler component
type ComponentHandleError struct {
	id     string
	client cloudwatch.Client
	next   string
}

type ErrorHandleInput struct {
	err error
}

type ErrorHandleOutput struct{}

func (c *ComponentHandleError) ID() string { return c.id }

func (c *ComponentHandleError) Do(ctx context.Context, input interface{}) (output interface{}, err error) {
	in, ok := input.(ErrorHandleInput)
	if !ok {
		return nil, fmt.Errorf("failed to cast PutObjectInput")
	}
	return c.do(ctx, in)
}

func (c *ComponentHandleError) do(ctx context.Context, input ErrorHandleInput) (output ErrorHandleOutput, err error) {
	return ErrorHandleOutput{}, nil
}
