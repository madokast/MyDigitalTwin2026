package probe

import (
	"dt2026/api/notify/qqbot"
	"dt2026/httpx"
	"net/http"
)

type ProbeQQBotResponse struct {
	Ok     bool `json:"ok"`
	Status int  `json:"status"`
}

func (r *ProbeQQBotResponse) GetStatus() int {
	return r.Status
}

func ProbeQQBot(sender *qqbot.Sender) httpx.Response {
	err := sender.SendMessage("Probe test message")
	if err != nil {
		return httpx.NewInternalServerError("failed to send message via QQ Bot: " + err.Error())
	}

	return &ProbeQQBotResponse{
		Ok:     true,
		Status: http.StatusOK,
	}
}
