package doctor

import (
	"hospital_management_system/internal/pkg/helpers"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Usecase interface {
	Create(req *DoctorCreateRequest) (*Doctor, error)
	CreateTx(tx *gorm.DB, req *DoctorCreateRequest) (*Doctor, error)
}

type usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) Usecase {
	return &usecase{repo: repo}
}

func (u *usecase) Create(req *DoctorCreateRequest) (*Doctor, error) {
	return u.CreateTx(nil, req) // call CreateTx without transaction
}

// CreateTx supports optional transaction
func (u *usecase) CreateTx(tx *gorm.DB, req *DoctorCreateRequest) (*Doctor, error) {
	// Convert UserID string to UUID
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, helpers.NewAppError(400, "Invalid user_id")
	}

	var status DoctorStatus
	switch req.Status {
	case "active":
		status = DoctorActive
	case "inactive":
		status = DoctorInactive
	case "on_leave":
		status = DoctorOnLeave
	default:
		return nil, helpers.NewAppError(http.StatusBadRequest, "Invalid status")
	}

	if req.Experience < 0 {
		return nil, helpers.NewAppError(400, "Experience cannot be negative")
	}
	if req.Fee < 0 {
		return nil, helpers.NewAppError(400, "Fee cannot be negative")
	}

	doctor := &Doctor{
		UserID:         userID,
		Specialization: req.Specialization,
		Experience:     req.Experience,
		Fee:            req.Fee,
		ProfileImage:   req.ProfileImage,
		Status:         status,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	var created *Doctor
	if tx != nil {
		created, err = u.repo.CreateTx(tx, doctor)
	} else {
		created, err = u.repo.Create(doctor)
	}

	if err != nil {
		return nil, helpers.NewAppError(500, err.Error())
	}
	return created, nil
}
