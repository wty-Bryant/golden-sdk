package service

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
)

type ResourceManager struct {
	projects map[string]Project
}

// CreateProject stores resources metadata into database
func (rm *ResourceManager) CreateProject(ctx context.Context, input *CreateProjectInput) error {
	if rm.projects == nil {
		rm.projects = make(map[string]Project)
	}

	rm.projects[input.ID] = Project{
		Name:      input.Name,
		ID:        input.ID,
		Resources: input.Resources,
	}

	return nil
}

// CreateProjectResources is called to initialize resources when workflow is triggerred
func (rm *ResourceManager) CreateProjectResources(ctx context.Context, input *CreateProjectResourcesInput) error {
	resources := rm.projects[input.ProjectID].Resources
	for _, r := range resources {
		switch r.Type {
		case "S3:Bucket":
			_, err := rm.createBucket(ctx, bucketInput{
				bucket: r.Properties["BucketName"],
				region: r.Properties["Region"],
			})
			if err != nil {
				return err
			}
		default:
		}
	}
	return nil
}

type CreateProjectResourcesInput struct {
	ProjectID string `json:"ProjectId"`
}

// ListProjects lists all projects info
func (rm *ResourceManager) ListProjects(ctx context.Context) ([]Project, error) {
	projects := make([]Project, 0)
	for _, p := range rm.projects {
		projects = append(projects, p)
	}
	return projects, nil
}

// createBucket creates s3 bucket resources
func (rm *ResourceManager) createBucket(ctx context.Context, input bucketInput) (*ResourceMetadata, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(input.region))
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(cfg)
	_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(input.bucket),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(input.region),
		},
	})
	if err != nil {
		return nil, err
	}
	return &ResourceMetadata{
		Type: Bucket,
		Name: input.bucket,
	}, nil
}

type Project struct {
	Name      string     `json:"name"`
	ID        string     `json:"id"`
	Resources []Resource `json:"resources"`
}

type CreateProjectInput struct {
	Name      string     `json:"name"`
	ID        string     `json:"id"`
	Resources []Resource `json:"resources"`
}

type Resource struct {
	Type       string            `json:"type"`
	Properties map[string]string `json:"properties"`
}

type bucketInput struct {
	bucket string
	region string
}

type ResourceMetadata struct {
	ID   string
	Name string
	Type ResourceType
	ARN  string
}
