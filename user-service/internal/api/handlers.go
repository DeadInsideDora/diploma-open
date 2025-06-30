package api

import (
	"encoding/json"
	"log"
	"net/http"
	"service/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

func AuthMeHandler(s domain.IDatabaseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := AuthUserFromCookie(s, r)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		writeUserData(w, user)
	}
}

func RegisterHandler(s domain.IDatabaseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name     string `json:"name"`
			Login    string `json:"login"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}
		if len(req.Name) == 0 {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "hash error", http.StatusInternalServerError)
			return
		}

		err = s.CreateUser(req.Login, req.Name, string(hash))
		if err != nil {
			if err.Error() == "user with login already exists" {
				http.Error(w, "user exists", http.StatusConflict)
			} else {
				log.Println(err)
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func LoginHandler(s domain.IDatabaseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		user, err := s.GetUserByLogin(req.Login)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, GetCookie(user))

		writeUserData(w, user)
	}
}

func GetDataHandler(s domain.IDatabaseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := AuthUserFromCookie(s, r)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		writeUserData(w, user)
	}
}

func UpdateCardsHandler(s domain.IDatabaseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := AuthUserFromCookie(s, r)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var cards []string
		if err := json.NewDecoder(r.Body).Decode(&cards); err != nil {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		if err := s.UpdateCards(user.ID, cards); err != nil {
			if err.Error() == "can't marshal cards" {
				http.Error(w, "invalid input", http.StatusBadRequest)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		user.Cards = cards
		writeUserData(w, user)
	}
}

func UpdateMapInfoHandler(s domain.IDatabaseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := AuthUserFromCookie(s, r)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var req domain.MapInfo
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		if err := s.UpdateMapInfo(user.ID, req); err != nil {
			if err.Error() == "can't marshal map info" {
				http.Error(w, "invalid input", http.StatusBadRequest)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		user.Info = &req
		writeUserData(w, user)
	}
}

func UpdateExchangeHandler(s domain.IDatabaseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := AuthUserFromCookie(s, r)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var req struct {
			Exchange int `json:"exchange"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		if err := s.UpdateExchange(user.ID, req.Exchange); err != nil {
			if err.Error() == "can't marshal cards" {
				http.Error(w, "invalid input", http.StatusBadRequest)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		user.Exchange = int64(req.Exchange)
		writeUserData(w, user)
	}
}

func writeUserData(w http.ResponseWriter, user *domain.User) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(struct {
		Cards    []string        `json:"cards"`
		Info     *domain.MapInfo `json:"map_info"`
		Name     string          `json:"name"`
		Login    string          `json:"login"`
		Exchange int64           `json:"exchange"`
	}{Cards: user.Cards, Info: user.Info, Name: user.Name, Login: user.Login, Exchange: user.Exchange}); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
