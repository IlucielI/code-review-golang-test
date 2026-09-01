package main

import (
	"html/template"
	"net/http"
)

// handleRender merender body post.
func (s *Server) handleRender(w http.ResponseWriter, r *http.Request) {
	body := r.URL.Query().Get("body")

	tmpl := template.Must(template.New("post").Parse("<article>{{.Body}}</article>"))
	_ = tmpl.Execute(w, struct {
		Body template.HTML
	}{
		Body: template.HTML(body),
	})
}
