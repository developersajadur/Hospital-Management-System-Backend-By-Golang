package user

import (
	"hospital_management_system/internal/pkg/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase Usecase
}

func NewHandler(uc Usecase) *Handler {
	return &Handler{usecase: uc}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Error(c, helpers.NewAppError(400, err.Error()))
		return
	}

	user, err := h.usecase.Register(&req)
	if err != nil {
		helpers.Error(c, err) 
		return
	}

	helpers.Success(c, 201, "User registered successfully", user)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Error(c, helpers.NewAppError(http.StatusBadRequest, err.Error()))
		return
	}

	token, err := h.usecase.Login(&req)
	if err != nil {
			helpers.Error(c, helpers.NewAppError(http.StatusBadRequest, err.Error()))
		return
	}

	helpers.Success(c, http.StatusOK, "Login successful", map[string]string{
		"token": token,
	})
}
