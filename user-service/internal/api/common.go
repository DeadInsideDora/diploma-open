package api

import (
	"fmt"
	"log"
	"net/http"
	"service/internal/domain"
	"strconv"
	"strings"
	"time"
)

func AuthUserFromCookie(s domain.IDatabaseService, r *http.Request) (*domain.User, error) {
	c, err := r.Cookie("sessionid")
	if err != nil {
		return nil, err
	}
	log.Printf("cookie: %s", c.Value)
	splited := strings.Split(c.Value, "|")
	if len(splited) != 2 {
		return nil, fmt.Errorf("uncorrect cookie")
	}
	expiredTimestamp, err := strconv.ParseInt(splited[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("uncorrect cookie")
	}
	currentTimestamp := getCurrentTimestamp().Unix()
	if expiredTimestamp <= currentTimestamp {
		return nil, fmt.Errorf("expired cookie")
	}
	id, err := strconv.ParseInt(splited[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("uncorrect cookie")
	}
	return s.GetUserById(id)
}

func GetCookie(user *domain.User) *http.Cookie {
	expireTimestamp := getExpiredTimestamp()
	cookieValue := fmt.Sprintf("%d|%d", user.ID, expireTimestamp.Unix())

	return &http.Cookie{
		Name:     "sessionid",
		Value:    cookieValue,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Expires:  expireTimestamp,
		MaxAge:   int(time.Until(expireTimestamp).Seconds()),
	}
}

func getExpiredTimestamp() time.Time {
	return getCurrentTimestamp().Add(time.Hour * 1)
}

func getCurrentTimestamp() time.Time {
	loc := time.FixedZone("UTC+3", 3*60*60)
	return time.Now().In(loc)
}
