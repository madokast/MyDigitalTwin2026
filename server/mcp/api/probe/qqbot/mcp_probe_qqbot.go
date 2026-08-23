package qqbot

import (
	"dt2026/api/probe"
	"dt2026/httpx"
	"dt2026/lib/optional"
	"dt2026/server/http"
)

const Name = "probe_qqbot"

const Description = "Probe whether the QQ bot can send a message. Optional message: omit to send the default probe text; an empty or whitespace-only string is rejected. Do not use this to query records or write data."

type Input struct {
	Message *string `json:"message,omitempty" jsonschema:"Optional probe message. Omit to send the default text. An empty or whitespace-only string is rejected."`
}

type Output = probe.ProbeQQBotResponse

func ProbeQQBot(s *http.Server, in Input) (Output, *httpx.Error) {
	bot, err := s.QQBot()
	if err != nil {
		var zero Output
		return zero, err
	}

	var message optional.Optional[string]
	if in.Message != nil {
		message = optional.Some(*in.Message)
	}

	res := probe.ProbeQQBot(message, bot)
	if httpErr, ok := res.(*httpx.Error); ok {
		var zero Output
		return zero, httpErr
	}
	return res.(probe.ProbeQQBotResponse), nil
}
