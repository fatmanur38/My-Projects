package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"forum/utils"
)

func ReportPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderReportForm(w, r)
	} else if r.Method == http.MethodPost {
		processReport(w, r)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func renderReportForm(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("postId")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	data := struct {
		PostID string
	}{
		PostID: postID,
	}

	tmpl, err := template.ParseFiles("templates/report_post.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func processReport(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	postID := r.FormValue("postId")
	reason := r.FormValue("reason")
	moderatorID, _ := utils.GetUserIDFromCookie(r)

	if postID == "" || reason == "" {
		http.Error(w, "Post ID and reason are required", http.StatusBadRequest)
		return
	}

	if _, err := utils.Db.Exec("INSERT INTO post_reports (post_id, moderator_id, reason) VALUES (?, ?, ?)", postID, moderatorID, reason); err != nil {
		log.Println("Error reporting post:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ReviewReport handles the review of post reports.
func ReviewReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	reportID := r.FormValue("reportId")
	action := r.FormValue("action")
	redirectURL := r.FormValue("redirectURL")
	rejectionReason := r.FormValue("rejectionReason")

	err := processReportReview(reportID, action, rejectionReason)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// processReportReview processes the review of post reports.
func processReportReview(reportID, action, rejectionReason string) error {
	db := utils.Db

	var status, response, moderatorID, postAuthorUsername, postTitle string
	var postID int

	if action == "approve" {
		status = "Approved"
		response = "Approved"
	} else if action == "reject" {
		status = "Rejected"
		response = "Rejected: " + rejectionReason
	}

	// Retrieve the necessary information from the report
	row := db.QueryRow(`
        SELECT r.post_id, r.moderator_id, u.username, p.title
        FROM post_reports r
        JOIN posts p ON r.post_id = p.id
        JOIN users u ON p.author_id = u.id
        WHERE r.id = ?
    `, reportID)
	if err := row.Scan(&postID, &moderatorID, &postAuthorUsername, &postTitle); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No rows found for report ID: %v", reportID)
			return nil
		}
		return err
	}

	// Update the status of the report in the database
	_, err := db.Exec("UPDATE post_reports SET status = ? WHERE id = ?", status, reportID)
	if err != nil {
		return err
	}

	// Add the admin's response to the admin_responses table
	_, err = db.Exec(`
        INSERT INTO admin_responses (request_id, post_id, moderator_id, status, reason, response, post_author_username, post_title)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `, reportID, postID, moderatorID, status, rejectionReason, response, postAuthorUsername, postTitle)
	if err != nil {
		return err
	}

	// Delete the report from the post_reports table
	_, err = db.Exec("DELETE FROM post_reports WHERE id = ?", reportID)
	if err != nil {
		return err
	}

	return nil
}
