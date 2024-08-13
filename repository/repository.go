package repository

import (
	"errors"
	"time"
)

var (
	updateError = errors.New("El update ha fallado")
	deleteError = errors.New("Delete error")
)

type Repository interface {
	Migrate() error
	InsertRegistre(registre Registres) (*Registres, error)
	LlegirTots Registres() ([]Registres, error)
	LlegirRegistrePerID(id int64) (*Registres, error)
	ActualitzarRegistre(id int64, actualitzar Registres) error
	BorrarRegistre(id int64) error

}

type Registres struct {
	ID int64 `json:"id"`
	Data time.Time `json:"data_registre"
	Precipitacio int `json:"precipitacio"`
	TempMaxima int `json:"temp_maxima"`
	TempMinima int `json:"temp_minima"`
	Humitat int `json:"humitat"`

}
