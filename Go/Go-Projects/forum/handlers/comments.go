package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"forum/utils"
)

func submitCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID, err := utils.GetUserIDFromCookie(r)
	if err != nil || !utils.CheckLoginStatus(r) {
		http.Redirect(w, r, "/login?redirect="+url.QueryEscape(r.RequestURI), http.StatusSeeOther)
		return
	}

	postIDStr := r.FormValue("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	content := r.FormValue("comment")

	if content == "" {
		comments, err := fetchCommentsByID(postID, userID)
		if err != nil {
			http.Error(w, "Error fetching comments", http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"PostID":          postID,
			"Comments":        comments,
			"CommentErrorMsg": "Bu kısım boş geçilemez!",
		}

		utils.RenderTemplate(w, "templates/view_post.html", data)
		return
	}

	_, err = utils.Db.Exec("INSERT INTO comments (post_id, content, author_id) VALUES (?, ?, ?)", postID, content, userID)
	if err != nil {
		http.Error(w, "Error saving comment", http.StatusInternalServerError)
		return
	}

	comments, err := fetchCommentsByID(postID, userID)
	if err != nil {
		http.Error(w, "Error fetching comments", http.StatusInternalServerError)
		return
	}

	response := struct {
		Success  bool      `json:"success"`
		Comments []Comment `json:"comments"`
	}{
		Success:  true,
		Comments: comments,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func voteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID, err := utils.GetUserIDFromCookie(r)
	if err != nil || !utils.CheckLoginStatus(r) {
		http.Error(w, "Not logged in", http.StatusForbidden)
		return
	}

	commentID, err := strconv.Atoi(r.FormValue("comment_id"))
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	voteValue, err := strconv.Atoi(r.FormValue("vote"))
	if err != nil {
		http.Error(w, "Invalid vote value", http.StatusBadRequest)
		return
	}

	var existingVote int
	err = utils.Db.QueryRow("SELECT vote FROM comment_votes WHERE user_id = ? AND comment_id = ?", userID, commentID).Scan(&existingVote)

	if err == sql.ErrNoRows {
		_, err = utils.Db.Exec("INSERT INTO comment_votes (user_id, comment_id, vote) VALUES (?, ?, ?)", userID, commentID, voteValue)
	} else if err == nil {
		if existingVote == voteValue {
			_, err = utils.Db.Exec("DELETE FROM comment_votes WHERE user_id = ? AND comment_id = ?", userID, commentID)
		} else {
			_, err = utils.Db.Exec("UPDATE comment_votes SET vote = ? WHERE user_id = ? AND comment_id = ?", voteValue, userID, commentID)
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
        FROM comment_votes WHERE comment_id = ?`, commentID).Scan(&likes, &dislikes)
	if err != nil {
		http.Error(w, "Error fetching vote counts", http.StatusInternalServerError)
		return
	}

	response := map[string]int{"likes": likes, "dislikes": dislikes}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID, err := utils.GetUserIDFromCookie(r)
	if err != nil || !utils.CheckLoginStatus(r) {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	commentID, err := strconv.Atoi(r.FormValue("comment_id"))
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// Comment'ın yazarını ve postID'yi al
	var authorID, postID int
	err = utils.Db.QueryRow("SELECT author_id, post_id FROM comments WHERE id = ?", commentID).Scan(&authorID, &postID)
	if err != nil {
		http.Error(w, "Error fetching comment", http.StatusInternalServerError)
		return
	}

	// Yorumun yazarı mevcut kullanıcı mı diye kontrol et
	if userID != authorID {
		http.Error(w, "Unauthorized to delete this comment", http.StatusForbidden)
		return
	}

	// Yorumu sil
	_, err = utils.Db.Exec("DELETE FROM comments WHERE id = ?", commentID)
	if err != nil {
		http.Error(w, "Error deleting comment", http.StatusInternalServerError)
		return
	}

	// Yorum silindikten sonra ilgili post sayfasına yönlendir
	data := map[string]interface{}{
		"PostID": postID,
	}
	utils.RenderTemplate(w, "templates/view_post.html", data)
}
