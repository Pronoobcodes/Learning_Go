package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"notes-api/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var request authRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "request body is invalid")
		return
	}

	user, err := h.authService.Register(
		r.Context(),
		request.Email,
		request.Password,
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRegistration) {
			writeError(w, http.StatusBadRequest, "invalid_registration", err.Error())
			return
		}

		http.Error(w, "could not register user", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"message": "Registration successful",
		"user": map[string]any{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.authService.LoginWithToken(
		r.Context(),
		request.Email,
		request.Password,
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		http.Error(w, "Login failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Login successful",
		"token":   token,
	})
}
