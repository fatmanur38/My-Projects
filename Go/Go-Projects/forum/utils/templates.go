package utils

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Printf("error parsing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Using a buffer to ensure template execution completes without error before writing to the response.
	buffer := new(bytes.Buffer)
	if err := t.Execute(buffer, data); err != nil {
		log.Printf("error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buffer.WriteTo(w) // Only write to the ResponseWriter if there were no errors
}
