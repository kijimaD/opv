package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

//go:embed index.html
var content embed.FS

type PomodoroData struct {
	TaskTitle      string `json:"taskTitle"`
	RemainingTime  int    `json:"remainingTime"`
	TotalTime      int    `json:"totalTime"`
	IsActive       bool   `json:"isActive"`
	TodayPoints    int    `json:"todayPoints"`
}

func main() {
	http.HandleFunc("/", serveHTML)
	http.HandleFunc("/api/pomodoro", pomodoroAPI)
	http.HandleFunc("/api/debug", debugAPI)

	fmt.Println("Starting server on :8007")
	log.Fatal(http.ListenAndServe(":8007", nil))
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	data, err := content.ReadFile("index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

var execCommand = exec.Command

func getEmacsValue(expr string) (string, error) {
	cmd := execCommand("emacsclient", "-e", expr)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	
	result := strings.TrimSpace(string(output))
	// Remove surrounding quotes if present
	if len(result) >= 2 && result[0] == '"' && result[len(result)-1] == '"' {
		result = result[1 : len(result)-1]
	}
	return result, nil
}

func pomodoroAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	data := PomodoroData{}

	// Get remaining time in seconds
	remainingStr, err := getEmacsValue("(org-pomodoro-remaining-seconds)")
	if err == nil && remainingStr != "nil" {
		// Try parsing as float first (Emacs might return decimal)
		if remaining, err := strconv.ParseFloat(remainingStr, 64); err == nil {
			data.RemainingTime = int(remaining)
		}
	}

	// Get total session length in minutes and convert to seconds
	lengthStr, err := getEmacsValue("org-pomodoro-length")
	if err == nil && lengthStr != "nil" {
		// Try parsing as float first
		if length, err := strconv.ParseFloat(lengthStr, 64); err == nil {
			data.TotalTime = int(length) * 60 // Convert minutes to seconds
		}
	}

	// Check if pomodoro is active
	activeStr, err := getEmacsValue("(org-pomodoro-active-p)")
	if err == nil {
		data.IsActive = activeStr == "t"
	}

	// Get task title
	taskTitle, err := getEmacsValue("org-clock-heading")
	if err == nil && taskTitle != "nil" {
		data.TaskTitle = taskTitle
	}

	// Get today's points
	pointsStr, err := getEmacsValue("(kd/pmd-today-point-display)")
	if err == nil && pointsStr != "nil" {
		// Try parsing as float first
		if points, err := strconv.ParseFloat(pointsStr, 64); err == nil {
			data.TodayPoints = int(points)
		}
	}

	json.NewEncoder(w).Encode(data)
}

func debugAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	debug := make(map[string]interface{})

	// Get raw values from Emacs
	remainingRaw, _ := getEmacsValue("(org-pomodoro-remaining-seconds)")
	lengthRaw, _ := getEmacsValue("org-pomodoro-length")
	activeRaw, _ := getEmacsValue("(org-pomodoro-active-p)")
	headingRaw, _ := getEmacsValue("org-clock-heading")
	pointsRaw, _ := getEmacsValue("(kd/pmd-today-point-display)")
	pomoTimeRaw, _ := getEmacsValue("(kd/org-pomodoro-time)")

	debug["raw"] = map[string]string{
		"remaining": remainingRaw,
		"length":    lengthRaw,
		"active":    activeRaw,
		"heading":   headingRaw,
		"points":    pointsRaw,
		"pomoTime":  pomoTimeRaw,
	}

	// Parse values
	data := PomodoroData{}
	if remainingRaw != "nil" {
		if remaining, err := strconv.ParseFloat(remainingRaw, 64); err == nil {
			data.RemainingTime = int(remaining)
		}
	}
	if lengthRaw != "nil" {
		if length, err := strconv.ParseFloat(lengthRaw, 64); err == nil {
			data.TotalTime = int(length) * 60
		}
	}
	data.IsActive = activeRaw == "t"
	if headingRaw != "nil" {
		data.TaskTitle = headingRaw
	}
	if pointsRaw != "nil" {
		if points, err := strconv.ParseFloat(pointsRaw, 64); err == nil {
			data.TodayPoints = int(points)
		}
	}

	debug["parsed"] = data
	json.NewEncoder(w).Encode(debug)
}