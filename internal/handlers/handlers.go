package handlers

import (
	"encoding/json"
	"final-project/internal/usecase"
	"net/http"
)

type Handler struct {
	uc *usecase.Usecase
}

func NewHandler(uc *usecase.Usecase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) GetRate(w http.ResponseWriter, r *http.Request) {
	result, err := h.uc.GetData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	js, err := json.MarshalIndent(result, " ", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
