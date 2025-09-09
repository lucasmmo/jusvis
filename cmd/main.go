package main

import (
	"jusvis/internal/auth"
	"jusvis/internal/citizen"
	"jusvis/internal/occurrence"
	"net/http"
)

func main() {
	occurrenceRepository := occurrence.NewMemoRepository()
	citizenRepository := citizen.NewMemoRepository()

	mux := http.NewServeMux()

	auth.Routes(mux, citizenRepository)
	occurrence.Routes(mux, occurrenceRepository, citizenRepository)

	http.ListenAndServe(":8080", mux)
}
