package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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

var authApi string = "http://20.244.56.144/evaluation-service/auth"

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

func getAccessToken() {
	reqBody, err := json.Marshal(getAccessTokenRequest())
	if err != nil {
		log.Fatalf("Error marshaling token request: %v", err)
	}

	resp, err := http.Post(authApi, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatalf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Auth API failed with status %d", resp.StatusCode)
	}

	var token ResponseAuth
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		log.Fatalf("Failed to decode token response: %v", err)
	}

	// Print the token
	fmt.Println("Access Token:", token.AccessToken)
	fmt.Println("Token Type:", token.TokenType)
	fmt.Println("Expires In:", token.ExpiresIn)
}

func main() {
	getAccessToken()
}
