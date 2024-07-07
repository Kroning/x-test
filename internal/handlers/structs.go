package handlers

import "github.com/Kroning/x-test/internal/entities"

type ResponseStatus struct {
	ErrorCode int    `json:"error_code,required"`
	Message   string `json:"message,required"`
}

type ResponseStatusWithID struct {
	ErrorCode int    `json:"error_code,required"`
	Message   string `json:"message,required"`
	Id        string `json:"id"`
}

type CompanyRequest struct {
	Name              string `json:"name" binding:"required"`
	Description       string `json:"description"`
	AmountOfEmployees int    `json:"amount_of_employees" binding:"required"`
	Registered        bool   `json:"registered" binding:"required"`
	Type              string `json:"type" binding:"required"`
}

type CompanyResponse struct {
	Status  ResponseStatus
	Company entities.Company
}
