package handler

import "github.com/julienschmidt/httprouter"

func NewRouter(routes Routes) *httprouter.Router {

	router := httprouter.New()
	for _, route := range routes {
		var handle httprouter.Handle

		handle = route.HandlerFunc
		handle = logger(handle)

		router.Handle(route.Method, route.Path, handle)
	}

	return router
}
