package probe

import (
	"dt2026/api/envkeys"
	"dt2026/api/httpx"
	"dt2026/api/notify/qqbot"
	"dt2026/env"
	"net/http"
)

type ProbeQQBotResponse struct {
	Ok     bool `json:"ok"`
	Status int  `json:"status"`
}

func (r ProbeQQBotResponse) IsOk() bool {
	return r.Ok
}

func (r ProbeQQBotResponse) GetStatus() int {
	return r.Status
}

func ProbeQQBot(w http.ResponseWriter, r *http.Request) {

	appID, ok := env.Get(envkeys.QQBotAppID)
	if !ok {
		httpx.WriteJSON(w, httpx.NewInternalServerError("QQ Bot App ID not set"))
		return
	}

	appSecret, ok := env.Get(envkeys.QQBotAppSecret)
	if !ok {
		httpx.WriteJSON(w, httpx.NewInternalServerError("QQ Bot App Secret not set"))
		return
	}

	userOpenID, ok := env.Get(envkeys.QQBotUserOpenID)
	if !ok {
		httpx.WriteJSON(w, httpx.NewInternalServerError("QQ Bot User OpenID not set"))
		return
	}

	sender := qqbot.NewSender(appID, appSecret, userOpenID)
	err := sender.SendMessage("Probe test message")
	if err != nil {
		httpx.WriteJSON(w, httpx.NewInternalServerError("failed to send message via QQ Bot: "+err.Error()))
		return
	}

	response := ProbeQQBotResponse{
		Ok:     true,
		Status: http.StatusOK,
	}
	httpx.WriteJSON(w, response)
}
