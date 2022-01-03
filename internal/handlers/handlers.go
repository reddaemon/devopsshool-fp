package handlers

import (
	"encoding/json"
	"final-project/internal/middleware"
	"final-project/internal/models"
	"final-project/internal/usecase"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
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
	log.Printf("user: %#v", u.Email)
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

	u.ID = float64(uuid.ClockSequence())
	if err != nil {
		log.Println("cannot create uuid")
	}

	ts, err := h.uc.CreateToken(u.ID)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	saveErr := h.uc.CreateAuth(ctx, u.ID, ts)
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Vary", "Origin")
	w.Header().Add("Vary", "Access-Control-Request-Method")
	w.Header().Add("Vary", "Access-Control-Request-Headers")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST,OPTIONS")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	au, err := middleware.ExtractTokenMetadata(r)

	if err != nil {
		log.Println("cannot extract token metadata...")
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte("unauthorized"))
		return
	}

	deleted, err := h.uc.DeleteAuth(ctx, au.AccessUuid)
	if err != nil { //if any goes wrong
		log.Printf("%d", deleted)
		log.Println("cannot delete auth for accessUUid")
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte("unauthorized"))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Vary", "Origin")
	w.Header().Add("Vary", "Access-Control-Request-Method")
	w.Header().Add("Vary", "Access-Control-Request-Headers")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST,OPTIONS")
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte("Successfully logged out"))
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	providedToken := middleware.ExtractToken(r)
	ctx := r.Context()
	mapToken := map[string]string{"refresh_token": providedToken}
	_, err := json.Marshal(mapToken)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(err.Error()))
		return
	}
	refreshtoken := mapToken["refresh_token"]

	//verify the token
	token, err := jwt.Parse(refreshtoken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		log.Printf("token expired error: %s", err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte("Refresh token expired"))
		return
	}

	//is token valid?
	if _, ok := token.Claims.(jwt.StandardClaims); !ok && !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) // the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(float64) // convert interface to string
		if !ok {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Header().Add("Content-Type", "application/json")
		}

		userId, err := strconv.ParseFloat(fmt.Sprintf("%.f", claims["user_id"]), 64)
		//userId := claims["user_id"]
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte("Error occurred"))
			return
		}
		refreshUuidStr := fmt.Sprintf("%f", refreshUuid)

		// Delete the previous Refresh Token
		deleted, err := h.uc.DeleteAuth(ctx, refreshUuidStr)
		if err != nil {
			log.Printf("deleted: %d", deleted)
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte("unauthorized"))
			return
		}
		// Create new pairs of refresh and access tokens
		ts, err := h.uc.CreateToken(userId)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			return
		}
		// save the tokens metadata to redis
		err = h.uc.CreateAuth(ctx, userId, ts)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte(err.Error()))
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		resp, _ := json.Marshal(tokens)

		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")
		w.Header().Add("Vary", "Access-Control-Request-Headers")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST,OPTIONS")
		w.WriteHeader(http.StatusCreated)
		w.Write(resp)

	}
}
