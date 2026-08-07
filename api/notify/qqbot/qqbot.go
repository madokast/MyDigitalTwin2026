package qqbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const TokenURL = "https://bots.qq.com/app/getAppAccessToken"

var defaultAPIBases = []string{
	"https://api.sgroup.qq.com",
	"https://api.bot.qq.com",
}

type Sender struct {
	appID       string
	appSecret   string
	userOpenID  string
	token       string
	tokenExpiry time.Time
	apiBases    []string
	mu          sync.RWMutex
	client      *http.Client
}

func NewSender(appID, appSecret, userOpenID string) *Sender {
	return &Sender{
		appID:      appID,
		appSecret:  appSecret,
		userOpenID: userOpenID,
		apiBases:   defaultAPIBases,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (s *Sender) getToken(forceRefresh bool) (string, error) {
	// 非强制刷新时，先检查缓存的 Token 是否有效
	s.mu.RLock()
	if !forceRefresh && s.token != "" && time.Now().Before(s.tokenExpiry) {
		token := s.token
		s.mu.RUnlock()
		return token, nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	// 双重检查（防止并发时重复刷新）
	if !forceRefresh && s.token != "" && time.Now().Before(s.tokenExpiry) {
		return s.token, nil
	}

	// 保存旧 Token 及过期时间，用于刷新失败时回退
	oldToken := s.token
	oldExpiry := s.tokenExpiry

	// 构造刷新请求
	payload := map[string]string{
		"appId":        s.appID,
		"clientSecret": s.appSecret,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		// 刷新失败，若旧 Token 仍有效则返回旧 Token
		if oldToken != "" && time.Now().Before(oldExpiry) {
			return oldToken, nil
		}
		return "", fmt.Errorf("failed to marshal token request: %w", err)
	}

	resp, err := s.client.Post(TokenURL, "application/json", bytes.NewReader(data))
	if err != nil {
		if oldToken != "" && time.Now().Before(oldExpiry) {
			return oldToken, nil
		}
		return "", fmt.Errorf("failed to fetch token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		if oldToken != "" && time.Now().Before(oldExpiry) {
			return oldToken, nil
		}
		return "", fmt.Errorf("failed to read token response: %w", err)
	}

	var parsed struct {
		AccessToken string      `json:"access_token"`
		ExpiresIn   json.Number `json:"expires_in"`
		Message     string      `json:"message"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		if oldToken != "" && time.Now().Before(oldExpiry) {
			return oldToken, nil
		}
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	if parsed.AccessToken == "" {
		if oldToken != "" && time.Now().Before(oldExpiry) {
			return oldToken, nil
		}
		msg := parsed.Message
		if msg == "" {
			msg = "unknown error"
		}
		return "", fmt.Errorf("token fetch failed: %s (body: %s)", msg, string(body))
	}

	// 修正过期时间计算：max(1, expiresIn - 60)
	expiresIn, err := parsed.ExpiresIn.Int64()
	if err != nil {
		expiresIn = 7200 // 默认值
	}
	if expiresIn <= 0 {
		expiresIn = 7200 // 默认值
	}
	effective := expiresIn - 60
	if effective < 1 {
		effective = 1
	}
	s.token = parsed.AccessToken
	s.tokenExpiry = time.Now().Add(time.Duration(effective) * time.Second)

	return s.token, nil
}

func (s *Sender) SendMessage(text string) error {
	path := "/v2/users/" + url.PathEscape(s.userOpenID) + "/messages"
	payload := map[string]any{
		"content":  text,
		"msg_type": 0,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	var lastErr error
	// 最多尝试两次：第一次正常，第二次强制刷新 Token
	for attempt := 0; attempt < 2; attempt++ {
		// 获取 Token，attempt > 0 时强制刷新
		token, err := s.getToken(attempt > 0)
		if err != nil {
			// 如果第一次获取失败，继续尝试第二次（强制刷新）
			if attempt == 0 {
				lastErr = fmt.Errorf("failed to get token (will retry): %w", err)
				continue
			}
			// 第二次仍然失败，返回错误
			return fmt.Errorf("failed to get token after retry: %w", err)
		}

		// 重置 unauthorized 标志，用于判断本次尝试是否因 401 失败
		var isUnauthorized bool
		for _, base := range s.apiBases {
			req, err := http.NewRequest("POST", base+path, bytes.NewReader(data))
			if err != nil {
				lastErr = fmt.Errorf("request creation failed: %w", err)
				continue
			}
			req.Header.Set("Authorization", "QQBot "+token)
			req.Header.Set("Content-Type", "application/json")

			resp, err := s.client.Do(req)
			if err != nil {
				lastErr = fmt.Errorf("request failed: %w", err)
				continue
			}
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				lastErr = fmt.Errorf("failed to read response: %w", err)
				continue
			}

			// 成功
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return nil
			}

			// 处理 401 未授权：标记并跳出内部循环
			if resp.StatusCode == http.StatusUnauthorized {
				isUnauthorized = true
				var errResp struct {
					Message string `json:"message"`
					Msg     string `json:"msg"`
					Code    int    `json:"code"`
				}
				if err := json.Unmarshal(body, &errResp); err == nil && errResp.Message != "" {
					lastErr = fmt.Errorf("API error (code %d): %s", errResp.Code, errResp.Message)
				} else {
					lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
				}
				break // 401 无需尝试其他 base，直接跳出内部循环
			}

			// 处理其他非 401 错误
			var errResp struct {
				Message string `json:"message"`
				Msg     string `json:"msg"`
				Code    int    `json:"code"`
			}
			if err := json.Unmarshal(body, &errResp); err == nil {
				if errResp.Message != "" {
					lastErr = fmt.Errorf("API error (code %d): %s", errResp.Code, errResp.Message)
				} else if errResp.Msg != "" {
					lastErr = fmt.Errorf("API error (code %d): %s", errResp.Code, errResp.Msg)
				} else {
					lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
				}
			} else {
				lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
			}
			// 非 401 错误继续尝试下一个 base
		}

		// 如果本次尝试是因为 401 且是第一次尝试，则继续外层循环（强制刷新后重试）
		if isUnauthorized && attempt == 0 {
			lastErr = fmt.Errorf("unauthorized, will refresh token and retry: %w", lastErr)
			continue
		}

		// 否则（非 401 错误，或已经是第二次尝试），跳出外层循环，最终返回错误
		break
	}

	return fmt.Errorf("QQ Bot sendMessage failed after all attempts: %w", lastErr)
}
