package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

// User represents the user data structure
type User struct {
	Email     string `json:"email"`
	Name      string `json:"name,omitempty"`
	Surname   string `json:"surname,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	Role      string `json:"role,omitempty"`
}

// UserDatabase represents the in-memory database of users
type UserDatabase struct {
	Users []User `json:"users"`
	mu    sync.Mutex
}

// Global variables
var (
	userDB     UserDatabase
	dbFilePath = "db.json"
)

// LoadDatabase loads the database from the JSON file
func LoadDatabase() error {
	file, err := os.ReadFile(dbFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// If the file doesn't exist, create an empty database
			userDB = UserDatabase{Users: []User{}}
			return SaveDatabase()
		}
		return err
	}

	return json.Unmarshal(file, &userDB)
}

// SaveDatabase saves the database to the JSON file
func SaveDatabase() error {
	data, err := json.MarshalIndent(userDB, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dbFilePath, data, 0644)
}

// FindUserByEmail finds a user by email
func FindUserByEmail(email string) (User, bool) {
	userDB.mu.Lock()
	defer userDB.mu.Unlock()

	for _, user := range userDB.Users {
		if user.Email == email {
			return user, true
		}
	}
	return User{}, false
}

// AuthenticateHandler handles authentication requests
func AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Parse user data
	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	// Check if email exists
	if existingUser, found := FindUserByEmail(user.Email); found {
		// User exists, return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"exists": true,
			"user":   existingUser,
		})
	} else {
		// User does not exist
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"exists": false,
		})
	}
}

// UpdateUserHandler handles user update requests (no longer updates passwords)
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Parse user data
	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	userDB.mu.Lock()
	defer userDB.mu.Unlock()

	// Find the user to check if they exist
	userFound := false
	for i, existingUser := range userDB.Users {
		if existingUser.Email == user.Email {
			// Update user data (excluding password)
			if user.Role != "" {
				userDB.Users[i].Role = user.Role
			}
			if user.Name != "" {
				userDB.Users[i].Name = user.Name
			}
			if user.Surname != "" {
				userDB.Users[i].Surname = user.Surname
			}
			userFound = true
			break
		}
	}

	if !userFound {
		// User does not exist, return error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "User not found. Only authorized users can be updated.",
		})
		return
	}

	// Save the updated database
	if err := SaveDatabase(); err != nil {
		http.Error(w, "Error saving database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
	})
}

func main() {
	// Load the database
	if err := LoadDatabase(); err != nil {
		log.Fatalf("Error loading database: %v", err)
	}

	// Set up HTTP routes
	http.HandleFunc("/authenticate", AuthenticateHandler)
	http.HandleFunc("/update", UpdateUserHandler)

	// Start the server
	port := 8080
	fmt.Printf("Server starting on port %d...\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
