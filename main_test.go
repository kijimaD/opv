package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServeHTML(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(serveHTML)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, "text/html; charset=utf-8")
	}

	body := rr.Body.String()
	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Errorf("handler returned unexpected body: missing DOCTYPE")
	}
	if !strings.Contains(body, "Pomodoro Timer") {
		t.Errorf("handler returned unexpected body: missing title")
	}
}

func TestPomodoroAPI(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/pomodoro", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(pomodoroAPI)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, "application/json")
	}

	cors := rr.Header().Get("Access-Control-Allow-Origin")
	if cors != "*" {
		t.Errorf("handler returned wrong CORS header: got %v want %v",
			cors, "*")
	}

	var data PomodoroData
	err = json.Unmarshal(rr.Body.Bytes(), &data)
	if err != nil {
		t.Errorf("handler returned invalid JSON: %v", err)
	}

	if data.TaskTitle == "" {
		t.Errorf("handler returned empty task title")
	}
	if data.RemainingTime < 0 {
		t.Errorf("handler returned negative remaining time: %v", data.RemainingTime)
	}
	if data.TotalTime <= 0 {
		t.Errorf("handler returned invalid total time: %v", data.TotalTime)
	}
	if data.RemainingTime > data.TotalTime {
		t.Errorf("remaining time (%v) is greater than total time (%v)", 
			data.RemainingTime, data.TotalTime)
	}
}

func TestPomodoroDataJSON(t *testing.T) {
	testData := PomodoroData{
		TaskTitle:     "Test Task",
		RemainingTime: 300,
		TotalTime:     1500,
		IsActive:      true,
		TodayPoints:   10,
	}

	jsonData, err := json.Marshal(testData)
	if err != nil {
		t.Fatalf("Failed to marshal PomodoroData: %v", err)
	}

	var decoded PomodoroData
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal PomodoroData: %v", err)
	}

	if decoded.TaskTitle != testData.TaskTitle {
		t.Errorf("TaskTitle mismatch: got %v want %v", decoded.TaskTitle, testData.TaskTitle)
	}
	if decoded.RemainingTime != testData.RemainingTime {
		t.Errorf("RemainingTime mismatch: got %v want %v", decoded.RemainingTime, testData.RemainingTime)
	}
	if decoded.TotalTime != testData.TotalTime {
		t.Errorf("TotalTime mismatch: got %v want %v", decoded.TotalTime, testData.TotalTime)
	}
	if decoded.IsActive != testData.IsActive {
		t.Errorf("IsActive mismatch: got %v want %v", decoded.IsActive, testData.IsActive)
	}
	if decoded.TodayPoints != testData.TodayPoints {
		t.Errorf("TodayPoints mismatch: got %v want %v", decoded.TodayPoints, testData.TodayPoints)
	}
}