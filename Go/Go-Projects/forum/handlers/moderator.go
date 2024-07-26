package handlers

import (
	"html/template"
	"log"
	"net/http"

	"forum/models"
	"forum/utils"
)

func ModeratorPage(w http.ResponseWriter, r *http.Request) {
	moderatorID, err := utils.GetUserIDFromCookie(r)
	if err != nil {
		log.Println("Error getting user ID from cookie:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tmpl, err := template.ParseFiles("templates/moderatorPage.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	responses, err := getAdminResponsesForModerator(moderatorID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, struct {
		AdminResponses []models.AdminResponse
	}{
		AdminResponses: responses,
	})
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func getAdminResponsesForModerator(moderatorID int) ([]models.AdminResponse, error) {
	rows, err := utils.Db.Query(`
        SELECT ar.id, ar.request_id, ar.post_id, ar.moderator_id, ar.status, ar.reason, ar.response, ar.post_title
        FROM admin_responses ar
        WHERE ar.moderator_id = ?
    `, moderatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var responses []models.AdminResponse
	for rows.Next() {
		var response models.AdminResponse
		if err := rows.Scan(&response.ID, &response.RequestID, &response.PostID, &response.ModeratorID, &response.Status, &response.Reason, &response.Response, &response.PostTitle); err != nil {
			return nil, err
		}
		responses = append(responses, response)
	}
	return responses, nil
}
