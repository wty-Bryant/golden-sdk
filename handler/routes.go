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
		Route{"CreateProject", "POST", "/create-project", createProjectHandler},
		Route{"CreateProjectResources", "POST", "/create-project-resources", createProjectResourcesHandler},
		Route{"ListProjects", "GET", "/list-projects", listProjectsHandler},
		Route{"CreateWorkflow", "POST", "/create-workflow", createWorkflowHandler},
		Route{"ListWorkflows", "GET", "/list-workflows", listWorkflowsHandler},
		Route{"CreateWorkflowTrigger", "POST", "/create-workflow-trigger", createWorkflowTriggerHandler},
		Route{"ListTriggers", "GET", "/list-workflow-triggers", listWorkflowTriggersHandler},
	}
	return routes
}
