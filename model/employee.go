package model

import (
	"sync"
)

type Employee struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Position    string
	Reflections []Reflection `gorm:"constraint:OnDelete:CASCADE;foreignKey:EmployeeID"`
	ImageName   string
	Email       string `gorm:"uniqueIndex;not null"`
	Bio         string
	mu          sync.Mutex `gorm:"-"`
}

type Reflection struct {
	ID         uint `gorm:"primaryKey"`
	Key        string
	Value      string
	EmployeeID uint
}

func (e *Employee) AddReflection(scoreName string, value string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.Reflections = append(e.Reflections, Reflection{
		Key:   scoreName,
		Value: value,
	})
}
