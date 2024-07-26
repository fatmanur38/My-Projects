package handlers

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"forum/utils"
)

type Post struct {
	ID                 int
	Title              string
	ImageURL           string
	Content            string
	AuthorName         string
	AuthorIcon         sql.NullString
	Likes              int
	Dislikes           int
	UserIcon           string
	Category           sql.NullString
	CreatedAt          time.Time
	FormattedCreatedAt string
	CommentCount       int
	ViewCount          int
	IsAdminOrModerator bool
	IsModerator        bool
}

type Comment struct {
	ID         int
	Content    string
	AuthorName string
	Likes      int
	Dislikes   int
	CanDelete  bool
	AuthorID   int
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		utils.RenderTemplate(w, "templates/uploadForm.html", nil)
		return
	}

	if r.Method != http.MethodPost {
		utils.RenderTemplate(w, "templates/uploadForm.html", map[string]interface{}{"Error": "Invalid request method"})
		return
	}

	r.ParseMultipartForm(20 << 20) // 20 MB limit

	title := r.FormValue("title")
	content := r.FormValue("content")
	categories := r.Form["categories"]

	file, handler, err := r.FormFile("file")
	var filePath string
	if err == nil { // If file is selected
		defer file.Close()

		// Check file size
		if handler.Size > 20<<20 { // 20 MB
			utils.RenderTemplate(w, "templates/uploadForm.html", map[string]interface{}{"Error": "File is bigger than 20 MB"})
			return
		}

		// Check file type
		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/gif":  true,
		}
		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			utils.RenderTemplate(w, "templates/uploadForm.html", map[string]interface{}{"Error": "Error reading file"})
			return
		}
		fileType := http.DetectContentType(buff)
		if !allowedTypes[fileType] {
			utils.RenderTemplate(w, "templates/uploadForm.html", map[string]interface{}{"Error": "Invalid file type"})
			return
		}
		// Reset file pointer
		file.Seek(0, 0)
		// If the file passes the checks
		// it's saved to the "uploads" directory.
		uploadDir := "./uploads"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.MkdirAll(uploadDir, os.ModePerm)
		}
		// The directory is created if it doesn't exist.
		filePath = filepath.Join(uploadDir, handler.Filename)
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0o666)
		if err != nil {
			utils.RenderTemplate(w, "templates/uploadForm.html", map[string]interface{}{"Error": "Error saving file"})
			return
		}
		defer f.Close()

		if _, err = io.Copy(f, file); err != nil {
			utils.RenderTemplate(w, "templates/uploadForm.html", map[string]interface{}{"Error": "Error writing file"})
			return
		}
	} else {
		// Use default image path if no file is selected
		filePath = "static/images/default-image.png"
	}

	userID, err := utils.GetUserIDFromCookie(r)
	if err != nil {
		utils.RenderTemplate(w, "templates/uploadForm.html", map[string]interface{}{"Error": "Not logged in"})
		return
	}
	// It inserts the post's details
	result, err := utils.Db.Exec("INSERT INTO posts (title, image_url, author_id, content) VALUES (?, ?, ?, ?)", title, "/"+filePath, userID, content)
	if err != nil {
		log.Printf("Error inserting post: %v", err)
		utils.RenderTemplate(w, "templates/uploadForm.html", map[string]interface{}{"Error": "Database insert error"})
		return
	}

	postID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		utils.RenderTemplate(w, "templates/uploadForm.html", map[string]interface{}{"Error": "Error getting post ID"})
		return
	}

	// Fetch the user's role
	role, err := utils.GetUserRole(userID)
	if err != nil {
		log.Printf("Error getting user role: %v", err)
		utils.RenderTemplate(w, "templates/uploadForm.html", map[string]interface{}{"Error": "Error getting user role"})
		return
	}

	// For each selected category, it inserts a record into the "post_categories" table with the appropriate role
	for _, category := range categories {
		_, err = utils.Db.Exec("INSERT INTO post_categories (post_id, category, roles) VALUES (?, ?, ?)", postID, category, role)
		if err != nil {
			log.Printf("Error inserting post category: %v", err)
			utils.RenderTemplate(w, "templates/uploadForm.html", map[string]interface{}{"Error": "Database insert error"})
			return
		}
	}

	// it redirects the user to the main page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func viewPostHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userID, err := utils.GetUserIDFromCookie(r)
	var viewed bool
	if err == nil {
		err = utils.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM post_views WHERE user_id = ? AND post_id = ?)", userID, postID).Scan(&viewed)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if !viewed {
			_, err = utils.Db.Exec("INSERT INTO post_views (user_id, post_id) VALUES (?, ?)", userID, postID)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}

			_, err = utils.Db.Exec("UPDATE posts SET view_count = view_count + 1 WHERE id = ?", postID)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
		}
	}

	var post Post
	err = utils.Db.QueryRow(`
    SELECT p.id, p.title, p.image_url, p.content, u.username, p.created_at, u.usericon_url,
    COALESCE(SUM(CASE v.vote WHEN 1 THEN 1 ELSE 0 END), 0) AS likes,
    COALESCE(SUM(CASE v.vote WHEN -1 THEN 1 ELSE 0 END), 0) AS dislikes,
    p.view_count
    FROM posts p
    JOIN users u ON p.author_id = u.id
    LEFT JOIN votes v ON p.id = v.post_id
    WHERE p.id = ?
    GROUP BY p.id, u.username, u.usericon_url`, postID).Scan(&post.ID, &post.Title, &post.ImageURL, &post.Content, &post.AuthorName, &post.CreatedAt, &post.AuthorIcon, &post.Likes, &post.Dislikes, &post.ViewCount)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
		} else {
			http.Error(w, "Error fetching post details", http.StatusInternalServerError)
		}
		return
	}

	// Remove <p> tags from the content
	post.Content = strings.ReplaceAll(post.Content, "<p>", "")
	post.Content = strings.ReplaceAll(post.Content, "</p>", "")

	comments, err := fetchCommentsByID(postID, userID)
	if err != nil {
		http.Error(w, "Error fetching comments", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"PostID":     post.ID,
		"Title":      post.Title,
		"Image":      post.ImageURL,
		"Content":    post.Content,
		"AuthorName": post.AuthorName,
		"CreatedAt":  post.CreatedAt,
		"Likes":      post.Likes,
		"Dislikes":   post.Dislikes,
		"ViewCount":  post.ViewCount,
		"AuthorIcon": post.AuthorIcon,
		"Comments":   comments,
		"LoggedIn":   utils.CheckLoginStatus(r),
	}

	utils.RenderTemplate(w, "templates/view_post.html", data)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	loggedIn := utils.CheckLoginStatus(r)
	isadmin := utils.IsAdmin(r)
	ismoderator := utils.IsModerator(r)
	isadminormoderator := isadmin || ismoderator

	tmpl := template.Must(template.New("index.html").Funcs(template.FuncMap{
		"split": utils.Split,
	}).ParseFiles("templates/index.html"))

	sortOrder := r.URL.Query().Get("sort")

	query := `SELECT p.id, p.title, p.image_url, u.username, u.usericon_url, p.created_at,
                      COALESCE(likes.likes, 0) AS likes,
                      COALESCE(dislikes.dislikes, 0) AS dislikes,
                      GROUP_CONCAT(DISTINCT pc.category) AS categories,
                      (SELECT COUNT(*) FROM comments WHERE post_id = p.id) AS comment_count,
                      p.view_count
              FROM posts p
              JOIN users u ON p.author_id = u.id
              LEFT JOIN (SELECT post_id, COUNT(*) AS likes FROM votes WHERE vote = 1 GROUP BY post_id) likes ON p.id = likes.post_id
              LEFT JOIN (SELECT post_id, COUNT(*) AS dislikes FROM votes WHERE vote = -1 GROUP BY post_id) dislikes ON p.id = dislikes.post_id
              LEFT JOIN post_categories pc ON p.id = pc.post_id
              GROUP BY p.id, u.username, u.usericon_url, likes.likes, dislikes.dislikes`

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

	rows, err := utils.Db.Query(query)
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var categories sql.NullString
		if err := rows.Scan(&post.ID, &post.Title, &post.ImageURL, &post.AuthorName, &post.AuthorIcon, &post.CreatedAt, &post.Likes, &post.Dislikes, &categories, &post.CommentCount, &post.ViewCount); err != nil {
			log.Printf("Error scanning post: %v", err)
			http.Error(w, "Error scanning post", http.StatusInternalServerError)
			return
		}
		post.FormattedCreatedAt = utils.ConvertToIstanbulTime(post.CreatedAt).Format("2006-01-02 15:04:05")
		post.Category = categories
		posts = append(posts, post)
	}

	data := struct {
		LoggedIn           bool
		Posts              []Post
		SortOrder          string
		IsAdminOrModerator bool
	}{
		LoggedIn:           loggedIn,
		Posts:              posts,
		SortOrder:          sortOrder,
		IsAdminOrModerator: isadminormoderator,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if the user is logged in
	if !utils.CheckLoginStatus(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Check if the user is an admin or moderator
	if !utils.IsAdmin(r) && !utils.IsModerator(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Get the post ID from the form
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	postID := r.FormValue("postID")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	// Delete the post from the database
	_, err = utils.Db.Exec("DELETE FROM posts WHERE id = ?", postID)
	if err != nil {
		log.Printf("Error deleting post: %v", err)
		http.Error(w, "Error deleting post", http.StatusInternalServerError)
		return
	}

	// Redirect to the homepage after successful deletion
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func fetchCommentsByID(postID int, currentUserID int) ([]Comment, error) {
	var comments []Comment
	rows, err := utils.Db.Query(`
        SELECT c.id, c.author_id, c.content, u.username, 
        COALESCE(SUM(CASE cv.vote WHEN 1 THEN 1 ELSE 0 END), 0) AS likes,
        COALESCE(SUM(CASE cv.vote WHEN -1 THEN 1 ELSE 0 END), 0) AS dislikes
        FROM comments c
        JOIN users u ON c.author_id = u.id
        LEFT JOIN comment_votes cv ON c.id = cv.comment_id
        WHERE c.post_id = ?
        GROUP BY c.id, c.author_id, u.username
        ORDER BY c.created_at ASC
    `, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.AuthorID, &comment.Content, &comment.AuthorName, &comment.Likes, &comment.Dislikes); err != nil {
			return nil, err
		}
		comment.CanDelete = (comment.AuthorID == currentUserID) // Yorumu silme yetkisini kontrol et
		comments = append(comments, comment)
	}
	return comments, nil
}
