package auth

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/roystondz/students-api/internal/storage"
	"github.com/roystondz/students-api/internal/types"
	"github.com/roystondz/students-api/internal/utils/response"

	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	Storage storage.Storage
	JWTKey  []byte
}

func New(st storage.Storage, jwtKey string) *AuthHandler {
	return &AuthHandler{
		Storage: st,
		JWTKey:  []byte(jwtKey),
	}
}

func (a *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user types.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response.WriteJson(w, http.StatusBadRequest, response.CommonError(err))
		return
	}

	slog.Info("Signing up user", slog.String("username", user.Username))

	userId, err := a.Storage.SignUp(user.Username, user.Password)
	if err != nil {
		slog.Error("Failed to sign up user: " + err.Error())
		response.WriteJson(w, http.StatusInternalServerError, response.CommonError(err))
		return
	}

	// TODO: Check if we can use the same token generation logic here as well
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userId,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(a.JWTKey)
	if err != nil {
		slog.Error("Failed to sign token: " + err.Error())
		response.WriteJson(w, http.StatusInternalServerError, response.CommonError(err))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true, // Prevents frontend JS from reading the cookie (protects against XSS)
		Secure:   true, // Ensures cookie is only sent over HTTPS
		Path:     "/",
		SameSite: http.SameSiteLaxMode, // Helps protect against CSRF
	})

	response.WriteJson(w, http.StatusCreated, map[string]string{"message": "user created", "id": strconv.FormatInt(userId, 10), "token": tokenString})

}

func (a *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var user types.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response.WriteJson(w, http.StatusBadRequest, response.CommonError(err))
		return
	}

	slog.Info("Signing in user", slog.String("username", user.Username))

	userId, err := a.Storage.SignIn(user.Username, user.Password)
	if err != nil {
		slog.Error("Failed to sign in user: " + err.Error())
		if errors.Is(err, storage.ErrUserNotFound) {
			response.WriteJson(w, http.StatusUnauthorized, response.CommonError(err))
			return
		}
		response.WriteJson(w, http.StatusInternalServerError, response.CommonError(err))
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userId,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(a.JWTKey)
	if err != nil {
		slog.Error("Failed to sign token: " + err.Error())
		response.WriteJson(w, http.StatusInternalServerError, response.CommonError(err))
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true, // Prevents frontend JS from reading the cookie (protects against XSS)
		Secure:   true, // Ensures cookie is only sent over HTTPS
		Path:     "/",
		SameSite: http.SameSiteLaxMode, // Helps protect against CSRF
	})

	response.WriteJson(w, http.StatusOK, map[string]string{"message": "user signed in", "id": strconv.FormatInt(userId, 10), "token": tokenString})
}

// TODO: Find out ideas on how this can be implemented properly
func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true, // Prevents frontend JS from reading the cookie (protects against XSS)
		Secure:   true, // Ensures cookie is only sent over HTTPS
		Path:     "/",
		SameSite: http.SameSiteLaxMode, // Helps protect against CSRF
	})

	response.WriteJson(w, http.StatusOK, map[string]string{"message": "user logged out"})
}
