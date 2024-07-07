package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"forum/utils"

)

type User struct {
	ID       int
	Username string
	Email    string
	Password string
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.CheckLoginStatus(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	userID, err := utils.GetUserIDFromCookie(r)
	if err != nil {
		http.Error(w, "Not logged in", http.StatusForbidden)
		return
	}

	var username, email string
	var name, about, userIcon sql.NullString
	var createdAt time.Time
	err = utils.Db.QueryRow("SELECT username, email, name, about, usericon_url, created_at FROM users WHERE id = ?", userID).Scan(&username, &email, &name, &about, &userIcon, &createdAt)
	if err != nil {
		log.Printf("Error fetching user profile for userID %d: %v", userID, err)
		http.Error(w, "Error fetching user profile", http.StatusInternalServerError)
		return
	}

	posts, err := fetchPosts(userID, "p.author_id = ?")
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}

	likedPosts, err := fetchPosts(userID, "v.user_id = ? AND v.vote = 1")
	if err != nil {
		log.Printf("Error fetching liked posts: %v", err)
		http.Error(w, "Error fetching liked posts", http.StatusInternalServerError)
		return
	}

	dislikedPosts, err := fetchPosts(userID, "v.user_id = ? AND v.vote = -1")
	if err != nil {
		log.Printf("Error fetching disliked posts: %v", err)
		http.Error(w, "Error fetching disliked posts", http.StatusInternalServerError)
		return
	}

	commentedPosts, err := fetchPosts(userID, "p.id IN (SELECT post_id FROM comments WHERE author_id = ?)")
	if err != nil {
		log.Printf("Error fetching commented posts: %v", err)
		http.Error(w, "Error fetching commented posts", http.StatusInternalServerError)
		return
	}

	formattedJoinDate := utils.ConvertToIstanbulTime(createdAt).Format("2006-01-02")

	data := struct {
		LoggedIn       bool
		Posts          []Post
		LikedPosts     []Post
		DislikedPosts  []Post
		CommentedPosts []Post
		Username       string
		Email          string
		Name           string
		About          string
		UserIcon       string
		CreatedAt      string
	}{
		LoggedIn:       true,
		Posts:          posts,
		LikedPosts:     likedPosts,
		DislikedPosts:  dislikedPosts,
		CommentedPosts: commentedPosts,
		Username:       username,
		Email:          email,
		Name:           name.String,
		About:          about.String,
		UserIcon:       userIcon.String,
		CreatedAt:      formattedJoinDate,
	}

	funcMap := template.FuncMap{
		"split":  utils.Split,
	}

	tmpl, err := template.New("profile.html").Funcs(funcMap).ParseFiles("templates/profile.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}

func fetchPosts(userID int, condition string) ([]Post, error) {
	query := `SELECT p.id, p.title, p.image_url, u.username, u.usericon_url, p.created_at,
        COALESCE(SUM(CASE v.vote WHEN 1 THEN 1 ELSE 0 END), 0) AS likes,
        COALESCE(SUM(CASE v.vote WHEN -1 THEN 1 ELSE 0 END), 0) AS dislikes,
        (SELECT COUNT(*) FROM comments WHERE post_id = p.id) AS comment_count,
        GROUP_CONCAT(DISTINCT pc.category ORDER BY pc.category) AS categories,
        p.view_count
        FROM posts p
        JOIN users u ON p.author_id = u.id
        LEFT JOIN votes v ON p.id = v.post_id
        LEFT JOIN post_categories pc ON p.id = pc.post_id
        WHERE ` + condition + `
        GROUP BY p.id
        ORDER BY p.created_at DESC`

	rows, err := utils.Db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.ImageURL, &post.AuthorName, &post.AuthorIcon, &post.CreatedAt, &post.Likes, &post.Dislikes, &post.CommentCount, &post.Category, &post.ViewCount); err != nil {
			return nil, err
		}
		post.FormattedCreatedAt = utils.ConvertToIstanbulTime(post.CreatedAt).Format("2006-01-02 15:04:05")
		posts = append(posts, post)
	}
	return posts, nil
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var username, email string
		var name, about, userIcon sql.NullString

		userID, err := utils.GetUserIDFromCookie(r)
		if err != nil {
			log.Printf("Failed to get user ID from cookie: %v", err)
			http.Error(w, "Not logged in", http.StatusForbidden)
			return
		}

		err = utils.Db.QueryRow("SELECT username, email, name, about, usericon_url FROM users WHERE id = ?", userID).Scan(&username, &email, &name, &about, &userIcon)
		if err != nil {
			log.Printf("Error fetching user profile for userID %d: %v", userID, err)
			http.Error(w, "Error fetching user profile", http.StatusInternalServerError)
			return
		}

		session, err := utils.Store.Get(r, "session-name")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Values["username"] = username
		session.Values["email"] = email
		session.Values["name"] = name.String
		session.Values["about"] = about.String
		session.Values["userIcon"] = userIcon.String
		session.Save(r, w)

		data := struct {
			Username string
			Email    string
			Name     string
			About    string
			UserIcon string
			ErrorMsg string
		}{
			Username: username,
			Email:    email,
			Name:     name.String,
			About:    about.String,
			UserIcon: userIcon.String,
		}

		utils.RenderTemplate(w, "templates/edit.html", data)
	} else if r.Method == http.MethodPost {
		userID, err := utils.GetUserIDFromCookie(r)
		if err != nil {
			http.Error(w, "Not logged in", http.StatusForbidden)
			return
		}

		r.ParseMultipartForm(10 << 20)

		name := r.FormValue("name")
		about := r.FormValue("about")
		email := r.FormValue("email")
		username := r.FormValue("username")

		var filePath string
		file, handler, err := r.FormFile("file")
		if err == nil {
			defer file.Close()

			uploadDir := "./uploads"
			if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
				os.MkdirAll(uploadDir, os.ModePerm)
			}

			filePath = filepath.Join(uploadDir, handler.Filename)
			f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0o666)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error saving file: %v", err), http.StatusInternalServerError)
				return
			}
			defer f.Close()
			io.Copy(f, file)
			filePath = "/" + filePath
		}

		var exists bool
		err = utils.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ? AND id != ?)", email, userID).Scan(&exists)
		if err != nil {
			log.Printf("Error checking email uniqueness: %v", err)
			http.Error(w, "Error updating user profile", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "Email already in use by another account.", http.StatusConflict)
			return
		}

		err = utils.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ? AND id != ?)", username, userID).Scan(&exists)
		if err != nil {
			log.Printf("Error checking username uniqueness: %v", err)
			http.Error(w, "Error updating user profile", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "Username already in use by another account.", http.StatusConflict)
			return
		}

		var updateQuery string
		var args []interface{}
		if filePath != "" {
			updateQuery = "UPDATE users SET name = ?, about = ?, email = ?, username = ?, usericon_url = ? WHERE id = ?"
			args = []interface{}{name, about, email, username, filePath, userID}
		} else {
			updateQuery = "UPDATE users SET name = ?, about = ?, email = ?, username = ? WHERE id = ?"
			args = []interface{}{name, about, email, username, userID}
		}
		_, err = utils.Db.Exec(updateQuery, args...)
		if err != nil {
			log.Printf("Error updating user profile: %v", err)
			http.Error(w, "Error updating user profile: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func deletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if !utils.CheckLoginStatus(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, err := utils.GetUserIDFromCookie(r)
	if err != nil {
		log.Printf("Error retrieving user ID: %v", err)
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	postID := r.FormValue("postID")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	// Check if the current user is the author of the post
	var authorID int
	err = utils.Db.QueryRow("SELECT author_id FROM posts WHERE id = ?", postID).Scan(&authorID)
	if err != nil {
		log.Printf("Error fetching post: %v", err)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	if authorID != userID {
		http.Error(w, "Unauthorized to delete this post", http.StatusUnauthorized)
		return
	}

	_, err = utils.Db.Exec("DELETE FROM posts WHERE id = ?", postID)
	if err != nil {
		log.Printf("Error deleting post: %v", err)
		http.Error(w, "Error deleting post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
