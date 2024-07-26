package models

type User struct {
	ID       int
	Username string
	Email    string
	Role     string
}

type Post struct {
	ID        int
	Title     string
	Content   string
	AuthorID  int
	Author    string
	CreatedAt string
	ImageURL  string
}

type ModeratorRequest struct {
	ID     int
	User   User
	Status string
	ImageURL  string
}

type UserRequest struct {
	ID     int
	User   User
	Status string
}

type Report struct {
	ID            int
	PostID        int
	ModeratorID   int
	Status        string
	Reason        string
	PostTitle     string
	ModeratorName string
	AuthorID      int
	ImageURL  string
}

type AdminResponse struct {
	ID                 int
	RequestID          int
	PostID             int
	ModeratorID        int
	Status             string
	Reason             string
	Response           string
	PostAuthorUsername string
	PostTitle          string // Yeni alan
}
