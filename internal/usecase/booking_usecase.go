package usecase

import (
	"errors"
	"hospital_management_system/internal/dto"
	"hospital_management_system/internal/infra/repository"
	"hospital_management_system/internal/models"
	"hospital_management_system/internal/pkg/helpers"
	"hospital_management_system/internal/pkg/utils"
	"net/http"

	"gorm.io/gorm"
)

type BookingUsecase interface {
	Create(req *dto.CreateBookingRequest) (*models.Booking, error)
	GetByID(id string) (*models.Booking, error)
	GetAll() ([]models.Booking, error)
	UpdateStatus(id string, req *dto.UpdateBookingStatusRequest) (*models.Booking, error)
	Delete(id string) error
}

type bookingUsecase struct {
	bookingRepo repository.BookingRepository
	patientRepo repository.PatientRepository
	roomRepo    repository.RoomRepository
	serviceRepo repository.ServiceRepository
}

func BookingNewUsecase(
	bookingRepo repository.BookingRepository,
	patientRepo repository.PatientRepository,
	roomRepo repository.RoomRepository,
	serviceRepo repository.ServiceRepository,
) BookingUsecase {
	return &bookingUsecase{
		bookingRepo: bookingRepo,
		patientRepo: patientRepo,
		roomRepo:    roomRepo,
		serviceRepo: serviceRepo,
	}
}


func (u *bookingUsecase) Create(req *dto.CreateBookingRequest) (*models.Booking, error) {

	_, err := u.patientRepo.GetPatientByID(req.PatientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewAppError(http.StatusNotFound, "Patient not found")
		}
		return nil, helpers.NewAppError(http.StatusInternalServerError, "Database error")
	}

	booking := &models.Booking{
		BookingType: models.BookingType(req.BookingType),
		PatientID:   models.UUIDFromString(req.PatientID),
		Status:      models.BookingPending,
	}

	if req.BookingType == "room" {

		if req.RoomID == nil {
			return nil, helpers.NewAppError(http.StatusBadRequest, "room_id is required")
		}

		roomID := models.UUIDFromString(*req.RoomID)

		room, err := u.roomRepo.GetRoomByID(roomID.String())
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, helpers.NewAppError(http.StatusNotFound, "Room not found")
			}
			return nil, helpers.NewAppError(http.StatusInternalServerError, "Database error")
		}

		if !room.Availability {
			return nil, helpers.NewAppError(http.StatusBadRequest, "Room is not available")
		}

		booking.RoomID = utils.UUIDPtr(req.RoomID)
		booking.CheckInDate = req.CheckInDate
		booking.CheckOutDate = req.CheckOutDate
	}

	if req.BookingType == "service" {

		if req.ServiceID == nil {
			return nil, helpers.NewAppError(http.StatusBadRequest, "service_id is required")
		}

		serviceID := models.UUIDFromString(*req.ServiceID)

		_, err := u.serviceRepo.GetServiceByID(serviceID.String())
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, helpers.NewAppError(http.StatusNotFound, "Service not found")
			}
			return nil, helpers.NewAppError(http.StatusInternalServerError, "Database error")
		}

		day := req.ScheduledAt.Format("2006-01-02")

		count, _ := u.bookingRepo.CountServiceBookingsForDay(*req.ServiceID, day)
		serial := int(count) + 1

		booking.ServiceID = utils.UUIDPtr(req.ServiceID)
		booking.ScheduledAt = req.ScheduledAt
		booking.SerialNumber = &serial
	}

	return u.bookingRepo.Create(booking)
}


func (u *bookingUsecase) GetByID(id string) (*models.Booking, error) {
	b, err := u.bookingRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewAppError(http.StatusNotFound, "Booking not found")
		}
		return nil, helpers.NewAppError(http.StatusInternalServerError, "Database error")
	}
	return b, nil
}

func (u *bookingUsecase) GetAll() ([]models.Booking, error) {
	return u.bookingRepo.GetAll()
}

func (u *bookingUsecase) UpdateStatus(id string, req *dto.UpdateBookingStatusRequest) (*models.Booking, error) {
	b := &models.Booking{ID: models.UUIDFromString(id)}
	b.Status = models.BookingStatus(req.Status)

	return u.bookingRepo.Update(b)
}

func (u *bookingUsecase) Delete(id string) error {
	return u.bookingRepo.Delete(id)
}
