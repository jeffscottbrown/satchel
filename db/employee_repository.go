package db

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jeffscottbrown/satchel/model"
	"github.com/jeffscottbrown/satchel/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type employeeRepository struct {
	db *gorm.DB
}

// SaveEmployee implements repository.EmployeeRepository.
func (r *employeeRepository) SaveEmployee(employee *model.Employee) error {
	// Use the context.Background() to avoid passing nil context
	ctx := context.Background()
	if err := r.db.WithContext(ctx).Save(employee).Error; err != nil {
		slog.Error("failed to save employee", slog.Any("error", err))
		return err
	}
	slog.Info("employee saved successfully", slog.String("name", employee.Name))
	return nil
}

// GetEmployeeByName implements repository.EmployeeRepository.
func (r *employeeRepository) GetEmployeeByEmail(email string) (model.Employee, error) {
	var employee model.Employee
	err := r.db.WithContext(context.Background()).Preload("Scores").Where("email = ?", email).First(&employee).Error
	if err != nil {
		return model.Employee{}, err
	}
	return employee, nil
}

// GetEmployees implements repository.EmployeeRepository.
func (r *employeeRepository) GetEmployees() ([]model.Employee, error) {
	var employees []model.Employee
	err := r.db.WithContext(context.Background()).Find(&employees).Error
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func NewEmployeeRepository(db *gorm.DB) repository.EmployeeRepository {
	return &employeeRepository{db: db}
}

func InitializeDatabase() {
	dbUser := os.Getenv("SATCHEL_DB_USER")
	if dbUser == "" {
		slog.Error("environment variable SATCHEL_DB_USER not set")
		os.Exit(1)
	}
	dbPassword := os.Getenv("SATCHEL_DB_PASSWORD")
	if dbPassword == "" {
		slog.Error("environment variable SATCHEL_DB_PASSWORD not set")
		os.Exit(1)
	}
	dbName := os.Getenv("SATCHEL_DB_NAME")
	if dbName == "" {
		slog.Error("environment variable SATCHEL_DB_NAME not set")
		os.Exit(1)
	}
	dbHost := os.Getenv("SATCHEL_DB_HOST")
	if dbHost == "" {
		slog.Error("environment variable SATCHEL_DB_HOST not set")
		os.Exit(1)
	}
	dbPort := os.Getenv("SATCHEL_DB_PORT")
	if dbPort == "" {
		slog.Error("environment variable SATCHEL_DB_PORT not set")
		os.Exit(1)
	}
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbUser, dbPassword, dbName, dbHost, dbPort)

	var db *gorm.DB
	var err error
	for i := 0; i < 3; i++ {
		db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
		if err == nil {
			repository.SetRepository(NewEmployeeRepository(db))
			break
		}
		if i < 2 {
			slog.Error("failed to connect to database, retrying", slog.Int("attempt", i+1), slog.Any("error", err))

			time.Sleep(3 * time.Second)
		}
	}
	if err != nil {
		slog.Error("could not connect to database after 3 attempts", slog.Any("error", err))
		os.Exit(-1)
	}
	if err := db.AutoMigrate(&model.Employee{}, &model.Score{}); err != nil {
		slog.Error("failed to auto-migrate database", slog.Any("error", err))
		os.Exit(-1)
	}
	slog.Info("database initialized successfully")
}
