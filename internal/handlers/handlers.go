package handlers

import (
	"encoding/json"
	"final-project/internal/models"
	"final-project/internal/usecase"
	"fmt"
	"net/http"
	"sort"
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

	var tmpInternalResult []string
	var tmpFinalResult [][]string

	for _, item := range result {
		for _, k := range item.Valute {
			sortRecords(item.Valute, "Name", 1)
			tmpInternalResult = append(tmpInternalResult, item.Date, k.Name, k.ID, k.NumCode, k.CharCode, k.Nominal, k.Value)
			tmpFinalResult = append(tmpFinalResult, tmpInternalResult)
			tmpInternalResult = nil
		}
	}

	fmt.Println(tmpFinalResult)

	js, err := json.MarshalIndent(tmpFinalResult, " ", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)
}

func sortRecords(records []models.Valute, orderField string, orderBy int) {
	sort.Slice(records, func(i, j int) bool {
		if orderField == "Id" {
			if orderBy == -1 {
				return records[i].ID > records[j].ID
			} else {
				return records[i].ID < records[j].ID
			}
		} else if orderField == "NumCode" {
			if orderBy == -1 {
				return records[i].NumCode > records[j].NumCode
			} else {
				return records[i].NumCode < records[j].NumCode
			}
		} else if orderField == "Name" {
			if orderBy == -1 {
				return records[i].Name > records[j].Name
			} else {
				return records[i].Name < records[j].Name
			}
		}
		return records[i].ID > records[j].ID
	})
}
