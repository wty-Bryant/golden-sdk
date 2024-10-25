package handler

import "github.com/julienschmidt/httprouter"

/*
Define all the routes here.
A new Route entry passed to the routes slice will be automatically
translated to a handler with the NewRouter() function
*/
type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc httprouter.Handle
}

type Routes []Route

func AllRoutes() Routes {
	routes := Routes{
		Route{"Index", "GET", "/", index},
		Route{"CreateProject", "POST", "/resourceManager/createProject", createProjectHandler},
		Route{"CreateProjectResources", "POST", "/resourceManager/createProjectResources", createProjectResourcesHandler},
		Route{"ListProjects", "GET", "/resourceManager/listProjects", listProjectsHandler},
		Route{"CreateWorkflow", "POST", "/workflowManager/createWorkflow", createWorkflowHandler},
		Route{"ListWorkflows", "GET", "/workflowManager/listWorkflows", listWorkflowsHandler},
		Route{"CreateWorkflowTrigger", "POST", "/workflowTriggerManager/createTrigger", createWorkflowTriggerHandler},
		Route{"ListTriggers", "GET", "/workflowTriggerManager/listTriggers", listWorkflowTriggersHandler},
	}
	return routes
}
