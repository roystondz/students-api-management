package sqlite

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"github.com/roystondz/students-api/internal/config"
	"github.com/roystondz/students-api/internal/storage"
	"github.com/roystondz/students-api/internal/types"
)

type Sqlite struct {
	db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		age INTEGER NOT NULL
	)
	`)
	if err != nil {
		return nil, err
	}

	return &Sqlite{db: db}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.db.Prepare(`
	INSERT INTO students (name, email, age) VALUES (?, ?, ?)
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}
	lastId, _ := result.LastInsertId()
	return lastId, nil
}

func (s *Sqlite) GetStudent(id int64) (types.Student, error) {
	stmt, err := s.db.Prepare(`SELECT * FROM students WHERE id = ? LIMIT 1`)
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.Student{}, storage.ErrStudentNotFound
		}
		return types.Student{}, err
	}
	return student, nil
}

func (s *Sqlite) GetAll() ([]types.Student, error) {
	stmt, err := s.db.Prepare(`SELECT * FROM students`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []types.Student
	for rows.Next() {
		var student types.Student
		err = rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	return students, nil
}

func (s *Sqlite) DeleteStudent(id int64) error {
	stmt, err := s.db.Prepare(`DELETE FROM students WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return storage.ErrStudentNotFound
	}
	return nil
}

func (s *Sqlite) UpdateStudent(id int64, name string, email string, age int) (types.Student, error) {
	stmt, err := s.db.Prepare(`UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?`)
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(name, email, age, id)
	if err != nil {
		return types.Student{}, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return types.Student{}, err
	}
	if rowsAffected == 0 {
		return types.Student{}, storage.ErrStudentNotFound
	}

	return types.Student{Id: id, Name: name, Email: email, Age: age}, nil
}