package model

import (
	"sync"
)

type Employee struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Position  string
	Scores    []Score `gorm:"foreignKey:EmployeeID"`
	ImageName string
	Email     string
	mu        sync.Mutex `gorm:"-"`
}

type Score struct {
	ID         uint `gorm:"primaryKey"`
	Key        string
	Value      string
	EmployeeID uint
}

func (e *Employee) AddScore(scoreName string, value string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.Scores = append(e.Scores, Score{
		Key:   scoreName,
		Value: value,
	})
}
