package auth

import (
	"encoding/json"
	"fmt"
	"jusvis/internal/citizen"
	"jusvis/pkg/token"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Routes(mux *http.ServeMux, citizenRepository citizen.Repository) {
	controller := controller{citizenRepository: citizenRepository}
	mux.HandleFunc("POST /auth/login", controller.Login)
	mux.HandleFunc("POST /auth/register", controller.Register)
}

type controller struct {
	citizenRepository citizen.Repository
}

func (h controller) Login(w http.ResponseWriter, r *http.Request) {
	var body map[string]any

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "cannot decode json body", http.StatusBadRequest)
		return
	}

	email := body["email"].(string)
	password := body["password"].(string)

	if email == "" || password == "" {
		http.Error(w, "cannot decode json body", http.StatusBadRequest)
		return
	}

	cit, err := h.citizenRepository.GetByEmail(email)
	if err != nil {
		http.Error(w, "cannot find user", http.StatusBadRequest)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(cit.HashPassword), []byte(password)); err != nil {
		fmt.Println(cit.HashPassword, password)
		http.Error(w, "cannot validate password", http.StatusInternalServerError)
		return
	}

	tokenString, err := token.Generate(cit.ID)
	if err != nil {
		http.Error(w, "cannot generate token", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Bearer %s", tokenString)
}

func (h controller) Register(w http.ResponseWriter, r *http.Request) {
	var body map[string]any

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "cannot decode json body", http.StatusBadRequest)
		return
	}

	email := body["email"].(string)
	password := body["password"].(string)

	if email == "" || password == "" {
		http.Error(w, "cannot decode json body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "cannot hash password", http.StatusInternalServerError)
		return
	}
	cit := &citizen.Citizen{
		User: &citizen.User{
			ID:           uuid.NewString(),
			Email:        email,
			HashPassword: string(hashedPassword),
		},
		Address: nil,
	}

	if err := h.citizenRepository.Save(cit); err != nil {
		http.Error(w, "cannot save user", http.StatusInternalServerError)
		return
	}

}
