package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
)

func signInHandler(res http.ResponseWriter, req *http.Request) {

	var payload SignInRequest

	// Only POST method is permitted
	if req.Method != http.MethodPost {
		writeJson(res, ErrorResponse{Error: fmt.Sprintf("Method not allowed")}, http.StatusMethodNotAllowed)
		return
	}

	// Deserializing the incoming json
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		writeJson(res, "JSON decoding error", http.StatusBadRequest)
		return
	}

	// Validation of a password
	if getEnvPassword() != "" && payload.Password == "" {
		writeJson(res, ErrorResponse{Error: fmt.Sprintf("Empty password")}, http.StatusUnauthorized)
		return
	}

	if payload.Password == getEnvPassword() {

		// Correct password, preparing JWT token
		token, err := generateToken(user.ID, user.Name)

		if err != nil {
			writeJson(res, ErrorResponse{Error: fmt.Sprintf("Internal security issue")}, http.StatusInternalServerError)
			return
		}

		// Sending successfull response
		writeJson(res, TokenResponse{Token: token}, http.StatusOK)

	} else {
		// Incorrect password has been provided
		writeJson(res, ErrorResponse{Error: fmt.Sprintf("Incorrect password")}, http.StatusUnauthorized)
		return

	}
}

// Generate toke JWT token
func generateToken(id int, name string) (string, error) {

	claims := jwt.MapClaims{
		"ID":   id,
		"Name": name,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := jwtToken.SignedString([]byte(generateSecretKey(name)))

	if err != nil {
		return "", nil
	}

	return signedToken, nil

}

// Generate a secret key
func generateSecretKey(source string) []byte {

	// Generating a hash from source
	hashedSource := (sha256.Sum256([]byte(source)))

	return hashedSource[:]

}
