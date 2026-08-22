package probe

import (
	"dt2026/app/notify/qqbot"
	"dt2026/httpx"
	"dt2026/lib/optional"
	"errors"
	"net/http"
)

type ProbeQQBotResponse struct {
	Ok     bool `json:"ok"`
	Status int  `json:"status"`
}

func (r ProbeQQBotResponse) GetStatus() int {
	return r.Status
}

func ProbeQQBot(message optional.Optional[string], sender *qqbot.Sender) httpx.Response {
	err := sender.SendMessage(message.Or("Probe test message"))
	if err != nil {
		if errors.Is(err, qqbot.EmptyMessageErr) {
			return httpx.NewBadRequestError("failed to send message via QQ Bot: " + err.Error())
		} else {
			return httpx.NewInternalServerError("failed to send message via QQ Bot: " + err.Error())
		}
	}

	return ProbeQQBotResponse{
		Ok:     true,
		Status: http.StatusOK,
	}
}
