package patient

import (
	"hospital_management_system/internal/pkg/helpers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Usecase interface {
	Create(req *PatientCreateRequest) (*Patient, error)
	CreateTx(txTx interface{}, req *PatientCreateRequest) (*Patient, error) // interface{} for DB transaction
}

type usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) Usecase {
	return &usecase{repo: repo}
}

type PatientCreateRequest struct {
	UserID         string `json:"user_id"`
	Age            int    `json:"age"`
	Gender         string `json:"gender"`
	Address        string `json:"address"`
	MedicalHistory string `json:"medical_history,omitempty"`
}

// Create patient
func (u *usecase) Create(req *PatientCreateRequest) (*Patient, error) {
	patient := &Patient{
		UserID:         uuidFromString(req.UserID),
		Age:            req.Age,
		Gender:         Gender(req.Gender),
		Address:        req.Address,
		MedicalHistory: req.MedicalHistory,
	}
	return u.repo.Create(patient)
}

// Create patient in transaction
func (u *usecase) CreateTx(tx interface{}, req *PatientCreateRequest) (*Patient, error) {
	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return nil, helpers.NewAppError(500, "Invalid transaction object")
	}

	patient := &Patient{
		UserID:         uuidFromString(req.UserID),
		Age:            req.Age,
		Gender:         Gender(req.Gender),
		Address:        req.Address,
		MedicalHistory: req.MedicalHistory,
	}

	return u.repo.CreateTx(gormTx, patient)
}

// helper to convert string to UUID
func uuidFromString(id string) uuid.UUID {
	uid, _ := uuid.Parse(id)
	return uid
}
