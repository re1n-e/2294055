package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Request struct {
	Email          string `json:"email"`
	Name           string `json:"name"`
	MobileNo       string `json:"mobileNo"`
	GithubUsername string `json:"githubUsername"`
	RollNo         string `json:"rollNo"`
	AccessCode     string `json:"accessCode"`
}

type Registered struct {
	Email        string `json:"email"`
	Name         string `json:"name"`
	RollNo       string `json:"rollNo"`
	AccessCode   string `json:"accessCode"`
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
}

var regApi string = "http://20.244.56.144/evaluation-service/register"

func createRequestHeader() *Request {
	return &Request{
		Email:          "raghavendrabargali@gmail.com",
		Name:           "Raghavendra Singh Bargali",
		MobileNo:       "7505137114",
		GithubUsername: "re1n-e",
		RollNo:         "2294055",
		AccessCode:     "QAhDUr",
	}
}

func get_registered_data() {
	reqBody, err := json.Marshal(createRequestHeader())
	if err != nil {
		log.Fatalf("Failed to marshal request: %v", err)
	}

	resp, err := http.Post(regApi, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatalf("Failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var reg Registered
	if err := json.NewDecoder(resp.Body).Decode(&reg); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	// Print the important info
	fmt.Println("ClientID:", reg.ClientID)
	fmt.Println("ClientSecret:", reg.ClientSecret)
}

func main() {
	get_registered_data()
}
