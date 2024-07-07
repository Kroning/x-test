package entities

import (
	"fmt"
	"github.com/google/uuid"
)

// Should be separate from CompanyDTO, but had no time to implement it
type Company struct {
	Id                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	AmountOfEmployees int       `json:"amount_of_employees" db:"amount_of_employees"`
	Registered        bool      `json:"registered"`
	Type              string    `json:"type"`
	/*CreatedAt time.Time
	UpdatedAt time.Time*/
}

const (
	NameMaxLength        = 15
	DescriptionMaxLength = 300
)

func (c *Company) Validate() error {
	if c.Name == "" || len(c.Name) > NameMaxLength {
		return fmt.Errorf("name must not be empty and not more then %d characters", NameMaxLength)
	}
	if len(c.Description) > DescriptionMaxLength {
		return fmt.Errorf("description must not be empty and not more then %d characters", DescriptionMaxLength)
	}
	if c.AmountOfEmployees <= 0 {
		return fmt.Errorf("amount of employees must be more then 0")
	}
	return nil
}

func NewCompany(id uuid.UUID, name, description string, amountOfEmployees int, registered bool, companyType string) *Company {
	if id == uuid.Nil {
		id = uuid.New()
	}

	return &Company{
		Id:                id,
		Name:              name,
		Description:       description,
		AmountOfEmployees: amountOfEmployees,
		Registered:        registered,
		Type:              companyType,
	}
}
