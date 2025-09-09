package occurrence

import (
	"encoding/json"
	"fmt"
	"jusvis/internal/citizen"
	"jusvis/internal/middleware"
	"net/http"
)

func Routes(mux *http.ServeMux, occurrenceRepository Repository, citizenRepository citizen.Repository) {
	controller := controller{
		occurrenceRepository: occurrenceRepository,
		citizenRepository:    citizenRepository,
	}
	mux.HandleFunc("POST /occurrence", middleware.Authorize(controller.Create))
	mux.HandleFunc("GET /occurrence/{id}", middleware.Authorize(controller.GetByID))
}

type controller struct {
	occurrenceRepository Repository
	citizenRepository    citizen.Repository
}

func (c controller) Create(w http.ResponseWriter, r *http.Request) {
	var body map[string]any

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("cannot decode json body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	ocType := body["type"].(string)
	userId := getUserIdFromToken(*r)

	createCommand := NewCreateCommand(c.occurrenceRepository, c.citizenRepository)
	if err := createCommand.Do(ocType, userId); err != nil {
		http.Error(w, fmt.Sprintf("cannot create an occurrence: %s", err.Error()), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c controller) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	oc, err := c.occurrenceRepository.GetByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot get an occurrence: %s", err.Error()), http.StatusBadRequest)
		return
	}

	cit, err := c.citizenRepository.GetByID(oc.RelatedBy)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot get an citizen: %s", err.Error()), http.StatusBadRequest)
		return
	}

	res := map[string]any{
		"id":   oc.ID,
		"type": string(oc.Type),
		"related_by": map[string]string{
			"id":    cit.ID,
			"email": cit.Email,
		},
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, fmt.Sprintf("cannot encode the occurrence: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func getUserIdFromToken(r http.Request) string {
	if userId := r.Header.Get("X-User-ID"); userId != "" {
		return userId
	}
	return ""
}
