package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/roystondz/students-api/internal/types"
	"github.com/roystondz/students-api/internal/utils/response"
)

func Create() http.HandlerFunc {
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

		response.WriteJson(w, http.StatusCreated, map[string]string{"message": "student created"})
	}
}
