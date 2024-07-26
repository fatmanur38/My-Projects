package utils

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("something-very-secret"))

func SetLoginCookie(w http.ResponseWriter, userID int, sessionToken string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:   "loggedin",
		Value:  "true",
		Path:   "/",
		MaxAge: maxAge,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "userid",
		Value:  fmt.Sprintf("%d", userID),
		Path:   "/",
		MaxAge: maxAge,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  sessionToken,
		Path:   "/",
		MaxAge: maxAge,
	})
}

func GetUserIDFromCookie(r *http.Request) (int, error) {
	cookie, err := r.Cookie("userid")
	if err != nil {
		return 0, err
	}
	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func CheckLoginStatus(r *http.Request) bool {
	loggedinCookie, err := r.Cookie("loggedin")
	if err != nil || loggedinCookie.Value != "true" {
		return false
	}

	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		return false
	}

	var dbSessionToken string
	var tokenExpires time.Time
	userID, err := GetUserIDFromCookie(r)
	if err != nil {
		return false
	}

	err = Db.QueryRow("SELECT session_token, token_expires FROM users WHERE id = ?", userID).Scan(&dbSessionToken, &tokenExpires)
	if err != nil || dbSessionToken != sessionToken.Value || tokenExpires.Before(time.Now()) {
		return false
	}

	return true
}

func GetUserRole(userID int) (string, error) {
	var role string
	query := "SELECT role FROM users WHERE id = ?"

	err := Db.QueryRow(query, userID).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // User not found
		}
		log.Printf("Error retrieving user role: %v", err)
		return "", err
	}

	return role, nil
}

func IsAdmin(r *http.Request) bool {
	userID, err := GetUserIDFromCookie(r)
	if err != nil {
		return false
	}

	var role string
	err = Db.QueryRow("SELECT role FROM users WHERE id = ?", userID).Scan(&role)
	if err != nil {
		log.Println(err)
		return false
	}

	return role == "admin"
}

func IsModerator(r *http.Request) bool {
	userID, err := GetUserIDFromCookie(r)
	if err != nil {
		return false
	}

	var role string
	err = Db.QueryRow("SELECT role FROM users WHERE id = ?", userID).Scan(&role)
	if err != nil {
		log.Println(err)
		return false
	}

	return role == "Moderator"
}
