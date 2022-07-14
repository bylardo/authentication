package helpers

import (
	"encoding/json"
	"net/http"

	"authentication.com/models"
)

func PrepareAPIResponse(statusCode int, message string, data interface{}) models.APIResponse {

	var apiResponse models.APIResponse

	if statusCode != http.StatusAccepted && statusCode != http.StatusOK {
		var errorResponse models.ErrorResponse
		errorResponse.Error = message

		apiResponse.Status = statusCode
		apiResponse.Data = errorResponse
	}

	if statusCode == http.StatusAccepted || statusCode == http.StatusOK {
		apiResponse.Status = statusCode
		apiResponse.Message = message
		apiResponse.Data = data
	}
	return apiResponse
}

func SendAPIResponse(data models.APIResponse, w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
	return
}
