package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"forum/utils"
)

func voteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID, err := utils.GetUserIDFromCookie(r)
	if err != nil || !utils.CheckLoginStatus(r) {
		http.Error(w, "Not logged in", http.StatusForbidden)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	voteValue, err := strconv.Atoi(r.FormValue("vote"))
	if err != nil {
		http.Error(w, "Invalid vote value", http.StatusBadRequest)
		return
	}

	var existingVote int
	err = utils.Db.QueryRow("SELECT vote FROM votes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&existingVote)

	if err == sql.ErrNoRows {
		_, err = utils.Db.Exec("INSERT INTO votes (user_id, post_id, vote) VALUES (?, ?, ?)", userID, postID, voteValue)
	} else if err == nil {
		if existingVote == voteValue {
			_, err = utils.Db.Exec("DELETE FROM votes WHERE user_id = ? AND post_id = ?", userID, postID)
		} else {
			_, err = utils.Db.Exec("UPDATE votes SET vote = ? WHERE user_id = ? AND post_id = ?", voteValue, userID, postID)
		}
	}

	if err != nil {
		http.Error(w, "Error processing vote", http.StatusInternalServerError)
		return
	}

	var likes, dislikes int
	err = utils.Db.QueryRow(`SELECT
        COALESCE(SUM(CASE vote WHEN 1 THEN 1 ELSE 0 END), 0) AS likes,
        COALESCE(SUM(CASE vote WHEN -1 THEN 1 ELSE 0 END), 0) AS dislikes
        FROM votes WHERE post_id = ?`, postID).Scan(&likes, &dislikes)
	if err != nil {
		http.Error(w, "Error fetching vote counts", http.StatusInternalServerError)
		return
	}

	response := map[string]int{"likes": likes, "dislikes": dislikes}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func viewCountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID, err := utils.GetUserIDFromCookie(r)
	if err != nil || !utils.CheckLoginStatus(r) {
		http.Error(w, "Not logged in", http.StatusForbidden)
		return
	}

	var requestData struct {
		PostID int `json:"post_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	postID := requestData.PostID

	var viewed bool
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

	w.WriteHeader(http.StatusOK)
}
