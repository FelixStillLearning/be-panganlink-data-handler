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
	GetForecast(komoditasID string, periods int) (interface{}, error)
	UpdateData(komoditasID string, tanggal string, hargaAktual float64) (interface{}, error)
	DeleteData(komoditasID string, tanggal string) (interface{}, error)
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

func (s *aiService) makeDeleteRequest(endpoint string) (interface{}, error) {
	url := fmt.Sprintf("%s%s", s.baseURL, endpoint)
	
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
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
	return s.makePostRequest("/api/v1/ai/recommend", map[string]string{"komoditas_id": komoditasID})
}

func (s *aiService) GetForecast(komoditasID string, periods int) (interface{}, error) {
	payload := map[string]interface{}{
		"komoditas_id": komoditasID,
		"periods":      periods,
	}
	return s.makePostRequest("/api/v1/ai/forecast", payload)
}

func (s *aiService) UpdateData(komoditasID string, tanggal string, hargaAktual float64) (interface{}, error) {
	payload := map[string]interface{}{
		"komoditas_id": komoditasID,
		"tanggal":      tanggal,
		"harga_aktual": hargaAktual,
	}
	return s.makePostRequest("/api/v1/ai/update_data", payload)
}

func (s *aiService) DeleteData(komoditasID string, tanggal string) (interface{}, error) {
	endpoint := fmt.Sprintf("/api/v1/ai/delete_data?komoditas_id=%s&tanggal=%s", komoditasID, tanggal)
	return s.makeDeleteRequest(endpoint)
}
