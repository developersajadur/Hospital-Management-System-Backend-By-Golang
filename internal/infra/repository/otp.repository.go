package repository

import (
	"hospital_management_system/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines database operations for OTP
type OtpRepository interface {
	SaveOTP(otp *models.OTP) error
	GetOTPByCode(userID uuid.UUID, code string, purpose string) (*models.OTP, error)
	MarkOTPUsed(id uuid.UUID) error
}

// repository implementation
type otpRepo struct {
	db *gorm.DB
}

// NewRepository creates a new repository instance
func OtpNewRepository(db *gorm.DB) OtpRepository {
	return &otpRepo{db: db}
}

// SaveOTP inserts new OTP into DB
func (r *otpRepo) SaveOTP(otp *models.OTP) error {
	return r.db.Create(otp).Error
}

// GetOTPByCode fetches OTP by userID, code, and purpose
func (r *otpRepo) GetOTPByCode(userID uuid.UUID, code string, purpose string) (*models.OTP, error) {
	var otp models.OTP
	err := r.db.
		Where("user_id = ? AND code = ? AND purpose = ? AND is_used = false AND is_deleted = false", userID, code, purpose).
		First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

// MarkOTPUsed updates OTP to mark it as used
func (r *otpRepo) MarkOTPUsed(id uuid.UUID) error {
	return r.db.Model(&models.OTP{}).
		Where("id = ?", id).
		Update("is_used", true).Error
}
