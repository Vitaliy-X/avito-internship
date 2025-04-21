package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const baseURL = "http://localhost:8080"

func TestIntegrationHTTPFlow(t *testing.T) {
	token := getToken(t, "moderator")

	pvzID := "11111111-1111-1111-1111-111111111111"
	createPVZ(t, token, pvzID)

	employeeToken := getToken(t, "employee")

	createReception(t, employeeToken, pvzID)

	for i := 0; i < 50; i++ {
		addProduct(t, employeeToken, pvzID, "одежда")
	}

	closeReception(t, employeeToken, pvzID)
}

func getToken(t *testing.T, role string) string {
	body := []byte(`{"role":"` + role + `"}`)
	resp, err := http.Post(baseURL+"/dummyLogin", "application/json", bytes.NewReader(body))
	assert.NoError(t, err)
	defer resp.Body.Close()

	data := struct {
		Token string `json:"token"`
	}{}
	decodeJSON(t, resp.Body, &data)
	return data.Token
}

func createPVZ(t *testing.T, token, pvzID string) {
	payload := map[string]interface{}{
		"id":               pvzID,
		"registrationDate": "2025-04-21T15:00:00Z",
		"city":             "Москва",
	}
	postJSON(t, "/pvz", token, payload)
}

func createReception(t *testing.T, token, pvzID string) {
	payload := map[string]string{
		"pvzId": pvzID,
	}
	postJSON(t, "/receptions", token, payload)
}

func addProduct(t *testing.T, token, pvzID, typ string) {
	payload := map[string]string{
		"pvzId": pvzID,
		"type":  typ,
	}
	postJSON(t, "/products", token, payload)
}

func closeReception(t *testing.T, token, pvzID string) {
	url := fmt.Sprintf("%s/pvz/%s/close_last_reception", baseURL, pvzID)
	req, _ := http.NewRequest(http.MethodPost, url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func postJSON(t *testing.T, path, token string, payload interface{}) {
	data, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, baseURL+path, bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated)
}

func decodeJSON(t *testing.T, r io.Reader, v interface{}) {
	dec := json.NewDecoder(r)
	err := dec.Decode(v)
	assert.NoError(t, err)
}
