package handlers

import (
	"encoding/json"
	"final-project/internal/models"
	"final-project/internal/usecase"
	"log"
	"net/http"
	"sort"

	"github.com/google/uuid"

	"github.com/go-chi/chi/v5"
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

func (h *Handler) PullRate(w http.ResponseWriter, r *http.Request) {
	dd := chi.URLParam(r, "dd")
	mm := chi.URLParam(r, "mm")
	yyyy := chi.URLParam(r, "yyyy")
	log.Printf("%s %s %s", dd, mm, yyyy)
	date := dd + "/" + mm + "/" + yyyy
	h.uc.PullDataByPeriod(date)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var u models.User
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.uc.AddUser(ctx, &u)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	var u models.User
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&u)
	log.Printf("user: %#v", u)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = h.uc.VerifyUser(ctx, &u)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	u.ID, err = uuid.NewRandom()
	if err != nil {
		log.Println("cannot create uuid")
	}

	ts, err := h.uc.CreateToken(u.ID.Time())
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	saveErr := h.uc.CreateAuth(ctx, int64(u.ID.Time()), ts)
	if saveErr != nil {
		w.Write([]byte(saveErr.Error()))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}

	resp, _ := json.Marshal(tokens)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
