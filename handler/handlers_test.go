package handler

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateProject(t *testing.T) {
	input := strings.NewReader("{\n    \"name\": \"Backup Workflow Resources\"\n    \"id\": \"project-001\",\n    \"resources\": [\n        {\n            \"type\": \"S3:Bucket\",\n            \"properties\": {\n               \"BucketName\": \"\"\n            }\n        }\n    ]\n}")
	// A request with a non-existant isdn
	req1, err := http.NewRequest("POST", "/create-project", input)
	if err != nil {
		t.Fatal(err)
	}

	rr1 := newRequestRecorder(req1, "POST", "/create-project", createProjectHandler)
	if rr1.Code != 422 {
		t.Fatalf("Expected response code to be 422, got %v", rr1.Code)
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
