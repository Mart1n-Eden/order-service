package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"order-service/internal/service"
)

type Handler struct {
	//consumer *consumer.Consumer
	//cache *cache.Cache
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) InitRoutes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", h.HomePage)
	r.HandleFunc("/order", h.GetOrder).Methods("GET")

	return r
}
