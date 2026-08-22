package server

import (
	"context"
	"dt2026/app/envkeys"
	"dt2026/app/notify/qqbot"
	"dt2026/app/records"
	"dt2026/env"
	"dt2026/httpx"
	"encoding/json/v2"
	"fmt"
	"net/http"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	pool     func() (*pgxpool.Pool, *httpx.Error)
	qqbot    func() (*qqbot.Sender, *httpx.Error)
	testMode bool
}

func NewServer() *Server {
	tm, _ := env.Get(envkeys.TestMode)
	testMode := tm == "true"

	return &Server{
		pool: sync.OnceValues(func() (*pgxpool.Pool, *httpx.Error) {
			url, ok := env.Get(envkeys.DatabaseUrl)
			if !ok {
				return nil, httpx.NewInternalServerError("database url not set")
			}
			pool, err := pgxpool.New(context.Background(), url)
			if err != nil {
				return nil, httpx.NewInternalServerError("failed to connect to database: " + err.Error())
			}
			_, err = pool.Exec(context.Background(), records.CreateRecordSQL)
			if err != nil {
				pool.Close()
				return nil, httpx.NewInternalServerError("failed to create records table: " + err.Error())
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
			return sender, nil
		}),
		testMode: testMode,
	}
}

func (s *Server) DB() (*pgxpool.Pool, *httpx.Error) {
	return s.pool()
}

func (s *Server) QQBot() (*qqbot.Sender, *httpx.Error) {
	return s.qqbot()
}

func (s *Server) JsonUnmarshal(r *http.Request, v any) *httpx.Error {
	err := json.UnmarshalRead(r.Body, v, json.RejectUnknownMembers(true))
	if err != nil {
		return httpx.NewBadRequestError(fmt.Sprintf("failed to unmarshal request body: %s", err.Error()))
	}
	return nil
}
