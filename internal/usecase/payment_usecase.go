package usecase

import (
	"encoding/json"
	"fmt"
	"hospital_management_system/config"
	"hospital_management_system/internal/dto"
	"hospital_management_system/internal/infra/repository"
	"hospital_management_system/internal/models"
	"hospital_management_system/internal/pkg/helpers"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

type PaymentUsecase interface {
	InitPayment(req *dto.InitPaymentRequest) (*dto.InitPaymentResponse, error)
	HandleSuccessCallback(req dto.SSLCallbackRequest) error
	HandleFailCallback(req dto.SSLCallbackRequest) error
	GetAll() ([]models.Payment, error)
}

type paymentUsecase struct {
	paymentRepo repository.PaymentRepository
	bookingRepo repository.BookingRepository
}

func NewPaymentUsecase(paymentRepo repository.PaymentRepository, bookingRepo repository.BookingRepository) PaymentUsecase {
	return &paymentUsecase{paymentRepo, bookingRepo}
}

func (u *paymentUsecase) InitPayment(req *dto.InitPaymentRequest) (*dto.InitPaymentResponse, error) {
	booking, err := u.bookingRepo.GetByID(req.BookingID)
	if err != nil {
		return nil, helpers.NewAppError(404, "Booking not found")
	}
	if booking.TotalPrice == nil {
		return nil, helpers.NewAppError(400, "Total price missing")
	}

	tranID := uuid.New().String()

	payment := &models.Payment{
		ID:        uuid.New(),
		BookingID: booking.ID,
		TranID:    tranID,
		Amount:    *booking.TotalPrice,
		Status:    models.PaymentInitiated,
	}

	if err := u.paymentRepo.Create(payment); err != nil {
		return nil, helpers.NewAppError(500, "Failed to create payment")
	}

	payload := url.Values{}
	payload.Set("store_id", config.ENV.SSLStoreID)
	payload.Set("store_passwd", config.ENV.SSLStorePassword)
	payload.Set("total_amount", fmt.Sprintf("%.2f", *booking.TotalPrice))
	payload.Set("currency", "BDT")
	payload.Set("tran_id", tranID)
	payload.Set("success_url", config.ENV.BaseURL+"/payments/success")
	payload.Set("fail_url", config.ENV.BaseURL+"/payments/fail")
	payload.Set("cancel_url", config.ENV.BaseURL+"/payments/cancel")
	payload.Set("cus_name", "Customer")
	payload.Set("cus_email", "customer@test.com")
	payload.Set("cus_phone", "01700000000")
	payload.Set("cus_add1", "Customer Address")
	payload.Set("cus_city", "Dhaka")
	payload.Set("cus_country", "Bangladesh")
	payload.Set("shipping_method", "NO")
	payload.Set("num_of_item", "1")
	payload.Set("product_name", "Hospital Service")
	payload.Set("product_category", "Healthcare")
	payload.Set("product_profile", "general")

	resp, err := http.PostForm("https://sandbox.sslcommerz.com/gwprocess/v4/api.php", payload)
	if err != nil {
		return nil, helpers.NewAppError(500, "SSLCommerz request failed")
	}
	defer resp.Body.Close()

	var sslResp struct {
		Status         string `json:"status"`
		GatewayPageURL string `json:"GatewayPageURL"`
		FailedReason   string `json:"failedreason"`
	}
	json.NewDecoder(resp.Body).Decode(&sslResp)

	if sslResp.Status != "SUCCESS" || sslResp.GatewayPageURL == "" {
		return nil, helpers.NewAppError(500, sslResp.FailedReason)
	}

	return &dto.InitPaymentResponse{
		RedirectURL: sslResp.GatewayPageURL,
		TranID:      tranID,
	}, nil
}

func (u *paymentUsecase) HandleSuccessCallback(req dto.SSLCallbackRequest) error {
	payment, err := u.paymentRepo.GetByTranID(req.TranID)
	if err != nil {
		return helpers.NewAppError(404, "Payment record not found")
	}

	t, _ := time.Parse("2006-01-02 15:04:05", req.PaymentDate)
	payment.Status = models.PaymentSuccess
	payment.Method = req.CardType
	payment.BankTranID = req.BankTranID
	payment.ValidationID = req.ValID
	payment.TransactionAt = &t

	if err := u.paymentRepo.Update(payment); err != nil {
		return helpers.NewAppError(500, "Failed to update payment")
	}

	return u.bookingRepo.UpdateStatus(payment.BookingID.String(), models.BookingConfirmed)
}

func (u *paymentUsecase) HandleFailCallback(req dto.SSLCallbackRequest) error {
	payment, err := u.paymentRepo.GetByTranID(req.TranID)
	if err != nil {
		return nil
	}
	payment.Status = models.PaymentFailed
	return u.paymentRepo.Update(payment)
}


func (u *paymentUsecase) GetAll() ([]models.Payment, error) {
	return u.paymentRepo.GetAll()
}
