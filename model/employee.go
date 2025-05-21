package model

import (
	"sync"
)

type Employee struct {
	ID        uint           `gorm:"primaryKey"`
	Name      string         `yaml:"name"`
	Position  string         `yaml:"position"`
	Scores    map[string]any `yaml:"scores" gorm:"type:jsonb"`
	ImageName string         `yaml:"image"`
	mu        sync.Mutex     `gorm:"-"`
}

func (e *Employee) AddScore(scoreName string, value any) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.Scores == nil {
		e.Scores = make(map[string]any)
	}

	e.Scores[scoreName] = value
}
