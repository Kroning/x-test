package postgresql

import (
	"errors"
	"fmt"
	"github.com/Kroning/x-test/internal/config"
	"github.com/Kroning/x-test/internal/entities"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	// it's needed to initialize pgx driver
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	DB *sqlx.DB
}

var (
	ErrCompanyNotFound    = errors.New("company not found")
	ErrMoreThanOneCompany = errors.New("more than one company found")
)

var (
	maxIdleConns = 20
	maxOpenConns = 100
)

func InitAndMigrate(config config.Postgres) (*Database, error) {
	connectionStr := CreateConnectionString(
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Dbname,
	)

	db, err := Initialize(connectionStr)
	if err != nil {
		return nil, err
	}

	err = db.MigrateUp(config.MigrationsDirPath)
	if err != nil {
		return nil, err
	}

	if config.MaxIdleConns != 0 {
		maxIdleConns = config.MaxIdleConns
	}

	if config.MaxOpenConns != 0 {
		maxOpenConns = config.MaxOpenConns
	}

	db.DB.SetMaxIdleConns(maxIdleConns)
	db.DB.SetMaxOpenConns(maxOpenConns)

	return db, nil
}

func Initialize(connectionStr string) (*Database, error) {
	db, err := sqlx.Connect("pgx", connectionStr)
	if err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

func CreateConnectionString(host string, port string, user string, password string, dbname string) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}

func (dbObj *Database) MigrateUp(migrationsPath string) error {
	var err error
	driver, err := postgres.WithInstance(dbObj.DB.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	var m *migrate.Migrate

	m, err = migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err.Error() != "no change" {
		return err
	}
	return nil
}

func (dbObj *Database) CreateCompany(company *entities.Company) error {
	result, err := dbObj.DB.Exec(`INSERT INTO "company" (id, name, description, amount_of_employees, registered, type) VALUES ($1, $2, $3, $4, $5, $6)`,
		company.Id, company.Name, company.Description, company.AmountOfEmployees, company.Registered, company.Type,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return fmt.Errorf("expected to affect 1 row, affected %d", rows)
	}

	return nil
}

func (dbObj *Database) PatchCompany(company *entities.Company) error {
	result, err := dbObj.DB.Exec(`UPDATE "company" SET name = $2, description = $3, amount_of_employees = $4, registered = $5, type = $6 WHERE id = $1`,
		company.Id, company.Name, company.Description, company.AmountOfEmployees, company.Registered, company.Type,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return fmt.Errorf("expected to affect 1 row, affected %d", rows)
	}

	return nil
}

func (dbObj *Database) DeleteCompany(id string) error {
	result, err := dbObj.DB.Exec(`DELETE FROM "company" WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return fmt.Errorf("expected to affect 1 row, affected %d", rows)
	}

	return nil
}

func (dbObj *Database) GetCompany(id string) (*entities.Company, error) {
	companies := make([]entities.Company, 0)
	err := dbObj.DB.Select(&companies, `SELECT * FROM "company" where id = $1 ORDER BY id`, id)
	if err != nil {
		return nil, err
	}
	if len(companies) == 0 {
		return nil, ErrCompanyNotFound
	}
	if len(companies) > 1 {
		return nil, ErrMoreThanOneCompany
	}

	return &companies[0], nil
}
