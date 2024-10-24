package handler

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateProject(t *testing.T) {
	input := strings.NewReader("{\n    \"name\": \"Backup Workflow Resources\",\n    \"id\": \"project-001\",\n    \"resources\": [\n        {\n            \"type\": \"S3:Bucket\",\n            \"properties\": {\n               \"BucketName\": \"resource-bucket\",\n	 \"Region\": \"us-west-2\"\n           }\n        }\n    ]\n}")
	req1, err := http.NewRequest("POST", "/create-project", input)
	if err != nil {
		t.Fatal(err)
	}

	rr1 := newRequestRecorder(req1, "POST", "/create-project", createProjectHandler)
	if rr1.Code != 200 {
		t.Fatalf("Expected response code to be 200, got %v", rr1.Code)
	}
}

func TestListProjects(t *testing.T) {
	createProjectInput := strings.NewReader("{\n    \"name\": \"Backup Workflow Resources\",\n    \"id\": \"project-001\",\n    \"resources\": [\n        {\n            \"type\": \"S3:Bucket\",\n            \"properties\": {\n               \"BucketName\": \"resource-bucket\",\n	 \"Region\": \"us-west-2\"\n           }\n        }\n    ]\n}")
	createProjectReq, err := http.NewRequest("POST", "/create-project", createProjectInput)
	if err != nil {
		t.Fatal(err)
	}

	createProjectRR := newRequestRecorder(createProjectReq, "POST", "/create-project", createProjectHandler)
	if createProjectRR.Code != 200 {
		t.Fatalf("Expected response code to be 200, got %v", createProjectRR.Code)
	}

	listProjectsInput := strings.NewReader("{}")
	listProjectsReq, err := http.NewRequest("GET", "/list-projects", listProjectsInput)
	if err != nil {
		t.Fatal(err)
	}

	rr1 := newRequestRecorder(listProjectsReq, "GET", "/list-projects", listProjectsHandler)
	if rr1.Code != 200 {
		t.Fatalf("Expected response code to be 200, got %v", rr1.Code)
	}
}

func TestCreateProjectResources(t *testing.T) {
	// create a existing bucket causing error
	in := strings.NewReader("{\n    \"name\": \"Backup Workflow Resources\",\n    \"id\": \"project-001\",\n    \"resources\": [\n        {\n            \"type\": \"S3:Bucket\",\n            \"properties\": {\n               \"BucketName\": \"resource-bucket\",\n	 \"Region\": \"us-west-2\"\n           }\n        }\n    ]\n}")
	req, err := http.NewRequest("POST", "/create-project", in)
	if err != nil {
		t.Fatal(err)
	}

	rr := newRequestRecorder(req, "POST", "/create-project", createProjectHandler)
	if rr.Code != 200 {
		t.Fatalf("Expected response code to be 200, got %v", rr.Code)
	}

	input := strings.NewReader("{\n    \"projectId\": \"project-001\"\n}")
	req1, err := http.NewRequest("POST", "/create-project-resources", input)
	if err != nil {
		t.Fatal(err)
	}

	rr1 := newRequestRecorder(req1, "POST", "/create-project-resources", createProjectResourcesHandler)
	if rr1.Code != 400 {
		t.Fatalf("Expected response code to be 400, got %v", rr1.Code)
	}
}

func TestCreateWorkflow(t *testing.T) {
	in := strings.NewReader("{\n    \"id\": \"workflow_backup\",\n    \"name\": \"Backup Service\",\n    \"input\": {\n        \"path\": \"/tmp/backup/\"\n    },\n    \"variables\": {\n        \"listOfFiles\": {\n            \"type\": \"listOfString\"\n        }\n    },\n    \"steps\": [\n        {\n            \"id\": \"step-1\",\n            \"parameters\": {\n               \"bucket_name\": {\n                  \"type\": \"string\"\n               }\n            },\n            \"workflow_step_type\": \"S3:PutObject\",\n            \"next\": \"step-2\"\n        }\n    ],\n    \"status\": \"Active\",\n    \"output\": {\n       \n    }\n}")
	req, err := http.NewRequest("POST", "/create-workflow", in)
	if err != nil {
		t.Fatal(err)
	}

	rr := newRequestRecorder(req, "POST", "/create-workflow", createWorkflowHandler)
	if rr.Code != 200 {
		t.Fatalf("Expected response code to be 200, got %v", rr.Code)
	}
}

