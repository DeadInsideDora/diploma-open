package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"service/internal/domain"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgxService struct {
	db *pgxpool.Pool
}

func NewPgxService(url string) (*PgxService, error) {
	db, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}

	return &PgxService{db: db}, nil
}

func (service *PgxService) GetUserById(id int64) (*domain.User, error) {
	row := service.db.QueryRow(context.Background(),
		"SELECT id, login, name, password_hash, cards, map_info, exchange FROM users WHERE id = $1", id)
	return extractUser(row)
}

func (service *PgxService) GetUserByLogin(login string) (*domain.User, error) {
	row := service.db.QueryRow(context.Background(),
		"SELECT id, login, name, password_hash, cards, map_info, exchange FROM users WHERE login = $1", login)
	return extractUser(row)
}

func (service *PgxService) CreateUser(login, name, passwordHash string) error {
	log.Printf("creating user: log=%s, name=%s", login, name)
	_, err := service.db.Exec(context.Background(),
		"INSERT INTO users (login, name, password_hash) VALUES ($1, $2, $3)", login, name, passwordHash)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return fmt.Errorf("user with login already exists")
		}
		return err
	}
	return nil
}

func (service *PgxService) UpdateCards(id int, cards []string) error {
	cardsJSON, err := json.Marshal(cards)
	if err != nil {
		return fmt.Errorf("can't marshal cards")
	}
	_, err = service.db.Exec(context.Background(), "UPDATE users SET cards=$1 WHERE id=$2", cardsJSON, id)
	return err
}

func (service *PgxService) UpdateMapInfo(id int, mapInfo domain.MapInfo) error {
	mapInfoJson, err := json.Marshal(mapInfo)
	if err != nil {
		return fmt.Errorf("can't marshal map info")
	}
	_, err = service.db.Exec(context.Background(), "UPDATE users SET map_info=$1 WHERE id=$2", mapInfoJson, id)
	return err
}

func (service *PgxService) UpdateExchange(id int, exchange int) error {
	_, err := service.db.Exec(context.Background(), "UPDATE users SET exchange=$1 WHERE id=$2", exchange, id)
	return err
}

func (service *PgxService) Close() {
	service.db.Close()
}

func extractUser(row pgx.Row) (*domain.User, error) {
	var u domain.User
	var cards []byte
	var map_info []byte
	err := row.Scan(&u.ID, &u.Login, &u.Name, &u.Password, &cards, &map_info, &u.Exchange)

	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(cards, &u.Cards); err != nil {
		return nil, err
	}
	if len(map_info) != 0 {
		if err := json.Unmarshal(map_info, &u.Info); err != nil {
			return nil, err
		}
	}

	return &u, nil
}
