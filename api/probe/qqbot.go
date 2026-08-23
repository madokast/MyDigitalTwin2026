package probe

import (
	"dt2026/api/notify/qqbot"
	"dt2026/httpx"
	"dt2026/lib/optional"
	"errors"
	"net/http"
	"strings"
)

type ProbeQQBotResponse struct {
	Ok      bool   `json:"ok" jsonschema:"Whether the request succeeded"`
	Status  int    `json:"status" jsonschema:"HTTP status code"`
	Message string `json:"message" jsonschema:"The message that was actually sent"`
}

func (r ProbeQQBotResponse) GetStatus() int {
	return r.Status
}

func ProbeQQBot(message optional.Optional[string], sender *qqbot.Sender) httpx.Response {
	sent := strings.TrimSpace(message.Or("Probe test message"))
	err := sender.SendMessage(sent)
	if err != nil {
		if errors.Is(err, qqbot.EmptyMessageErr) {
			return httpx.NewBadRequestError("failed to send message via QQ Bot: " + err.Error())
		} else {
			return httpx.NewInternalServerError("failed to send message via QQ Bot: " + err.Error())
		}
	}

	return ProbeQQBotResponse{
		Ok:      true,
		Status:  http.StatusOK,
		Message: sent,
	}
}
