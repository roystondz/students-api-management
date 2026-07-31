package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/roystondz/students-api/internal/storage"
	"github.com/roystondz/students-api/internal/types"
	"github.com/roystondz/students-api/internal/utils/response"
)

func Create(st storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)

		slog.Info("Creating a student")
		if errors.Is(err, io.EOF) {
			slog.Error("Please provide request body")
			response.WriteJson(w, http.StatusBadRequest, response.CommonError(errors.New("Please provide request body")))
			return
		}

		if err != nil {
			slog.Error("Invalid request body: " + err.Error())
			response.WriteJson(w, http.StatusBadRequest, response.CommonError(errors.New("Invalid request body: "+err.Error())))
			return
		}

		// Validating Request
		if err := validator.New().Struct(student); err != nil {
			validateErros := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErros))
			return
		}

		lastId, err := st.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			slog.Error("Failed to create student: " + err.Error())
			response.WriteJson(w, http.StatusInternalServerError, response.CommonError(errors.New("Failed to create student: "+err.Error())))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]string{"message": "student created", "id": strconv.FormatInt(lastId, 10)})
	}
}

func Get(st storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		slog.Info("Getting a student: ", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("Invalid student id: " + err.Error())
			response.WriteJson(w, http.StatusBadRequest, response.CommonError(errors.New("Invalid student id: "+err.Error())))
			return
		}

		student, e := st.GetStudent(intId)
		if e != nil {
			if errors.Is(e, storage.ErrStudentNotFound) {
				response.WriteJson(w, http.StatusNotFound, response.CommonError(e))
				return
			}
			slog.Error("Failed to get student: ", slog.Int64("id", intId))
			response.WriteJson(w, http.StatusInternalServerError, response.CommonError(errors.New("Failed to get student: "+e.Error())))
			return
		}

		response.WriteJson(w, http.StatusOK, student)

	}
}

func GetAll(st storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		students, err := st.GetAll()
		if err != nil {
			slog.Error("Failed to get students: " + err.Error())
			response.WriteJson(w, http.StatusInternalServerError, response.CommonError(errors.New("Failed to get students: "+err.Error())))
			return
		}
		response.WriteJson(w, http.StatusOK, students)
	}
}

func Delete(st storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		slog.Info("Deleting a student: ", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("Invalid student id: " + err.Error())
			response.WriteJson(w, http.StatusBadRequest, response.CommonError(errors.New("Invalid student id: "+err.Error())))
			return
		}

		err = st.DeleteStudent(intId)
		if err != nil {
			if errors.Is(err, storage.ErrStudentNotFound) {
				response.WriteJson(w, http.StatusNotFound, response.CommonError(err))
				return
			}
			slog.Error("Failed to delete student: " + err.Error())
			response.WriteJson(w, http.StatusInternalServerError, response.CommonError(errors.New("Failed to delete student: "+err.Error())))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"message": "student deleted"})
	}
}

func Update(st storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		slog.Info("Updating a student: ", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("Invalid student id: " + err.Error())
			response.WriteJson(w, http.StatusBadRequest, response.CommonError(errors.New("Invalid student id: "+err.Error())))
			return
		}

		var student types.Student
		err = json.NewDecoder(r.Body).Decode(&student)
		if err != nil {
			slog.Error("Failed to decode request body: " + err.Error())
			response.WriteJson(w, http.StatusBadRequest, response.CommonError(errors.New("Failed to decode request body: "+err.Error())))
			return
		}

		student, err = st.UpdateStudent(intId, student.Name, student.Email, student.Age)
		if err != nil {
			if errors.Is(err, storage.ErrStudentNotFound) {
				response.WriteJson(w, http.StatusNotFound, response.CommonError(err))
				return
			}
			slog.Error("Failed to update student: " + err.Error())
			response.WriteJson(w, http.StatusInternalServerError, response.CommonError(errors.New("Failed to update student: "+err.Error())))
			return
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}
