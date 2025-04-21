package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	Radius    string `json:"radius"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

var areas []Area

// --- User Struct and Data ---
type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
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

	if area.Latitude == "" || area.Longitude == "" || area.Name == "" || area.Radius == "" {
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
	if newUser.Username == "" || newUser.Email == "" || newUser.Phone == "" || newUser.Latitude == "" || newUser.Longitude == "" {
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

// --- User Handlers ---

// ... (keep existing getUsers, getUser, createUser functions) ...

// Update a user's latitude and longitude
func updateUserLocation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Define a temporary struct to decode latitude and longitude from the request body
	var locationUpdate struct {
		// Use pointers to string to differentiate between empty string and field not provided,
		// although for PATCH requiring both might be simpler/intended here.
		// If you *always* expect both, you can use string directly.
		// Let's assume for this PATCH, both are required.
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
	}

	// Decode the request body into the temporary struct
	err = json.NewDecoder(r.Body).Decode(&locationUpdate)
	if err != nil {
		// Check for EOF which means empty body, could be handled differently if needed
		if err == io.EOF {
			http.Error(w, "Request body cannot be empty", http.StatusBadRequest)
		} else {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		}
		return
	}

	// Validate that both latitude and longitude were provided
	// Note: This checks for empty strings. If "0" is a valid coordinate, this is fine.
	if locationUpdate.Latitude == "" || locationUpdate.Longitude == "" {
		http.Error(w, "Both latitude and longitude fields are required in the request body", http.StatusBadRequest)
		return
	}

	// --- Optional: Add validation for coordinate format if needed ---
	// Example (very basic, might need a regex or library for robust validation):
	// _, errLat := strconv.ParseFloat(locationUpdate.Latitude, 64)
	// _, errLon := strconv.ParseFloat(locationUpdate.Longitude, 64)
	// if errLat != nil || errLon != nil {
	//     http.Error(w, "Invalid format for latitude or longitude", http.StatusBadRequest)
	//     return
	// }
	// --- End Optional Validation ---

	// Find the user and update their location fields
	found := false
	var updatedUser User   // To store the user data to return
	for i := range users { // Iterate by index to modify the original slice element
		if users[i].ID == id {
			users[i].Latitude = locationUpdate.Latitude   // Update Latitude
			users[i].Longitude = locationUpdate.Longitude // Update Longitude
			updatedUser = users[i]                        // Get the updated user data
			found = true
			break // Exit loop once found and updated
		}
	}

	// Handle case where user was not found
	if !found {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Respond with the updated user data
	w.WriteHeader(http.StatusOK) // 200 OK
	json.NewEncoder(w).Encode(updatedUser)
}

// --- Main Function ---

func main() {
	// Add some dummy data to start with
	areas = append(areas, Area{ID: 1, Radius: "42", Name: "Brisbane", Latitude: "-27.492887", Longitude: "153.055914"})
	areas = append(areas, Area{ID: 2, Radius: "37", Name: "Sydney", Latitude: "-33.837386", Longitude: "151.059379"})
	areas = append(areas, Area{ID: 3, Radius: "53", Name: "Melbourne", Latitude: "-37.822437", Longitude: "145.011258"})
	areas = append(areas, Area{ID: 4, Radius: "123", Name: "Test", Latitude: "-37.822437", Longitude: "145.011258"})

	// Add some dummy user data to start with (including Location)
	users = append(users, User{ID: 1, Username: "alice", Email: "alice@example.com", Phone: "0488079008", Latitude: "-27.492887", Longitude: "153.055914"})   // Example Melbourne location
	users = append(users, User{ID: 2, Username: "bob", Email: "bob@example.com", Phone: "0488079009", Latitude: "-33.837386", Longitude: "151.059379"})       // Example Sydney location
	users = append(users, User{ID: 3, Username: "peter", Email: "peter@example.com", Phone: "0488079010", Latitude: "-37.822437", Longitude: "145.011258"})   // Example some location
	users = append(users, User{ID: 4, Username: "paul", Email: "paul@example.com", Phone: "0488079011", Latitude: "-27.492887", Longitude: "153.055914"})     // Example some location
	users = append(users, User{ID: 5, Username: "daniel", Email: "daniel@example.com", Phone: "0488079012", Latitude: "-37.822437", Longitude: "145.011258"}) // Example some location

	// Initialize the router
	r := mux.NewRouter()

	// Define the area endpoints
	r.HandleFunc("/areas", getAreas).Methods("GET")     // Get all areas
	r.HandleFunc("/areas/{id}", getArea).Methods("GET") // Get specific area
	r.HandleFunc("/areas", createArea).Methods("POST")  // Create a new area

	// Define the user endpoints
	r.HandleFunc("/users", getUsers).Methods("GET")                  // Get all users
	r.HandleFunc("/users/{id}", getUser).Methods("GET")              // Get specific user
	r.HandleFunc("/users", createUser).Methods("POST")               // Create a new user
	r.HandleFunc("/users/{id}", updateUserLocation).Methods("PATCH") // Update user location

	// Start the server
	fmt.Println("Server is running on port 8000...")
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", r))
}
