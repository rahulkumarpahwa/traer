package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Test struct for requests/responses
type TestCreateUserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TestLoginReq struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password"`
}

type TestLoginResp struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}

type TestUpdateUserReq struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func testAllRoutes() {
	// Wait for server to start
	time.Sleep(1 * time.Second)

	baseURL := "http://localhost:8080"
	client := &http.Client{}

	fmt.Println("\n=== TESTING ALL BACKEND ROUTES ===")

	// Test 1: Health Check
	fmt.Println("1️⃣  Testing GET / (Health Check)")
	testHealthCheck(client, baseURL)

	// Test 2: Create User (Public endpoint)
	fmt.Println("\n2️⃣  Testing POST /users (Create User)")
	userEmail := "testuser@example.com"
	userName := "testuser"
	userPassword := "TestPass123!"

	testCreateUser(client, baseURL, TestCreateUserReq{
		Username: userName,
		Email:    userEmail,
		Password: userPassword,
	})

	// Test 3: Login to get token
	fmt.Println("\n3️⃣  Testing POST /login (Login)")
	token := testLogin(client, baseURL, TestLoginReq{
		Email:    userEmail,
		Password: userPassword,
	})

	if token == "" {
		fmt.Println("❌ Failed to get token, skipping protected routes")
		return
	}

	// Test 4-8: Protected routes with token
	fmt.Println("\n4️⃣  Testing GET /users/email (Get User by Email)")
	testGetUserByEmail(client, baseURL, userEmail, token)

	fmt.Println("\n5️⃣  Testing GET /users/username (Get User by Username)")
	testGetUserByUsername(client, baseURL, userName, token)

	fmt.Println("\n6️⃣  Testing PUT /users (Update User)")
	testUpdateUser(client, baseURL, TestUpdateUserReq{
		Username: "updateduser",
	}, token)

	fmt.Println("\n7️⃣  Testing POST /jobs/create (Create Job - Protected)")
	testJobCreate(client, baseURL, token)

	fmt.Println("\n8️⃣  Testing GET /jobs/active (Get Active Jobs)")
	testJobsActive(client, baseURL, token)

	fmt.Println("\n9️⃣  Testing GET /jobs/status (Get Job Status)")
	testJobStatus(client, baseURL, token)

	fmt.Println("\n🔟 Testing GET /instances/get (Get Instances)")
	testInstancesGet(client, baseURL, token)

	fmt.Println("\n✅ All route tests completed!\n")
}

func testHealthCheck(client *http.Client, baseURL string) {
	resp, err := client.Get(baseURL + "/")
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ Status: %d\n   Response: %s\n", resp.StatusCode, string(body))
}

func testCreateUser(client *http.Client, baseURL string, req TestCreateUserReq) {
	data, _ := json.Marshal(req)
	resp, err := client.Post(baseURL+"/users", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ Status: %d\n   Response: %s\n", resp.StatusCode, string(body))
}

func testLogin(client *http.Client, baseURL string, req TestLoginReq) string {
	data, _ := json.Marshal(req)
	resp, err := client.Post(baseURL+"/login", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ Status: %d\n   Response: %s\n", resp.StatusCode, string(body))

	// Extract token from cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "token" {
			fmt.Printf("✅ Token obtained: %s (first 50 chars)\n", cookie.Value[:min(50, len(cookie.Value))])
			return cookie.Value
		}
	}

	return ""
}

func testGetUserByEmail(client *http.Client, baseURL, email, token string) {
	req, _ := http.NewRequest("GET", baseURL+"/users/email?email="+email, nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: token})

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ Status: %d\n   Response: %s\n", resp.StatusCode, string(body)[:min(200, len(string(body)))])
}

func testGetUserByUsername(client *http.Client, baseURL, username, token string) {
	req, _ := http.NewRequest("GET", baseURL+"/users/username?username="+username, nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: token})

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ Status: %d\n   Response: %s\n", resp.StatusCode, string(body)[:min(200, len(string(body)))])
}

func testUpdateUser(client *http.Client, baseURL string, req TestUpdateUserReq, token string) {
	data, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("PUT", baseURL+"/users?id=1", bytes.NewBuffer(data))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.AddCookie(&http.Cookie{Name: "token", Value: token})

	resp, err := client.Do(httpReq)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ Status: %d\n   Response: %s\n", resp.StatusCode, string(body))
}

func testJobCreate(client *http.Client, baseURL string, token string) {
	jobReq := map[string]interface{}{
		"url":    "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		"output": "audio",
		"format": "mp3",
	}

	data, _ := json.Marshal(jobReq)
	httpReq, _ := http.NewRequest("POST", baseURL+"/jobs/create", bytes.NewBuffer(data))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.AddCookie(&http.Cookie{Name: "token", Value: token})

	resp, err := client.Do(httpReq)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ Status: %d\n   Response: %s\n", resp.StatusCode, string(body)[:min(300, len(string(body)))])
}

func testJobsActive(client *http.Client, baseURL string, token string) {
	req, _ := http.NewRequest("GET", baseURL+"/jobs/active", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: token})

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ Status: %d\n   Response: %s\n", resp.StatusCode, string(body)[:min(300, len(string(body)))])
}

func testJobStatus(client *http.Client, baseURL string, token string) {
	req, _ := http.NewRequest("GET", baseURL+"/jobs/status?job_id=1", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: token})

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ Status: %d\n   Response: %s\n", resp.StatusCode, string(body)[:min(300, len(string(body)))])
}

func testInstancesGet(client *http.Client, baseURL string, token string) {
	req, _ := http.NewRequest("GET", baseURL+"/instances/get", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: token})

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ Status: %d\n   Response: %s\n", resp.StatusCode, string(body)[:min(300, len(string(body)))])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func init() {
	// Only run tests if TEST_ROUTES environment variable is set
	if os.Getenv("TEST_ROUTES") == "1" {
		go testAllRoutes()
	}
}
