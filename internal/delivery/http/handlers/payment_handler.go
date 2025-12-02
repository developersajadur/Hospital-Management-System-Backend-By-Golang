package handlers

import (
	"hospital_management_system/internal/dto"
	"hospital_management_system/internal/pkg/helpers"
	"hospital_management_system/internal/pkg/utils"
	"hospital_management_system/internal/usecase"
	"net/http"
)

type PaymentHandler struct {
	uc usecase.PaymentUsecase
}

func NewPaymentHandler(uc usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{uc}
}

func (h *PaymentHandler) Init(w http.ResponseWriter, r *http.Request) {
	var req dto.InitPaymentRequest
	utils.BodyDecoder(w, r, &req)

	res, err := h.uc.InitPayment(&req)
	if err != nil {
		helpers.Error(w, err)
		return
	}

	helpers.Success(w, http.StatusOK, "Payment session created", res)
}

func (h *PaymentHandler) Success(w http.ResponseWriter, r *http.Request) {
	var cb dto.SSLCallbackRequest
	r.ParseForm()

	cb.TranID = r.FormValue("tran_id")
	cb.BankTranID = r.FormValue("bank_tran_id")
	cb.Amount = r.FormValue("amount")
	cb.CardType = r.FormValue("card_type")
	cb.ValID = r.FormValue("val_id")
	cb.PaymentDate = r.FormValue("tran_date")
	cb.Status = r.FormValue("status")

	if err := h.uc.HandleSuccessCallback(cb); err != nil {
		helpers.Error(w, err)
		return
	}

	helpers.Success(w, 200, "Payment successful", nil)
}

func (h *PaymentHandler) Fail(w http.ResponseWriter, r *http.Request) {
	var cb dto.SSLCallbackRequest
	r.ParseForm()
	cb.TranID = r.FormValue("tran_id")

	_ = h.uc.HandleFailCallback(cb)
	helpers.Success(w, 200, "Payment failed", nil)
}
