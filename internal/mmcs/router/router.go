package router

import (
	"fmt"
	"net/http"
	"os"
	"reflect"

	"github.com/anjanavikramrao/leave-application-backend/internal/mmcs/controllers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const (
	GetAll = "/%s"
	Get    = "/%s/{id}"
	Post   = "/%s"
	Put    = "/%s/{id}"
	Delete = "/%s/{id}"
)

type Resources struct {
	LeaveRequest controllers.LeaveController `resource:"leave" description:"Leave Request resource"`
}

var allResources = Resources{
	LeaveRequest: controllers.NewLeaveController(),
}

func setupRoutes() http.Handler {
	router := mux.NewRouter()

	resources := reflect.ValueOf(allResources)
	types := resources.Type()

	for i := 0; i < types.NumField(); i++ {
		controller := resources.Field(i).Interface()
		if controller != nil {
			if resourceType, found := types.Field(i).Tag.Lookup("resource"); found {
				if c, ok := controller.(controllers.ControllerInterfaceGet); ok {
					router.HandleFunc(fmt.Sprintf(GetAll, resourceType), http.HandlerFunc(c.Get)).Methods("GET")
					router.HandleFunc(fmt.Sprintf(Get, resourceType), http.HandlerFunc(c.Get)).Methods("GET")
				}
				if c, ok := controller.(controllers.ControllerInterfacePost); ok {
					router.HandleFunc(fmt.Sprintf(Post, resourceType), http.HandlerFunc(c.Post)).Methods("POST")
				}
				if c, ok := controller.(controllers.ControllerInterfacePut); ok {
					router.HandleFunc(fmt.Sprintf(Put, resourceType), http.HandlerFunc(c.Put)).Methods("PUT")
				}
				if c, ok := controller.(controllers.ControllerInterfaceDelete); ok {
					router.HandleFunc(fmt.Sprintf(Delete, resourceType), http.HandlerFunc(c.Delete)).Methods("DELETE")
				}
			}
		}
	}

	handler := cors.Default().Handler(router)
	return handler
}

func initialize(handler http.Handler) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Service is listening on port %s ....\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), handler)
}

func Serve() {
	handler := setupRoutes()
	initialize(handler)
}
