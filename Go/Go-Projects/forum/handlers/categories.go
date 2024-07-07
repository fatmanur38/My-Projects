package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"io"
	"os"

	"forum/utils"
)

func uploadCategoryHandler(w http.ResponseWriter, r *http.Request, category string) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		uploadDir := "./uploads"

		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.MkdirAll(uploadDir, os.ModePerm)
		}

		filePath := filepath.Join(uploadDir, handler.Filename)
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0o666)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error saving file info to database: %v", err), http.StatusInternalServerError)
			return
		}
		defer f.Close()
		io.Copy(f, file)

		userID, err := utils.GetUserIDFromCookie(r)
		if err != nil {
			http.Error(w, "Not logged in", http.StatusForbidden)
			return
		}

		_, err = utils.Db.Exec("INSERT INTO posts (title, image_url, author_id, category) VALUES (?, ?, ?, ?)", title, "/"+filePath, userID, category)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/"+category, http.StatusSeeOther)
	}
}

func renderCategoryPage(w http.ResponseWriter, r *http.Request, category string) {
	loggedIn := utils.CheckLoginStatus(r)
	tmpl, err := template.ParseFiles(fmt.Sprintf("templates/%s.html", category))
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sortOrder := r.URL.Query().Get("sort")
	query := `SELECT p.id, p.title, p.image_url, u.username, u.usericon_url, p.created_at, 
              COALESCE(SUM(CASE v.vote WHEN 1 THEN 1 ELSE 0 END), 0) AS likes,
              COALESCE(SUM(CASE v.vote WHEN -1 THEN 1 ELSE 0 END), 0) AS dislikes,
              p.view_count
              FROM posts p
              JOIN users u ON p.author_id = u.id
              LEFT JOIN votes v ON p.id = v.post_id
              JOIN post_categories pc ON p.id = pc.post_id
              WHERE pc.category = ?
              GROUP BY p.id, u.username, u.usericon_url`

	switch sortOrder {
	case "newest":
		query += " ORDER BY p.created_at DESC"
	case "oldest":
		query += " ORDER BY p.created_at ASC"
	case "most_liked":
		query += " ORDER BY likes DESC"
	case "least_liked":
		query += " ORDER BY likes ASC"
	default:
		query += " ORDER BY p.created_at DESC "
	}

	rows, err := utils.Db.Query(query, category)
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.ImageURL, &post.AuthorName, &post.AuthorIcon, &post.CreatedAt, &post.Likes, &post.Dislikes, &post.ViewCount); err != nil {
			log.Printf("Error scanning post: %v", err)
			http.Error(w, "Error scanning post", http.StatusInternalServerError)
			return
		}
		post.FormattedCreatedAt = utils.ConvertToIstanbulTime(post.CreatedAt).Format("2006-01-02 15:04:05")
		posts = append(posts, post)
	}

	data := struct {
		Category  string
		Posts     []Post
		LoggedIn  bool
		SortOrder string
	}{
		Category:  category,
		Posts:     posts,
		LoggedIn:  loggedIn,
		SortOrder: sortOrder,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func uploadMoviesHandler(w http.ResponseWriter, r *http.Request) {
	uploadCategoryHandler(w, r, "movies")
}

func uploadTurkishHandler(w http.ResponseWriter, r *http.Request) {
	uploadCategoryHandler(w, r, "turkish")
}

func uploadScienceHandler(w http.ResponseWriter, r *http.Request) {
	uploadCategoryHandler(w, r, "science")
}

func uploadFoodHandler(w http.ResponseWriter, r *http.Request) {
	uploadCategoryHandler(w, r, "food")
}

func uploadTechnologyHandler(w http.ResponseWriter, r *http.Request) {
	uploadCategoryHandler(w, r, "technology")
}

func uploadHealthHandler(w http.ResponseWriter, r *http.Request) {
	uploadCategoryHandler(w, r, "health")
}

func moviesHandler(w http.ResponseWriter, r *http.Request) {
	renderCategoryPage(w, r, "movies")
}

func turkishHandler(w http.ResponseWriter, r *http.Request) {
	renderCategoryPage(w, r, "turkish")
}

func scienceHandler(w http.ResponseWriter, r *http.Request) {
	renderCategoryPage(w, r, "science")
}

func foodHandler(w http.ResponseWriter, r *http.Request) {
	renderCategoryPage(w, r, "food")
}

func technologyHandler(w http.ResponseWriter, r *http.Request) {
	renderCategoryPage(w, r, "technology")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	renderCategoryPage(w, r, "health")
}
