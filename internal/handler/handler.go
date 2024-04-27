package handler

import (
	"html/template"
	"net/http"
)

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/index.html")
	if err != nil {
		http.Error(w, "Sorry..", http.StatusInternalServerError)
	}

	tmpl.Execute(w, nil)
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if content, ok := h.service.GetOrder(id); !ok {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(content)
	}
}
