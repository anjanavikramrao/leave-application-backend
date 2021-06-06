package controllers

import (
	"net/http"
)

type ControllerInterfaceGet interface {
	Get(http.ResponseWriter, *http.Request)
}

type ControllerInterfacePost interface {
	Post(http.ResponseWriter, *http.Request)
}

type ControllerInterfacePut interface {
	Put(http.ResponseWriter, *http.Request)
}

type ControllerInterfaceDelete interface {
	Delete(http.ResponseWriter, *http.Request)
}
