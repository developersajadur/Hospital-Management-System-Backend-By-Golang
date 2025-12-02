package dto

type InitPaymentRequest struct {
	BookingID string `json:"booking_id" validate:"required,uuid"`
}

type InitPaymentResponse struct {
	RedirectURL string `json:"redirect_url"`
	TranID      string `json:"tran_id"`
}

type SSLCallbackRequest struct {
	TranID      string `form:"tran_id"`
	BankTranID  string `form:"bank_tran_id"`
	Amount      string `form:"amount"`
	CardType    string `form:"card_type"`
	StoreAmount string `form:"store_amount"`
	ValID       string `form:"val_id"`
	PaymentDate string `form:"tran_date"`
	Status      string `form:"status"`
}
