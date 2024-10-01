// user.go
package models

import "log"

// User struct represents a user in the system
type User struct {
	Username string
	Password string
}

// RegisterUser registers a new user in the database
func RegisterUser(username, password string) (*User, error) {
	stmt, err := DB.Prepare("INSERT INTO users(username, password) VALUES(?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, password)
	if err != nil {
		return nil, err
	}

	return &User{Username: username, Password: password}, nil
}

// AuthenticateUser authenticates a user by checking their credentials
func AuthenticateUser(username, password string) bool {
	var storedPassword string
	err := DB.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPassword)
	if err != nil {
		log.Println("Error fetching user data:", err)
		return false
	}

	return storedPassword == password
}

// GetUserProfile retrieves a user's profile from the database
func GetUserProfile(username string) *User {
	var user User
	err := DB.QueryRow("SELECT username, password FROM users WHERE username = ?", username).Scan(&user.Username, &user.Password)
	if err != nil {
		log.Println("Error fetching user profile:", err)
		return nil
	}

	return &user
}
