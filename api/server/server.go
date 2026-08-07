package server

import (
	"context"
	"dt2026/api/envkeys"
	"dt2026/api/notify/qqbot"
	"dt2026/env"
	"dt2026/httpx"
	"net/http"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	db    func() (*pgxpool.Pool, *httpx.Error)
	qqbot func() (*qqbot.Sender, *httpx.Error)
}

func NewServer() *Server {
	return &Server{
		db: sync.OnceValues(func() (*pgxpool.Pool, *httpx.Error) {
			url, ok := env.Get(envkeys.DatabaseUrl)
			if !ok {
				return nil, httpx.NewInternalServerError("database url not set")
			}
			pool, err := pgxpool.New(context.Background(), url)
			if err != nil {
				return nil, httpx.NewInternalServerError("failed to connect to database: " + err.Error())
			}
			return pool, nil
		}),
		qqbot: sync.OnceValues(func() (*qqbot.Sender, *httpx.Error) {
			appID, ok := env.Get(envkeys.QQBotAppID)
			if !ok {
				return nil, httpx.NewInternalServerError("QQ Bot App ID not set")
			}

			appSecret, ok := env.Get(envkeys.QQBotAppSecret)
			if !ok {
				return nil, httpx.NewInternalServerError("QQ Bot App Secret not set")
			}

			userOpenID, ok := env.Get(envkeys.QQBotUserOpenID)
			if !ok {
				return nil, httpx.NewInternalServerError("QQ Bot User OpenID not set")
			}

			sender := qqbot.NewSender(appID, appSecret, userOpenID)
			err := sender.SendMessage("Probe test message")
			if err != nil {
				return nil, httpx.NewInternalServerError("failed to send message via QQ Bot: " + err.Error())
			}
			return sender, nil
		}),
	}
}

func (s *Server) DB() (*pgxpool.Pool, *httpx.Error) {
	return s.db()
}

func (s *Server) QQBot() (*qqbot.Sender, *httpx.Error) {
	return s.qqbot()
}

func Handle(w http.ResponseWriter, r *http.Request) {

}
