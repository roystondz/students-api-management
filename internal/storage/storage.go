package storage

import (
	"errors"

	"github.com/roystondz/students-api/internal/types"
)

var ErrStudentNotFound = errors.New("student not found")

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudent(id int64) (types.Student, error)
	GetAll() ([]types.Student, error)
	DeleteStudent(id int64) error
	UpdateStudent(id int64, name string, email string, age int) (types.Student, error)
}