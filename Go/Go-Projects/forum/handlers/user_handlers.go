package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	"forum/utils"
	"golang.org/x/crypto/bcrypt"
)

var tmpl = template.Must(template.ParseFiles("templates/sifre.html"))

func ServePasswordChangePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	session, err := utils.Store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Username": session.Values["username"],
		"Email":    session.Values["email"],
		"Name":     session.Values["name"],
		"About":    session.Values["about"],
		"UserIcon": session.Values["userIcon"],
	}
	tmpl.Execute(w, data)
}

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Unsupported request method.", http.StatusMethodNotAllowed)
		return
	}

	session, err := utils.Store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Error fetching session", http.StatusInternalServerError)
		return
	}

	sessionUserID, err := utils.GetUserIDFromCookie(r)
	if err != nil || sessionUserID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	oldPassword := r.FormValue("old_password")
	newPassword := r.FormValue("new_password")
	confirmPassword := r.FormValue("confirm_new_password")

	data := map[string]interface{}{
		"Username": session.Values["username"],
		"Email":    session.Values["email"],
		"Name":     session.Values["name"],
		"About":    session.Values["about"],
		"UserIcon": session.Values["userIcon"],
	}

	if newPassword != confirmPassword {
		data["ErrorMessage"] = "New passwords do not match."
		tmpl.Execute(w, data)
		return
	}

	var dbPasswordHash string
	err = utils.Db.QueryRow("SELECT password FROM users WHERE id = ?", sessionUserID).Scan(&dbPasswordHash)
	if err != nil {
		data["ErrorMessage"] = "User not found or database error."
		tmpl.Execute(w, data)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbPasswordHash), []byte(oldPassword))
	if err != nil {
		data["ErrorMessage"] = "Old password is incorrect."
		tmpl.Execute(w, data)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		data["ErrorMessage"] = "Error while hashing password."
		tmpl.Execute(w, data)
		return
	}

	_, err = utils.Db.Exec("UPDATE users SET password = ? WHERE id = ?", string(hashedPassword), sessionUserID)
	if err != nil {
		data["ErrorMessage"] = "Failed to update password."
		tmpl.Execute(w, data)
		return
	}

	http.Redirect(w, r, "/edit", http.StatusSeeOther)
}

func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	var user struct {
		Username sql.NullString
		Name     sql.NullString
		About    sql.NullString
		UserIcon sql.NullString
	}

	err := utils.Db.QueryRow("SELECT username, name, about, usericon_url FROM users WHERE username = ?", username).Scan(&user.Username, &user.Name, &user.About, &user.UserIcon)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		log.Printf("Error fetching user profile for username %s: %v", username, err)
		http.Error(w, "Error fetching user profile", http.StatusInternalServerError)
		return
	}

	var authorID int
	err = utils.Db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&authorID)
	if err != nil {
		log.Printf("Error fetching user ID for username %s: %v", username, err)
		http.Error(w, "Error fetching user ID", http.StatusInternalServerError)
		return
	}

	rows, err := utils.Db.Query(`
    SELECT p.id, p.title, p.image_url, p.created_at,
    COALESCE(SUM(CASE v.vote WHEN 1 THEN 1 ELSE 0 END), 0) AS likes,
    COALESCE(SUM(CASE v.vote WHEN -1 THEN 1 ELSE 0 END), 0) AS dislikes,
    GROUP_CONCAT(DISTINCT pc.category) AS categories
    FROM posts p
    LEFT JOIN votes v ON p.id = v.post_id
    LEFT JOIN post_categories pc ON p.id = pc.post_id
    WHERE p.author_id = ?
    GROUP BY p.id, p.title, p.image_url, p.created_at
    ORDER BY p.created_at DESC`, authorID)
	if err != nil {
		log.Printf("Error fetching posts for user %s: %v", username, err)
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.ImageURL, &post.CreatedAt, &post.Likes, &post.Dislikes, &post.Category); err != nil {
			log.Printf("Error scanning post for user %s: %v", username, err)
			http.Error(w, "Error scanning post", http.StatusInternalServerError)
			return
		}
		post.FormattedCreatedAt = utils.ConvertToIstanbulTime(post.CreatedAt).Format("2006-01-02 15:04:05")
		posts = append(posts, post)
	}

	data := struct {
		Username string
		Name     string
		About    string
		UserIcon sql.NullString
		Posts    []Post
		LoggedIn bool
	}{
		Username: user.Username.String,
		Name:     user.Name.String,
		About:    user.About.String,
		UserIcon: user.UserIcon,
		Posts:    posts,
		LoggedIn: utils.CheckLoginStatus(r),
	}

	funcMap := template.FuncMap{
		"split":  utils.Split,
	}

	tmpl, err := template.New("usersprofile.html").Funcs(funcMap).ParseFiles("templates/usersprofile.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}
