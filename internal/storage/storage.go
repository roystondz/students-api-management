package storage

import (
	"errors"

	"github.com/roystondz/students-api/internal/types"
)

var ErrStudentNotFound = errors.New("student not found")
var ErrUserNotFound = errors.New("user not found")
var ErrUserUnauthorised = errors.New("unauthorised access")

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudent(id int64) (types.Student, error)
	GetAll() ([]types.Student, error)
	DeleteStudent(id int64) error
	UpdateStudent(id int64, name string, email string, age int) (types.Student, error)
	SignUp(username string, password string) (int64, error)
	SignIn(username string, password string) (int64, error)
	Logout() (int64, error)
}