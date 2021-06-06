package models

type ErrorResponse struct {
	Code        uint   `json:"code"`
	Description string `json:"description"`
}
