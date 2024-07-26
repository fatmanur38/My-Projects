package handlers

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"

	"forum/utils"
)

var config Config

type Config struct {
	GoogleClientID       string `json:"google_client_id"`
	GoogleClientSecret   string `json:"google_client_secret"`
	GitHubClientID       string `json:"github_client_id"`
	GitHubClientSecret   string `json:"github_client_secret"`
	FacebookClientID     string `json:"facebook_client_id"`
	FacebookClientSecret string `json:"facebook_client_secret"`
}

func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Failed to open config file: %s", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Failed to decode config file: %s", err)
	}
}

var (
	googleOauthConfig   *oauth2.Config
	githubOauthConfig   *oauth2.Config
	facebookOauthConfig *oauth2.Config
)

// CSRF (Siteler Arası İstek Sahtekarlığı) Saldırıları
var (
	oauthStateString         = "random"
	oauthStateStringGitHub   = "random"
	facebookOauthStateString = "random"
)

// Paket yüklenirken otomatik olarak çalışır.
func init() {
	loadConfig()
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8040/auth/google/callback",
		ClientID:     config.GoogleClientID,
		ClientSecret: config.GoogleClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
	githubOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8040/auth/github/callback",
		ClientID:     config.GitHubClientID,
		ClientSecret: config.GitHubClientSecret,
		Scopes:       []string{"read:user", "user:email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
	facebookOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8040/auth/facebook/callback",
		ClientID:     config.FacebookClientID,
		ClientSecret: config.FacebookClientSecret,
		Scopes: []string{
			"email",
		},
		Endpoint: facebook.Endpoint,
	}
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	oauthStateString = generateNonce()
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != oauthStateString {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	resp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", token.AccessToken))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var userID int
	err = utils.Db.QueryRow("SELECT id FROM users WHERE email = ?", googleUser.Email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = utils.Db.Exec("INSERT INTO users (username, email) VALUES (?, ?)", googleUser.Name, googleUser.Email)
			if err != nil {
				fmt.Printf("Error registering user: %v", err)
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}
			err = utils.Db.QueryRow("SELECT id FROM users WHERE email = ?", googleUser.Email).Scan(&userID)
			if err != nil {
				fmt.Printf("Error fetching new user ID: %v", err)
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}
		} else {
			fmt.Printf("Error querying user: %v", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
	}

	sessionToken := utils.GenerateSessionToken()
	expiration := time.Now().Add(24 * time.Hour)

	_, err = utils.Db.Exec("UPDATE users SET session_token = ?, token_expires = ? WHERE id = ?", sessionToken, expiration, userID)
	if err != nil {
		http.Error(w, "Failed to update session token.", http.StatusInternalServerError)
		return
	}

	utils.SetLoginCookie(w, userID, sessionToken, int(time.Until(expiration).Seconds()))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleFacebookLogin(w http.ResponseWriter, r *http.Request) {
	facebookOauthStateString = generateNonce()
	url := facebookOauthConfig.AuthCodeURL(facebookOauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleFacebookCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != facebookOauthStateString {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := facebookOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	resp, err := http.Get(fmt.Sprintf("https://graph.facebook.com/me?access_token=%s&fields=id,name,email", token.AccessToken))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	var facebookUser struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&facebookUser); err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var userID int
	err = utils.Db.QueryRow("SELECT id FROM users WHERE email = ?", facebookUser.Email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = utils.Db.Exec("INSERT INTO users (username, email) VALUES (?, ?)", facebookUser.Name, facebookUser.Email)
			if err != nil {
				fmt.Printf("Error registering user: %v", err)
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}
			err = utils.Db.QueryRow("SELECT id FROM users WHERE email = ?", facebookUser.Email).Scan(&userID)
			if err != nil {
				fmt.Printf("Error fetching new user ID: %v", err)
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}
		} else {
			fmt.Printf("Error querying user: %v", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
	}

	sessionToken := utils.GenerateSessionToken()
	expiration := time.Now().Add(24 * time.Hour)

	_, err = utils.Db.Exec("UPDATE users SET session_token = ?, token_expires = ? WHERE id = ?", sessionToken, expiration, userID)
	if err != nil {
		http.Error(w, "Failed to update session token.", http.StatusInternalServerError)
		return
	}

	utils.SetLoginCookie(w, userID, sessionToken, int(time.Until(expiration).Seconds()))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	oauthStateStringGitHub = generateNonce()
	url := githubOauthConfig.AuthCodeURL(oauthStateStringGitHub, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != oauthStateStringGitHub {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := githubOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("oauthConfGitHub.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	client := githubOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	var githubUser struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if githubUser.Email == "" {
		// Fetch the user's emails if the email field is empty
		resp, err := client.Get("https://api.github.com/user/emails")
		if err != nil {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		defer resp.Body.Close()

		var emails []struct {
			Email   string `json:"email"`
			Primary bool   `json:"primary"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		for _, email := range emails {
			if email.Primary {
				githubUser.Email = email.Email
				break
			}
		}
	}

	if githubUser.Email == "" {
		http.Error(w, "Unable to fetch email from GitHub", http.StatusInternalServerError)
		return
	}

	// Continue with the rest of the login/registration process
	var userID int
	err = utils.Db.QueryRow("SELECT id FROM users WHERE email = ?", githubUser.Email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = utils.Db.Exec("INSERT INTO users (username, email) VALUES (?, ?)", githubUser.Name, githubUser.Email)
			if err != nil {
				fmt.Printf("Error registering user: %v", err)
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}
			err = utils.Db.QueryRow("SELECT id FROM users WHERE email = ?", githubUser.Email).Scan(&userID)
			if err != nil {
				fmt.Printf("Error fetching new user ID: %v", err)
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}
		} else {
			fmt.Printf("Error querying user: %v", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
	}

	sessionToken := utils.GenerateSessionToken()
	expiration := time.Now().Add(24 * time.Hour)

	_, err = utils.Db.Exec("UPDATE users SET session_token = ?, token_expires = ? WHERE id = ?", sessionToken, expiration, userID)
	if err != nil {
		http.Error(w, "Failed to update session token.", http.StatusInternalServerError)
		return
	}

	utils.SetLoginCookie(w, userID, sessionToken, int(time.Until(expiration).Seconds()))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func generateNonce() string {
	nonce := make([]byte, 16)
	rand.Read(nonce)
	return hex.EncodeToString(nonce)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	redirect := r.URL.Query().Get("redirect")

	if r.Method == http.MethodGet {
		utils.RenderTemplate(w, "templates/login.html", map[string]interface{}{
			"Redirect": redirect,
		})
		return
	}

	if r.Method != http.MethodPost {
		utils.RenderTemplate(w, "templates/login.html", map[string]interface{}{
			"LoginErrorMsg": "Invalid request method",
		})
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" && password == "" {
		handleGitHubLogin(w, r)
		return
	}

	if email == "" && password == "" {
		handleFacebookLogin(w, r)
		return
	}

	var dbEmail, dbPassword, role string
	var userID int
	err := utils.Db.QueryRow("SELECT id, email, password, role FROM users WHERE email = ?", email).Scan(&userID, &dbEmail, &dbPassword, &role)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RenderTemplate(w, "templates/login.html", map[string]interface{}{
				"LoginErrorMsg": "User not found",
				"Redirect":      redirect,
			})
			return
		}
		utils.RenderTemplate(w, "templates/login.html", map[string]interface{}{
			"LoginErrorMsg": fmt.Sprintf("Error querying user: %v", err),
			"Redirect":      redirect,
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	if err != nil {
		utils.RenderTemplate(w, "templates/login.html", map[string]interface{}{
			"LoginErrorMsg": "Invalid email or password",
			"Redirect":      redirect,
		})
		return
	}

	sessionToken := utils.GenerateSessionToken()
	expiration := time.Now().Add(24 * time.Hour)

	_, err = utils.Db.Exec("UPDATE users SET session_token = ?, token_expires = ? WHERE id = ?", sessionToken, expiration, userID)
	if err != nil {
		http.Error(w, "Failed to update session token.", http.StatusInternalServerError)
		return
	}

	utils.SetLoginCookie(w, userID, sessionToken, int(time.Until(expiration).Seconds()))

	if redirect != "" {
		http.Redirect(w, r, redirect, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		utils.RenderTemplate(w, "templates/register.html", nil)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")
	role := r.FormValue("role")

	if password != confirmPassword {
		utils.RenderTemplate(w, "templates/register.html", map[string]interface{}{
			"RegisterErrorMsg": "Passwords do not match",
			"Username":         username,
			"Email":            email,
		})
		return
	}

	emailRegex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	matched, err := regexp.MatchString(emailRegex, email)
	if err != nil || !matched {
		utils.RenderTemplate(w, "templates/register.html", map[string]interface{}{
			"RegisterErrorMsg": "Invalid email format",
			"Username":         username,
			"Email":            email,
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		utils.RenderTemplate(w, "templates/register.html", map[string]interface{}{
			"RegisterErrorMsg": fmt.Sprintf("Error hashing password: %v", err),
			"Username":         username,
			"Email":            email,
		})
		return
	}

	if role == "moderator_request" {
		role = "User"
	}

	res, err := utils.Db.Exec("INSERT INTO users (username, email, password, role) VALUES (?, ?, ?, ?)", username, email, string(hashedPassword), role)
	if err != nil {
		utils.RenderTemplate(w, "templates/register.html", map[string]interface{}{
			"RegisterErrorMsg": fmt.Sprintf("Error registering user: %v", err),
			"Username":         username,
			"Email":            email,
		})
		return
	}

	userID, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.FormValue("role") == "moderator_request" {
		_, err := utils.Db.Exec("INSERT INTO all_user_requests (user_id, status) VALUES (?, 'Pending')", userID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromCookie(r)
	if err != nil {
		log.Println(err)
	}
	utils.SetLoginCookie(w, userID, "", -1)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
