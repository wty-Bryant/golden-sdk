package service

type ResourceType string

const (
	Bucket ResourceType = "s3:bucket"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusDeactive Status = "deactive"
)

type ComponentType string

const (
	ComponentTypePutObject ComponentType = "s3:putObject"
)

type WorkflowTriggerType string

const (
	WorkflowTriggerTypeScheduled WorkflowTriggerType = "scheduled"
)
