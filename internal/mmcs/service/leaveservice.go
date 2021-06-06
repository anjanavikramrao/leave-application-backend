package service

import (
	"strconv"

	"github.com/anjanavikramrao/leave-application-backend/internal/mmcs/db"
	"github.com/anjanavikramrao/leave-application-backend/internal/mmcs/models"
)

const (
	LeaveTypePrivileged = iota
	LeaveTypeSick       = iota
)

type LeaveRequestService interface {
	Create(*models.LeaveRequest) (string, error)
	Update(string, *models.LeaveRequest) error
	Delete(string) error
	Get(id string) (*models.LeaveRequest, error)
}

type leaveRequestService struct {
	leaveDb db.LeaveDatabase
}

func NewLeaveRequestService() (LeaveRequestService, error) {
	db, err := db.NewLeaveDatabase()
	if err != nil {
		return nil, err
	}

	return &leaveRequestService{leaveDb: db}, nil
}

func (svc *leaveRequestService) Create(req *models.LeaveRequest) (string, error) {
	id, err := svc.leaveDb.Create(req)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(uint64(id), 10), nil
}

func (svc *leaveRequestService) Update(id string, req *models.LeaveRequest) error {
	identifier, err := getId(id)
	if err != nil {
		return err
	}
	return svc.leaveDb.Update(identifier, req)
}

func (svc *leaveRequestService) Delete(id string) error {
	identifier, err := getId(id)
	if err != nil {
		return err
	}
	return svc.leaveDb.Delete(identifier)
}

func (svc *leaveRequestService) Get(id string) (*models.LeaveRequest, error) {
	identifier, err := getId(id)
	if err != nil {
		return nil, err
	}
	return svc.leaveDb.Get(identifier)
}

func getId(identifier string) (uint, error) {
	id, err := strconv.ParseUint(identifier, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}
