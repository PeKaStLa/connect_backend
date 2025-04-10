package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Area struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

var areas []Area

func getAreas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(areas)
}

func getArea(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid area ID", http.StatusBadRequest)
		return
	}

	if id == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(areas)
		return
	}

	// find the book with the given id
	for _, area := range areas {
		if area.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(area)
			return
		}
	}
	http.Error(w, "Area not found", http.StatusNotFound)
}

// Add a new area
func createArea(w http.ResponseWriter, r *http.Request) {
	var area Area
	_ = json.NewDecoder(r.Body).Decode(&area)

	if area.Location == "" || area.Name == "" {
		http.Error(w, "Please provide a name and location", http.StatusBadRequest)
		return
	}

	area.ID = len(areas) + 1 // Assign an ID (we’re just winging it here)
	areas = append(areas, area)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(area)
}

func main() {
	// Add some dummy data to start with
	areas = append(areas, Area{ID: 1, Name: "Brisbane", Location: "-27.492887, 153.055914"})
	areas = append(areas, Area{ID: 2, Name: "Sydney", Location: "-33.837386, 151.059379"})
	areas = append(areas, Area{ID: 3, Name: "Melbourne", Location: "-37.822437, 145.011258"})

	// Initialize the router
	r := mux.NewRouter()

	// Define the endpoints
	r.HandleFunc("/areas", getAreas).Methods("GET")
	r.HandleFunc("/areas/{id}", getArea).Methods("GET")
	r.HandleFunc("/areas", createArea).Methods("POST")

	// Start the server
	fmt.Println("Server is running on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", r))
}
