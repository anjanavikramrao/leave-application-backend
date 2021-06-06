package models

import (
	"time"
)

const (
	LeaveStatusReady    = iota
	LeaveStatusDenied   = iota
	LeaveStatusApproved = iota
)

type LeaveRequest struct {
	Employee  string    `json:"employee"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Reason    string    `json:"reason"`
	Type      int       `json:"type"`
}
