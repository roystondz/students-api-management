package auth

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/roystondz/students-api/internal/storage"
	"github.com/roystondz/students-api/internal/types"
	"github.com/roystondz/students-api/internal/utils/response"
)


func SignUp(st storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){

		var user types.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.CommonError(err))
			return;
		}

		slog.Info("Signing up user", slog.String("username", user.Username))

		userId, err := st.SignUp(user.Username, user.Password)
		if err != nil {
			slog.Error("Failed to sign up user: " + err.Error())
			response.WriteJson(w, http.StatusInternalServerError, response.CommonError(err))
			return;
		}

		response.WriteJson(w, http.StatusCreated, map[string]string{"message": "user created", "id": strconv.FormatInt(userId, 10)})

	}
}

func SignIn(st storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		var user types.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.CommonError(err))
			return;
		}

		slog.Info("Signing in user", slog.String("username", user.Username))

		userId, err := st.SignIn(user.Username, user.Password)
		if err != nil {
			slog.Error("Failed to sign in user: " + err.Error())
			if errors.Is(err, storage.ErrUserNotFound){
				response.WriteJson(w, http.StatusUnauthorized, response.CommonError(err))
				return;
			}
			response.WriteJson(w, http.StatusInternalServerError, response.CommonError(err))
			return;
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"message": "user signed in", "id": strconv.FormatInt(userId, 10)})
	}
}
