package Log

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

// ---------- Auth + Logger Setup ----------

type AccessTokenRequest struct {
	Email        string `json:"email"`
	Name         string `json:"name"`
	RollNo       string `json:"rollNo"`
	AccessCode   string `json:"accessCode"`
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
}

type ResponseAuth struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type LogRequest struct {
	Stack   string `json:"stack"`
	Level   string `json:"level"`
	Package string `json:"package"`
	Message string `json:"message"`
}

var (
	authAPI    = "http://20.244.56.144/evaluation-service/auth"
	logAPI     = "http://20.244.56.144/evaluation-service/logs"
	httpClient = &http.Client{Timeout: 5 * time.Second}
	token      *ResponseAuth
	tokenOnce  sync.Once
)

func getAccessTokenRequest() *AccessTokenRequest {
	return &AccessTokenRequest{
		Email:        "raghavendrabargali@gmail.com",
		Name:         "Raghavendra Singh Bargali",
		RollNo:       "2294055",
		AccessCode:   "QAhDUr",
		ClientID:     "262ae3a5-1697-4f07-8d88-783f22f6ff71",
		ClientSecret: "EwpXRmDMBPjTdXjX",
	}
}

func fetchToken() *ResponseAuth {
	reqBody, _ := json.Marshal(getAccessTokenRequest())
	resp, err := httpClient.Post(authAPI, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatalf("Token fetch failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Fatalf("Auth failed with status %d", resp.StatusCode)
	}

	var t ResponseAuth
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		log.Fatalf("Failed to decode token: %v", err)
	}
	return &t
}

func getToken() string {
	tokenOnce.Do(func() {
		token = fetchToken()
	})
	return token.AccessToken
}

func Log(stack, level, pkg, message string) {
	logPayload := LogRequest{Stack: stack, Level: level, Package: pkg, Message: message}
	jsonData, _ := json.Marshal(logPayload)

	req, _ := http.NewRequest("POST", logAPI, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+getToken())
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Log failed: %v", err)
		return
	}
	defer resp.Body.Close()
}
