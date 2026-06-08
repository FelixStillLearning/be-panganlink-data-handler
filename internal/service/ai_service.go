package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AIService interface {
	GetRecommendation(komoditasID string) (interface{}, error)
	GetForecast(komoditasID string) (interface{}, error)
}

type aiService struct {
	baseURL    string
	httpClient *http.Client
}

func NewAIService(baseURL string) AIService {
	return &aiService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *aiService) makePostRequest(endpoint string, payload interface{}) (interface{}, error) {
	url := fmt.Sprintf("%s%s", s.baseURL, endpoint)
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("AI Service error (status %d): %s", resp.StatusCode, string(body))
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *aiService) GetRecommendation(komoditasID string) (interface{}, error) {
	return s.makePostRequest("/ai/recommend", map[string]string{"komoditas_id": komoditasID})
}

func (s *aiService) GetForecast(komoditasID string) (interface{}, error) {
	return s.makePostRequest("/ai/forecast", map[string]string{"komoditas_id": komoditasID})
}
