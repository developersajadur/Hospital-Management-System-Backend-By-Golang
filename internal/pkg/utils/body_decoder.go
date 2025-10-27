package utils

import (
	"encoding/json"
	"hospital_management_system/internal/pkg/helpers"
	"net/http"
)

func BodyDecoder(w http.ResponseWriter, r *http.Request, data interface{})  {
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "Can't decode Data from body"))
		return 

	}
}
