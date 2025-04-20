package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// --- Area Struct and Data ---
type Area struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	//Radius in Meter
	Radius   string `json:"radius"`
	Location string `json:"location"`
}

var areas []Area

// --- User Struct and Data ---
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Location string `json:"location"` // Added Location field
}

var users []User // Slice to store users

// --- Area Handlers ---

func getAreas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(areas)
}

func getArea(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Set header early
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid area ID", http.StatusBadRequest)
		return
	}

	// find the book with the given id
	for _, area := range areas {
		if area.ID == id {
			json.NewEncoder(w).Encode(area)
			return
		}
	}
	http.Error(w, "Area not found", http.StatusNotFound)
}

// Add a new area
func createArea(w http.ResponseWriter, r *http.Request) {
	var area Area
	// It's better practice to handle the error from Decode
	err := json.NewDecoder(r.Body).Decode(&area)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if area.Location == "" || area.Name == "" || area.Radius == "" {
		http.Error(w, "Please provide a name and location", http.StatusBadRequest)
		return
	}

	area.ID = len(areas) + 1 // Assign an ID (simple increment, not robust for production)
	areas = append(areas, area)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Use 201 Created for successful POST
	json.NewEncoder(w).Encode(area)
}

// --- User Handlers ---

// Get all users
func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Get a single user by ID
func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Set header early
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Find the user with the given id
	for _, user := range users {
		if user.ID == id {
			json.NewEncoder(w).Encode(user)
			return // Found and sent user
		}
	}

	// If loop completes, user was not found
	http.Error(w, "User not found", http.StatusNotFound)
}

// Create a new user
func createUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	// Decode the request body into the new user struct
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if newUser.Username == "" || newUser.Email == "" || newUser.Phone == "" || newUser.Location == "" {
		http.Error(w, "Username, Phone, Location and Email are required", http.StatusBadRequest)
		return
	}

	// Assign a simple ID (Not suitable for production/concurrent use)
	newUser.ID = len(users) + 1
	users = append(users, newUser) // Add the new user to the slice

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)  // Respond with 201 Created status
	json.NewEncoder(w).Encode(newUser) // Return the created user (with ID)
}

// --- Main Function ---

func main() {
	// Add some dummy data to start with
	areas = append(areas, Area{ID: 1, Radius: "85", Name: "Brisbane", Location: "-27.492887, 153.055914"})
	areas = append(areas, Area{ID: 2, Radius: "75", Name: "Sydney", Location: "-33.837386, 151.059379"})
	areas = append(areas, Area{ID: 3, Radius: "106", Name: "Melbourne", Location: "-37.822437, 145.011258"})

	// Add some dummy user data to start with (including Location)
	users = append(users, User{ID: 1, Username: "alice", Email: "alice@example.com", Phone: "0488079008", Location: "-37.8136, 144.9631"}) // Example Melbourne location
	users = append(users, User{ID: 2, Username: "bob", Email: "bob@example.com", Phone: "0488079009", Location: "-33.8688, 151.2093"})     // Example Sydney location
	users = append(users, User{ID: 3, Username: "peter", Email: "peter@example.com", Phone: "0488079010", Location: "-11.8688, 222.2093"}) // Example some location

	// Initialize the router
	r := mux.NewRouter()

	// Define the area endpoints
	r.HandleFunc("/areas", getAreas).Methods("GET")     // Get all areas
	r.HandleFunc("/areas/{id}", getArea).Methods("GET") // Get specific area
	r.HandleFunc("/areas", createArea).Methods("POST")  // Create a new area

	// Define the user endpoints
	r.HandleFunc("/users", getUsers).Methods("GET")     // Get all users
	r.HandleFunc("/users/{id}", getUser).Methods("GET") // Get specific user
	r.HandleFunc("/users", createUser).Methods("POST")  // Create a new user

	// Start the server
	fmt.Println("Server is running on port 8000...")
	log.Fatal(http.ListenAndServe("localhost:8000", r))
}
