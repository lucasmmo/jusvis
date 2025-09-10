package main

import (
	"jusvis/internal/auth"
	"jusvis/internal/citizen"
	"jusvis/internal/middleware"
	"jusvis/internal/occurrence"
	"net/http"
)

func main() {
	occurrenceRepository := occurrence.NewMemoRepository()
	citizenRepository := citizen.NewMemoRepository()

	mux := http.NewServeMux()

	auth.Routes(mux, citizenRepository)
	occurrence.Routes(mux, occurrenceRepository, citizenRepository)

	corsMux := middleware.Cors(mux)

	http.ListenAndServe(":8080", corsMux)
}
