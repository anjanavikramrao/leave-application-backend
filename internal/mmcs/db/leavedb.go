package db

import (
	"os"

	"github.com/anjanavikramrao/leave-application-backend/internal/mmcs/models"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var dbDialect string = "postgres"
var connectionString string = os.Getenv("CONNECTION_STRING")

type LeaveDatabase interface {
	Create(*models.LeaveRequest) (uint, error)
	Update(uint, *models.LeaveRequest) error
	Delete(uint) error
	Get(id uint) (*models.LeaveRequest, error)
	Close() error
}

type leaveRequestEntity struct {
	gorm.Model
	models.LeaveRequest
}

type leaveDatabase struct {
	db *gorm.DB
}

func NewLeaveDatabase() (LeaveDatabase, error) {
	db, err := gorm.Open(dbDialect, connectionString)
	if err != nil {
		return nil, err
	}

	result := db.AutoMigrate(&leaveRequestEntity{})
	if result.Error != nil {
		db.Close()
		return nil, result.Error
	}

	return &leaveDatabase{db: db}, nil
}

func (ldb *leaveDatabase) Create(req *models.LeaveRequest) (uint, error) {
	reqEntity := &leaveRequestEntity{LeaveRequest: *req}
	result := ldb.db.Create(reqEntity)
	if result.Error != nil {
		return 0, result.Error
	}

	return reqEntity.ID, nil
}

func (ldb *leaveDatabase) Update(id uint, req *models.LeaveRequest) error {
	reqEntity, err := ldb.findEntity(id)
	if err != nil {
		return err
	}

	reqEntity.LeaveRequest.Employee = req.Employee
	reqEntity.LeaveRequest.StartDate = req.StartDate
	reqEntity.LeaveRequest.EndDate = req.EndDate
	reqEntity.LeaveRequest.Reason = req.Reason
	reqEntity.LeaveRequest.Type = req.Type

	result := ldb.db.Save(reqEntity)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (ldb *leaveDatabase) Delete(id uint) error {
	reqEntity, err := ldb.findEntity(id)
	if err != nil {
		return err
	}

	result := ldb.db.Delete(reqEntity)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (ldb *leaveDatabase) Get(id uint) (*models.LeaveRequest, error) {
	reqEntity, err := ldb.findEntity(id)
	if err != nil {
		return nil, err
	}

	return &reqEntity.LeaveRequest, nil
}

func (ldb *leaveDatabase) Close() error {
	return ldb.db.Close()
}

func (ldb *leaveDatabase) findEntity(id uint) (*leaveRequestEntity, error) {
	reqEntity := &leaveRequestEntity{}
	result := ldb.db.First(reqEntity, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return reqEntity, nil
}
