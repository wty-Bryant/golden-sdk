package main

import (
	"github.com/golden-sdk/handler"
	"log"
	"net/http"
)

func main() {
	router := handler.NewRouter(handler.AllRoutes())
	log.Fatal(http.ListenAndServe(":8080", router))
}
