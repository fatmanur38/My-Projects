package utils

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var Db *sql.DB

func InitDatabase() {
	var err error
	Db, err = sql.Open("sqlite3", "database/database.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	CreateTables()
}

func GetUsernameBySessionToken(sessionToken string) (string, error) {
	var username string
	var tokenExpires time.Time
	err := Db.QueryRow("SELECT username, token_expires FROM users WHERE session_token = ?", sessionToken).Scan(&username, &tokenExpires)
	if err != nil {
		return "", err
	}
	if tokenExpires.Before(time.Now()) {
		return "", fmt.Errorf("session token expired")
	}
	return username, nil
}

func InsertDefaultAdmin() {
	adminEmail := "admin@gmail.com"
	adminUsername := "admin"
	adminPassword := "admin123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error hashing admin password: %v", err)
	}

	_, err = Db.Exec("INSERT INTO users (username, email, password, role) VALUES (?, ?, ?, ?) ON CONFLICT(email) DO NOTHING", adminUsername, adminEmail, string(hashedPassword), "admin")
	if err != nil {
		log.Fatalf("Error inserting default admin: %v", err)
	}
}

func CreateTables() {
	createUsersTable := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE,
		username TEXT UNIQUE,
		name TEXT,
		about TEXT,
		usericon_url TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		password TEXT,
		session_token TEXT,
		token_expires DATETIME,
		role TEXT NOT NULL
	);`
	_, err := Db.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}

	createPostsTable := `CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		content TEXT,
		author_id INTEGER,
		image_url TEXT,
		category TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		view_count INTEGER DEFAULT 0,
		FOREIGN KEY (author_id) REFERENCES users(id)
	);`
	_, err = Db.Exec(createPostsTable)
	if err != nil {
		log.Fatalf("Error creating posts table: %v", err)
	}

	createCommentsTable := `CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER,
		content TEXT,
		author_id INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (author_id) REFERENCES users(id)
	);`
	_, err = Db.Exec(createCommentsTable)
	if err != nil {
		log.Fatalf("Error creating comments table: %v", err)
	}

	createVotesTable := `CREATE TABLE IF NOT EXISTS votes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		post_id INTEGER,
		vote INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (post_id) REFERENCES posts(id),
		UNIQUE(user_id, post_id)
	);`
	_, err = Db.Exec(createVotesTable)
	if err != nil {
		log.Fatalf("Error creating votes table: %v", err)
	}

	createCommentVotesTable := `CREATE TABLE IF NOT EXISTS comment_votes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		comment_id INTEGER,
		vote INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (comment_id) REFERENCES comments(id),
		UNIQUE(user_id, comment_id)
	);`
	_, err = Db.Exec(createCommentVotesTable)
	if err != nil {
		log.Fatalf("Error creating comment_votes table: %v", err)
	}

	createPostViewsTable := `CREATE TABLE IF NOT EXISTS post_views (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		post_id INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (post_id) REFERENCES posts(id),
		UNIQUE(user_id, post_id)
	);`
	_, err = Db.Exec(createPostViewsTable)
	if err != nil {
		log.Fatalf("Error creating post_views table: %v", err)
	}

	createPostCategoriesTable := `CREATE TABLE IF NOT EXISTS post_categories (
		post_id INTEGER,
		category TEXT,
		roles TEXT NOT NULL,
		FOREIGN KEY (post_id) REFERENCES posts(id),
		PRIMARY KEY (post_id, category)
	);`
	_, err = Db.Exec(createPostCategoriesTable)
	if err != nil {
		log.Fatalf("Error creating post_categories table: %v", err)
	}

	createRolesTable := `CREATE TABLE IF NOT EXISTS roles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);`
	_, err = Db.Exec(createRolesTable)
	if err != nil {
		log.Fatalf("Error creating roles table: %v", err)
	}

	createPermissionsTable := `CREATE TABLE IF NOT EXISTS permissions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		role_id INTEGER NOT NULL,
		permission TEXT NOT NULL,
		FOREIGN KEY (role_id) REFERENCES roles(id)
	);`
	_, err = Db.Exec(createPermissionsTable)
	if err != nil {
		log.Fatalf("Error creating permissions table: %v", err)
	}

	createModeratorRequestsTable := `CREATE TABLE IF NOT EXISTS moderator_requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		status TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`
	_, err = Db.Exec(createModeratorRequestsTable)
	if err != nil {
		log.Fatalf("Error creating moderator_requests table: %v", err)
	}

	createAllUserRequestsTable := `CREATE TABLE IF NOT EXISTS all_user_requests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    status TEXT DEFAULT 'Pending',
    FOREIGN KEY (user_id) REFERENCES users(id)
);`
	_, err = Db.Exec(createAllUserRequestsTable)
	if err != nil {
		log.Fatalf("Error creating admin_responses table: %v", err)
	}

	createAdminResponsesTable := `CREATE TABLE IF NOT EXISTS admin_responses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		request_id INTEGER NOT NULL,
		post_id INTEGER,
		moderator_id INTEGER,
		status TEXT DEFAULT 'Pending',
		reason TEXT,
		post_title TEXT,
		response TEXT NOT NULL,
		post_author_username TEXT,
		FOREIGN KEY (request_id) REFERENCES post_reports(id),
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (moderator_id) REFERENCES users(id)
	);`
	_, err = Db.Exec(createAdminResponsesTable)
	if err != nil {
		log.Fatalf("Error creating admin_responses table: %v", err)
	}

	createPostReportsTable := `CREATE TABLE IF NOT EXISTS post_reports (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER,
		moderator_id INTEGER,
		status TEXT DEFAULT 'Pending',
		reason TEXT,
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (moderator_id) REFERENCES users(id)
	);`
	_, err = Db.Exec(createPostReportsTable)
	if err != nil {
		log.Fatalf("Error creating post_reports table: %v", err)
	}
}
