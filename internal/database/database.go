package database

import "github.com/Kroning/x-test/internal/entities"

type AbstractDB interface {
	CreateCompany(company *entities.Company) error
	PatchCompany(company *entities.Company) error
	DeleteCompany(id string) error
	GetCompany(id string) (*entities.Company, error)
}
