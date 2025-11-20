package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"tracker/pkg/db"

	"github.com/golang-jwt/jwt/v5"
)

const apiDateFormat = "20060102"

type NewTaskResponse struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

type TasksResponse struct {
	Tasks []db.Task `json:"tasks"`
}

type EmptyResponse struct{}

type SignInRequest struct {
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

// JWT Claim for user in token
type Claims struct {
	jwt.RegisteredClaims
	ID   int
	Name string
}

// User

type User struct {
	ID   int
	Name string
}

var user User

// Initialization of server handlers and user
func Init() {

	// Setting handlers

	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(taskDoneHandler))
	http.HandleFunc("/api/signin", signInHandler)

	// Setting user

	user = User{
		ID: 1,
		Name: "Andrew",
	}

}

// HTTP response JSON writer
func writeJson(w http.ResponseWriter, data interface{}, httpStatus int) {

	// Setting content type
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Setting http response code
	w.WriteHeader(httpStatus)

	// Creating JSON object and sending it as response
	out, err := json.Marshal(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(out)

	if err != nil {
		log.Println(err)
	}

}

// Authentication wrapper
func auth(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		var token string

		// Getting a set password
		pass := getEnvPassword()

		if len(pass) > 0 {

			// Getting a cookie
			cookie, err := req.Cookie("token")

			if err == nil {
				token = cookie.Value
			}

			// Validating a token
			var valid bool

			valid, err = ValidateToken(token, string(generateSecretKey(user.Name)))

			if err != nil {
				writeJson(res, ErrorResponse{Error: fmt.Sprintf("Internal security error")}, http.StatusInternalServerError)
				return
			}

			if !valid {
				// Authentication error
				writeJson(res, ErrorResponse{Error: fmt.Sprintf("Empty password")}, http.StatusUnauthorized)
				return
			}
		}

		next(res, req)
	})
}

// Getting passrord  from TODO_PASSWORD env variable
func getEnvPassword() string {

	envPassword := os.Getenv("TODO_PASSWORD")

	if envPassword == "" {
		envPassword = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiTmFtZSI6InVzZXIxIn0.NdVbWc1NZbtrS7IA4-3W-Tf8fg9pnk53vNHx9G1dU4A"
	}

	return envPassword
}

// JWT-token validation
func ValidateToken(tokenString string, secretKey string) (bool, error) {

	var valid bool

	claims := Claims{}

	jwtToken, err := jwt.ParseWithClaims(tokenString, &claims,

		func(token *jwt.Token) (interface{}, error) {
			// желательно проверять используемый метод
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

	if err != nil {

		return false, err
	}

	valid = jwtToken.Valid

	return valid, nil

}
