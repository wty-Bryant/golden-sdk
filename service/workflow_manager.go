package service

import (
	"archive/zip"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type WorkflowManager struct {
	workflows map[string]*Workflow
	triggers  map[string]*Trigger
}

func (wm *WorkflowManager) CreateWorkflow(ctx context.Context, input *CreateWorkflowInput) (CreateWorkflowOutput, error) {
	workflow := &Workflow{
		ID:         input.ID,
		Name:       input.Name,
		Status:     input.Status,
		Components: input.Components,
		Variables:  input.Variables,
	}
	if wm.workflows == nil {
		wm.workflows = make(map[string]*Workflow)
	}
	wm.workflows[input.ID] = workflow

	out := CreateWorkflowOutput{
		Metadata: make(map[string]ComponentMetadata),
	}
	for _, c := range input.Components {
		cm := ComponentMetadata{
			ID:   c.ID,
			Type: ComponentType(c.Type),
			Next: c.Next,
		}
		out.Metadata[c.ID] = cm
	}

	return out, nil
}

type CreateWorkflowInput struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	ProjectID  string           `json:"projectId"`
	Inputs     []WorkflowInput  `json:"input"`
	Variables  []Variable       `json:"variables"`
	Components []ComponentInfo  `json:"steps"`
	Status     Status           `json:"status"`
	Output     []WorkflowOutput `json:"output"`
}

type WorkflowInput struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type WorkflowOutput struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Variable struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	DefaultValue string `json:"defaultValue"`
}

type CreateWorkflowOutput struct {
	Metadata map[string]ComponentMetadata `json:"component-metadata"`
}

type ComponentInfo struct {
	ID      string     `json:"id"`
	Type    string     `json:"type"`
	Inputs  []Variable `json:"input"`
	Outputs []Variable `json:"output"`
	Next    string     `json:"next"` // next component id
}

func (wm *WorkflowManager) createWorkflowComponents(ctx context.Context, workflowID string) (map[string]Component, error) {
	componentsInfo := wm.workflows[workflowID].Components
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
	workflows := make([]Workflow, 0)
	for _, workflow := range wm.workflows {
		workflows = append(workflows, *workflow)
	}

	return ListWorkflowsOutput{
		Workflows: workflows,
	}, nil
}

type ListWorkflowsOutput struct {
	Workflows []Workflow `json:"workflows"`
}

func (wm *WorkflowManager) CreateWorkflowTrigger(ctx context.Context, input *CreateWorkflowTriggerInput) error {
	trigger := &Trigger{
		input: *input,
	}

	if wm.triggers == nil {
		wm.triggers = make(map[string]*Trigger)
	}
	wm.triggers[input.ID] = trigger
	return nil
}

type CreateWorkflowTriggerInput struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Config     map[string]interface{} `json:"config"`
	WorkflowID string                 `json:"workflowId"`
	Input      string                 `json:"input"`
	Status     Status                 `json:"status"`
}

type Trigger struct {
	input CreateWorkflowTriggerInput
}

func (wm *WorkflowManager) ListWorkflowTriggers() (ListWorkflowTriggersOutput, error) {
	triggers := make([]CreateWorkflowTriggerInput, 0)
	for _, t := range wm.triggers {
		triggers = append(triggers, t.input)
	}

	return ListWorkflowTriggersOutput{
		Triggers: triggers,
	}, nil
}

type ListWorkflowTriggersOutput struct {
	Triggers []CreateWorkflowTriggerInput `json:"triggers"`
}

type Workflow struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Endpoint   string          `json:"endpoint"`
	Status     Status          `json:"status"`
	Components []ComponentInfo // ID to Component
	Variables  []Variable
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

type ReadFileOutput struct {
	files []string
}

func (c *ComponentReadFile) ID() string { return c.id }

func (c *ComponentReadFile) Do(ctx context.Context, input interface{}) (output interface{}, err error) {
	in, ok := input.(ReadFileInput)
	if !ok {
		return nil, fmt.Errorf("failed to cast ReadFileInput")
	}
	return c.do(ctx, in)
}

func (c *ComponentReadFile) do(ctx context.Context, input ReadFileInput) (ReadFileOutput, error) {
	files := make([]string, 0)
	err := filepath.WalkDir(input.directory, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return ReadFileOutput{}, err
	}
	return ReadFileOutput{
		files: files,
	}, nil
}

// file zipper component
type ComponentZipFile struct {
	id    string
	files []string
	next  string
}

type ZipFileInput struct {
	files   []string
	zipFile string
}

type ZipFileOutput struct {
	zipFile string
}

func (c *ComponentZipFile) ID() string { return c.id }

func (c *ComponentZipFile) Do(ctx context.Context, input interface{}) (output interface{}, err error) {
	in, ok := input.(ZipFileInput)
	if !ok {
		return nil, fmt.Errorf("failed to cast ZipFileInput")
	}
	return c.do(ctx, in)
}

func (c *ComponentZipFile) do(ctx context.Context, input ZipFileInput) (output ZipFileOutput, err error) {
	// Create a new zip archive.
	zipFile, err := os.Create(input.zipFile)
	if err != nil {
		return ZipFileOutput{}, fmt.Errorf("error when creating zip file %s: %v", input.zipFile, err)
	}
	defer zipFile.Close()

	// Create a new writer for the zip archive.
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range input.files {
		var fileInfo os.FileInfo
		var header *zip.FileHeader
		var writer io.Writer
		// Open the file to be added to the archive.
		fileToZip, err := os.Open(file)
		defer fileToZip.Close()
		if err != nil {
			return ZipFileOutput{}, fmt.Errorf("error when opening file %s: %v", file, err)
		}

		// Create a new file header for the file.
		fileInfo, err = fileToZip.Stat()
		if err != nil {
			return ZipFileOutput{}, fmt.Errorf("error when stating file %s: %v", file, err)
		}

		header, err = zip.FileInfoHeader(fileInfo)
		if err != nil {
			return ZipFileOutput{}, fmt.Errorf("error when building info header of file %s: %v", file, err)
		}
		// Set the file header name to the name of the file.
		header.Name = file

		// Add the file header to the zip archive.
		writer, err = zipWriter.CreateHeader(header)
		if err != nil {
			return ZipFileOutput{}, fmt.Errorf("error when creating zip writer of file %s: %v", file, err)
		}

		// Write the file contents to the zip archive.
		_, err = io.Copy(writer, fileToZip)
		if err != nil {
			return ZipFileOutput{}, fmt.Errorf("error when zipping content of file %s: %v", file, err)
		}
	}

	return ZipFileOutput{
		zipFile: input.zipFile,
	}, nil
}

// s3 PutObject component
type ComponentPutObject struct {
	id     string
	client s3.Client
	next   string
}

type PutObjectInput struct {
	bucket string
	region string
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
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(input.region),
	)
	if err != nil {
		return PutObjectOutput{}, fmt.Errorf("failed to load aws config: %v", err)
	}

	client := s3.NewFromConfig(cfg)
	for _, file := range input.files {
		body, err := os.Open(file)
		if err != nil {
			return PutObjectOutput{}, fmt.Errorf("error when opening file %s: %v", file, err)
		}
		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(input.bucket),
			Key:    aws.String(file),
			Body:   body,
		})
		if err != nil {
			return PutObjectOutput{}, fmt.Errorf("error when put object %s: %v", file, err)
		}
	}
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
