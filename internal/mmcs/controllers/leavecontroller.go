package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/anjanavikramrao/leave-application-backend/internal/mmcs/models"
	"github.com/anjanavikramrao/leave-application-backend/internal/mmcs/service"
	"github.com/gorilla/mux"
)

var baseUrl string = os.Getenv("BASE_URL")

type LeaveController interface {
	ControllerInterfaceGet
	ControllerInterfacePost
	ControllerInterfacePut
	ControllerInterfaceDelete
}

type leaveController struct {
	svc service.LeaveRequestService
}

func NewLeaveController() LeaveController {
	svc, err := service.NewLeaveRequestService()
	if err != nil {
		panic(err)
	}
	return &leaveController{svc: svc}
}

func (ctl *leaveController) Get(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	if leaveId, ok := pathParams["id"]; ok {
		req, err := ctl.svc.Get(leaveId)
		if err != nil {
			status := http.StatusInternalServerError
			if err.Error() == "record not found" {
				status = http.StatusNotFound
			}
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(makeErrorResponse(uint(status), err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(req)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(makeErrorResponse(http.StatusMethodNotAllowed, "Get all leave requests not supported."))
}

func (ctl *leaveController) Post(w http.ResponseWriter, r *http.Request) {
	var req *models.LeaveRequest = &models.LeaveRequest{}
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(makeErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	id, err := ctl.svc.Create(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(makeErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	w.Header().Set("Location", getBaseUrl()+id)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

func (ctl *leaveController) Put(w http.ResponseWriter, r *http.Request) {
	var req *models.LeaveRequest = &models.LeaveRequest{}
	w.Header().Set("Content-Type", "application/json")

	pathParams := mux.Vars(r)
	if leaveId, ok := pathParams["id"]; ok {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(makeErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		err = ctl.svc.Update(leaveId, req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(makeErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(makeErrorResponse(http.StatusInternalServerError, "Missing leave request identifier."))
}

func (ctl *leaveController) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pathParams := mux.Vars(r)
	if leaveId, ok := pathParams["id"]; ok {
		err := ctl.svc.Delete(leaveId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(makeErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(makeErrorResponse(http.StatusInternalServerError, "Missing leave request identifier."))
}

func makeErrorResponse(code uint, description string) *models.ErrorResponse {
	return &models.ErrorResponse{
		Code:        code,
		Description: description,
	}
}

func getBaseUrl() string {
	if !strings.HasSuffix(baseUrl, "/") {
		baseUrl = baseUrl + "/"
	}

	return fmt.Sprintf("%sleave/", baseUrl)
}
