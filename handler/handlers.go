package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golden-sdk/service"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func init() {
	rm = service.ResourceManager{}
	wm = service.WorkflowManager{}
}

var rm service.ResourceManager
var wm service.WorkflowManager

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// A Logger function which simply wraps the handler function around some log messages
func logger(fn func(w http.ResponseWriter, r *http.Request, param httprouter.Params)) func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		start := time.Now()
		log.Printf("%s %s", r.Method, r.URL.Path)
		fn(w, r, param)
		log.Printf("Done in %v (%s %s)", time.Since(start), r.Method, r.URL.Path)
	}
}

func createProjectHandler(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	input := &service.CreateProjectInput{}
	if err := populateModelFromHandler(w, r, param, input); err != nil {
		writeErrorResponse(w, http.StatusUnprocessableEntity, "Unprocessable create project input")
		return
	}

	resourceManager := &rm
	_ = resourceManager
	if err := rm.CreateProject(context.Background(), input); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "failed to create project")
		return
	}
	writeOKResponse(w, *input)
}

func createProjectResourcesHandler(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	input := &service.CreateProjectResourcesInput{}
	if err := populateModelFromHandler(w, r, param, input); err != nil {
		writeErrorResponse(w, http.StatusUnprocessableEntity, "Unprocessable create project resources input")
		return
	}

	if err := rm.CreateProjectResources(context.Background(), input); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "failed to create project resources")
		return
	}
	writeOKResponse(w, *input)
}

func listProjectsHandler(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	projects, err := rm.ListProjects(context.Background())
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "failed to list projects")
		return
	}
	writeOKResponse(w, projects)
}

func createWorkflowHandler(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	input := &service.CreateWorkflowInput{}
	if err := populateModelFromHandler(w, r, param, input); err != nil {
		writeErrorResponse(w, http.StatusUnprocessableEntity, "Unprocessable create workflow input")
		return
	}

	output, err := wm.CreateWorkflow(context.Background(), input)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "failed to create workflow")
		return
	}
	writeOKResponse(w, output)
}

func listWorkflowsHandler(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	out, err := wm.ListWorkflows()
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "failed to list workflows")
		return
	}

	writeOKResponse(w, out)
}

func createWorkflowTriggerHandler(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	input := &service.CreateWorkflowTriggerInput{}
	if err := populateModelFromHandler(w, r, param, input); err != nil {
		writeErrorResponse(w, http.StatusUnprocessableEntity, "Unprocessable create workflow trigger input")
		return
	}

	err := wm.CreateWorkflowTrigger(context.Background(), input)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "failed to create workflow trigger")
		return
	}
	writeOKResponse(w, *input)
}

func listWorkflowTriggersHandler(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	out, err := wm.ListWorkflowTriggers()
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "failed to list workflow triggers")
		return
	}

	writeOKResponse(w, out)
}

// Writes the response as a standard JSON response with StatusOK
func writeOKResponse(w http.ResponseWriter, m interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&JsonResponse{Data: m}); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
	}
}

// Writes the error response as a Standard API JSON response with a response code
func writeErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errorCode)
	json.
		NewEncoder(w).
		Encode(&JsonErrorResponse{Error: &ApiError{Status: errorCode, Title: errorMsg}})
}

func populateModelFromHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params, model interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := r.Body.Close(); err != nil {
		return err
	}
	if err := json.Unmarshal(body, model); err != nil {
		return err
	}

	return nil
}
