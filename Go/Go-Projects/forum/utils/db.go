package utils

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func InitDatabase() {
	var err error
	Db, err = sql.Open("sqlite3", "database/database.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
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
		token_expires DATETIME
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
        FOREIGN KEY (post_id) REFERENCES posts(id),
        PRIMARY KEY (post_id, category)
    );`
	_, err = Db.Exec(createPostCategoriesTable)
	if err != nil {
		log.Fatalf("Error creating post_categories table: %v", err)
	}
}
