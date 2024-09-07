package repository

import (
	"errors"
	"time"
)

var (
	errUpdate = errors.New("el update ha fallado") //Por convención, las variables de error se nombran con el prefijo err
	errDelete = errors.New("delete error")
)

type Repository interface {
	Migrate() error
	InsertRegistre(registres Registres) (*Registres, error)
	LlegirTotsRegistres() ([]Registres, error)
	LlegirRegistrePerID(id int64) (*Registres, error)
	ActualitzarRegistre(id int64, actualitzar Registres) error
	BorrarRegistre(id int64) error
}

type Registres struct {
	ID           int64     `json:"id"`
	Data         time.Time `json:"data_registre"`
	Precipitacio int       `json:"precipitacio"`
	TempMaxima   int       `json:"temp_maxima"`
	TempMinima   int       `json:"temp_minima"`
	Humitat      int       `json:"humitat"`
}
