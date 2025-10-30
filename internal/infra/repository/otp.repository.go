package repository

import (
	"hospital_management_system/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines database operations for OTP
type OtpRepository interface {
	SaveOTP(otp *models.OTP) error
	GetOTPByCodeAndEmail(email string, code string) (*models.OTP, error)
	MarkOTPUsed(tx *gorm.DB, id uuid.UUID) error
	MarkUserVerified(tx *gorm.DB, email string) error
	Transaction(fn func(*gorm.DB) error) error
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
func (r *otpRepo) GetOTPByCodeAndEmail(email string, code string) (*models.OTP, error) {
	var otp models.OTP
	err := r.db.
		Where("email = ? AND code = ? AND is_used = false AND is_deleted = false", email, code).
		First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

// MarkOTPUsed updates OTP to mark it as used
func (r *otpRepo) MarkOTPUsed(tx *gorm.DB, id uuid.UUID) error {
	return tx.Model(&models.OTP{}).
		Where("id = ?", id).
		Update("is_used", true).Error
}

// MarkUserVerified updates user's is_verified to true
func (r *otpRepo) MarkUserVerified(tx *gorm.DB, email string) error {
	return tx.Model(&models.User{}).
		Where("email = ?", email).
		Update("is_verified", true).Error
}

// Transaction wraps database operations in a transaction
func (r *otpRepo) Transaction(fn func(*gorm.DB) error) error {
	return r.db.Transaction(fn)
}