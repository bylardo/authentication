package models

type ErrorResponse struct {
	Error interface{} `json:"error"`
}

type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
