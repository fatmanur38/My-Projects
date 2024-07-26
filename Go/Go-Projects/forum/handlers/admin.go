package handlers

import (
	"html/template"
	"log"
	"net/http"

	"forum/models"
	"forum/utils"
)

func AdminPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/adminPage.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	requests, err := getAllRequests()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	users, err := getAllUsers()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	posts, err := getAllPosts()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	allRequests, err := getAllUserRequests()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	moderators, err := getAllModerators()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	reports, err := getAllReports()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, struct {
		ModeratorRequests []models.ModeratorRequest
		Users             []models.User
		Posts             []models.Post
		AllRequests       []models.UserRequest
		Moderators        []models.User
		Reports           []models.Report
	}{
		ModeratorRequests: requests,
		Users:             users,
		Posts:             posts,
		AllRequests:       allRequests,
		Moderators:        moderators,
		Reports:           reports,
	})
}

func getAllRequests() ([]models.ModeratorRequest, error) {
	rows, err := utils.Db.Query(`
        SELECT COALESCE(mr.id, 0), COALESCE(u.id, 0), COALESCE(u.username, ''), COALESCE(u.email, ''), COALESCE(u.role, ''), COALESCE(mr.status, '')
        FROM moderator_requests mr
        JOIN users u ON mr.user_id = u.id
        WHERE mr.status = 'Pending'
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.ModeratorRequest
	for rows.Next() {
		var req models.ModeratorRequest
		var user models.User
		if err := rows.Scan(&req.ID, &user.ID, &user.Username, &user.Email, &user.Role, &req.Status); err != nil {
			return nil, err
		}
		req.User = user
		requests = append(requests, req)
	}
	return requests, nil
}

func getAllUserRequests() ([]models.UserRequest, error) {
	rows, err := utils.Db.Query(`
        SELECT COALESCE(r.id, 0), COALESCE(u.id, 0), COALESCE(u.username, ''), COALESCE(u.email, ''), COALESCE(u.role, ''), COALESCE(r.status, '')
        FROM all_user_requests r
        JOIN users u ON r.user_id = u.id
        WHERE r.status = 'Pending'
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.UserRequest
	for rows.Next() {
		var req models.UserRequest
		var user models.User
		if err := rows.Scan(&req.ID, &user.ID, &user.Username, &user.Email, &user.Role, &req.Status); err != nil {
			return nil, err
		}
		req.User = user
		requests = append(requests, req)
	}
	return requests, nil
}

func getAllUsers() ([]models.User, error) {
	rows, err := utils.Db.Query("SELECT COALESCE(id, 0), COALESCE(username, ''), COALESCE(email, ''), COALESCE(role, '') FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func getAllPosts() ([]models.Post, error) {
	rows, err := utils.Db.Query(`
        SELECT COALESCE(p.id, 0), COALESCE(p.title, ''), COALESCE(p.content, ''), COALESCE(u.username, ''), COALESCE(p.created_at, ''), COALESCE(p.image_url, '')
        FROM posts p
        JOIN users u ON p.author_id = u.id
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var author string
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &author, &post.CreatedAt, &post.ImageURL); err != nil {
			return nil, err
		}
		post.Author = author
		posts = append(posts, post)
	}
	return posts, nil
}

func ApproveReject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.FormValue("userId")
	action := r.FormValue("action")
	redirectURL := r.FormValue("redirectURL")

	err := processRequest(userID, action)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func processRequest(userID string, action string) error {
	db := utils.Db

	switch action {
	case "approve":
		_, err := db.Exec("UPDATE users SET role = 'Moderator' WHERE id = ?", userID)
		if err != nil {
			return err
		}
	case "reject":
		// Optionally handle any specific logic for reject here
	}

	// Remove the user request from all_user_requests table regardless of action
	if _, err := db.Exec("DELETE FROM all_user_requests WHERE user_id = ?", userID); err != nil {
		return err
	}

	return nil
}

func getAllModerators() ([]models.User, error) {
	rows, err := utils.Db.Query(`
        SELECT COALESCE(id, 0), COALESCE(username, ''), COALESCE(email, '')
        FROM users
        WHERE role = 'Moderator'
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var moderators []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, err
		}
		moderators = append(moderators, user)
	}
	return moderators, nil
}

func getAllReports() ([]models.Report, error) {
	rows, err := utils.Db.Query(`
        SELECT COALESCE(r.id, 0), COALESCE(p.title, ''), COALESCE(u.username, ''), COALESCE(r.reason, ''), COALESCE(r.status, ''), COALESCE(p.image_url, '')
        FROM post_reports r
        JOIN posts p ON r.post_id = p.id
        JOIN users u ON r.moderator_id = u.id
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []models.Report
	for rows.Next() {
		var report models.Report
		if err := rows.Scan(&report.ID, &report.PostTitle, &report.ModeratorName, &report.Reason, &report.Status, &report.ImageURL); err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	return reports, nil
}

func ChangeUserRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.FormValue("userId")
	newRole := r.FormValue("newRole")
	redirectURL := r.FormValue("redirectURL")

	// Check if the user is not an admin before proceeding
	user, err := getUserByID(userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if user.Role == "Admin" {
		http.Error(w, "Cannot demote admin role", http.StatusBadRequest)
		return
	}

	if _, err := utils.Db.Exec("UPDATE users SET role = ? WHERE id = ?", newRole, userID); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func getUserByID(userID string) (models.User, error) {
	var user models.User
	err := utils.Db.QueryRow("SELECT id, username, email, role FROM users WHERE id = ?", userID).
		Scan(&user.ID, &user.Username, &user.Email, &user.Role)
	if err != nil {
		return user, err
	}
	return user, nil
}

// func updateUserRole(userID, newRole string) error {
// 	_, err := utils.Db.Exec("UPDATE users SET role = ? WHERE id = ?", newRole, userID)
// 	return err
// }

// func DeleteUser(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	userID := r.FormValue("userId")

// 	err := deleteUser(userID)
// 	if err != nil {
// 		log.Println(err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}

// 	http.Redirect(w, r, "/adminPage", http.StatusSeeOther)
// }

// func deleteUser(userID string) error {
// 	_, err := utils.Db.Exec("DELETE FROM users WHERE id = ?", userID)
// 	return err
// }