func TestListWorkflows(t *testing.T) {
	createWorkflowInput := strings.NewReader("{\n    \"id\": \"workflow_backup\",\n    \"name\": \"Backup Service\",\n    \"input\": {\n        \"path\": \"/tmp/backup/\"\n    },\n    \"variables\": {\n        \"listOfFiles\": {\n            \"type\": \"listOfString\"\n        }\n    },\n    \"steps\": [\n        {\n            \"id\": \"step-1\",\n            \"parameters\": {\n               \"bucket_name\": {\n                  \"type\": \"string\"\n               }\n            },\n            \"workflow_step_type\": \"S3:PutObject\",\n            \"next\": \"step-2\"\n        }\n    ],\n    \"status\": \"Active\",\n    \"output\": {\n       \n    }\n}")
	createWorkflowReq, err := http.NewRequest("POST", "/create-workflow", createWorkflowInput)
	if err != nil {
		t.Fatal(err)
	}

	createWorkflowRR := newRequestRecorder(createWorkflowReq, "POST", "/create-workflow", createWorkflowHandler)
	if createWorkflowRR.Code != 200 {
		t.Fatalf("Expected response code to be 200, got %v", createWorkflowRR.Code)
	}

	listWorkflowsInput := strings.NewReader("{}")
	listWorkflowReq, err := http.NewRequest("GET", "/list-workflows", listWorkflowsInput)
	if err != nil {
		t.Fatal(err)
	}

	listWorkflowsRR := newRequestRecorder(listWorkflowReq, "GET", "/list-workflows", listWorkflowsHandler)
	if listWorkflowsRR.Code != 200 {
		t.Fatalf("Expected response code to be 200, got %v", listWorkflowsRR.Code)
	}
}

func TestCreateWorkflowTrigger(t *testing.T) {
	in := strings.NewReader("{\n    \"id\": \"trigger-01\",\n    \"name\": \"Scheduled-Backup-Service\",\n    \"workflow_trigger_type\": \"scheduled\",\n    \"trigger_conf\": {\n        \"runAt\": \"Sunday 12, 2024\",\n        \"repeat\": \"True\"\n    },\n    \"workflow_id\": \"workflow-01\",\n    \"input\": \"/tmp/user1/backlup\",\n    \"status\": \"Active\"\n}")
	req, err := http.NewRequest("POST", "/create-workflow-trigger", in)
	if err != nil {
		t.Fatal(err)
	}

	rr := newRequestRecorder(req, "POST", "/create-workflow-trigger", createWorkflowTriggerHandler)
	if rr.Code != 200 {
		t.Fatalf("Expected response code to be 200, got %v", rr.Code)
	}
}

func TestListWorkflowTriggers(t *testing.T) {
	createWorkflowTriggerInput := strings.NewReader("{\n    \"id\": \"trigger-01\",\n    \"name\": \"Scheduled-Backup-Service\",\n    \"workflow_trigger_type\": \"scheduled\",\n    \"trigger_conf\": {\n        \"runAt\": \"Sunday 12, 2024\",\n        \"repeat\": \"True\"\n    },\n    \"workflow_id\": \"workflow-01\",\n    \"input\": \"/tmp/user1/backlup\",\n    \"status\": \"Active\"\n}")
	createWorkflowTriggerReq, err := http.NewRequest("POST", "/create-workflow-trigger", createWorkflowTriggerInput)
	if err != nil {
		t.Fatal(err)
	}

	createWorkflowTriggerRR := newRequestRecorder(createWorkflowTriggerReq, "POST", "/create-workflow-trigger", createWorkflowTriggerHandler)
	if createWorkflowTriggerRR.Code != 200 {
		t.Fatalf("Expected response code to be 200, got %v", createWorkflowTriggerRR.Code)
	}

	listWorkflowTriggersInput := strings.NewReader("{}")
	listWorkflowTriggersReq, err := http.NewRequest("GET", "/list-workflow-triggers", listWorkflowTriggersInput)
	if err != nil {
		t.Fatal(err)
	}

	listWorkflowTriggersRR := newRequestRecorder(listWorkflowTriggersReq, "GET", "/list-workflow-triggers", listWorkflowTriggersHandler)
	if listWorkflowTriggersRR.Code != 200 {
		t.Fatalf("Expected response code to be 200, got %v", listWorkflowTriggersRR.Code)
	}
}

// Mocks a handler and returns a httptest.ResponseRecorder
func newRequestRecorder(req *http.Request, method string, strPath string, fnHandler func(w http.ResponseWriter, r *http.Request, param httprouter.Params)) *httptest.ResponseRecorder {
	router := httprouter.New()
	router.Handle(method, strPath, fnHandler)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	router.ServeHTTP(rr, req)
	return rr
}
