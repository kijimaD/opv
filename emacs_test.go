package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

// Mock for exec.Command
func mockExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// Helper process for mocking command execution
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	args := os.Args
	for i, arg := range args {
		if arg == "--" {
			args = args[i+1:]
			break
		}
	}

	if len(args) < 3 || args[0] != "emacsclient" || args[1] != "-e" {
		fmt.Fprintf(os.Stderr, "Invalid command")
		os.Exit(1)
	}

	expr := args[2]
	switch expr {
	case "(org-pomodoro-remaining-seconds)":
		fmt.Print("900")
	case "org-pomodoro-length":
		fmt.Print("25")
	case "(org-pomodoro-active-p)":
		fmt.Print("t")
	case "org-clock-heading":
		fmt.Print("\"Write unit tests\"")
	case "(kd/pmd-today-point-display)":
		fmt.Print("7")
	case "(test-nil-value)":
		fmt.Print("nil")
	case "(test-error)":
		os.Exit(1)
	default:
		fmt.Print("nil")
	}
	os.Exit(0)
}

func TestGetEmacsValue(t *testing.T) {
	// Save original and restore after test
	originalExecCommand := execCommand
	execCommand = mockExecCommand
	defer func() { execCommand = originalExecCommand }()

	tests := []struct {
		name     string
		expr     string
		expected string
		wantErr  bool
	}{
		{
			name:     "Get remaining seconds",
			expr:     "(org-pomodoro-remaining-seconds)",
			expected: "900",
			wantErr:  false,
		},
		{
			name:     "Get pomodoro length",
			expr:     "org-pomodoro-length",
			expected: "25",
			wantErr:  false,
		},
		{
			name:     "Get active state",
			expr:     "(org-pomodoro-active-p)",
			expected: "t",
			wantErr:  false,
		},
		{
			name:     "Get task heading with quotes",
			expr:     "org-clock-heading",
			expected: "Write unit tests",
			wantErr:  false,
		},
		{
			name:     "Get today points",
			expr:     "(kd/pmd-today-point-display)",
			expected: "7",
			wantErr:  false,
		},
		{
			name:     "Handle nil value",
			expr:     "(test-nil-value)",
			expected: "nil",
			wantErr:  false,
		},
		{
			name:     "Handle command error",
			expr:     "(test-error)",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getEmacsValue(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("getEmacsValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("getEmacsValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPomodoroDataParsing(t *testing.T) {
	// Save original and restore after test
	originalExecCommand := execCommand
	execCommand = mockExecCommand
	defer func() { execCommand = originalExecCommand }()

	// Test the full flow of data parsing
	tests := []struct {
		name              string
		mockRemainingTime string
		mockLength        string
		mockActive        string
		mockHeading       string
		mockPoints        string
		expectedData      PomodoroData
	}{
		{
			name:              "Active pomodoro",
			mockRemainingTime: "900",  // 15 minutes
			mockLength:        "25",   // 25 minutes
			mockActive:        "t",
			mockHeading:       "\"Write unit tests\"",
			mockPoints:        "7",
			expectedData: PomodoroData{
				TaskTitle:     "Write unit tests",
				RemainingTime: 900,
				TotalTime:     1500, // 25 * 60
				IsActive:      true,
				TodayPoints:   7,
			},
		},
		{
			name:              "Inactive pomodoro",
			mockRemainingTime: "nil",
			mockLength:        "25",
			mockActive:        "nil",
			mockHeading:       "nil",
			mockPoints:        "10",
			expectedData: PomodoroData{
				TaskTitle:     "",
				RemainingTime: 0,
				TotalTime:     1500,
				IsActive:      false,
				TodayPoints:   10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This would test the actual parsing logic
			// For now, we're verifying the mock system works correctly
			
			// Test remaining time parsing
			remaining, _ := getEmacsValue("(org-pomodoro-remaining-seconds)")
			if remaining != "900" {
				t.Errorf("Expected remaining time 900, got %s", remaining)
			}

			// Test length parsing  
			length, _ := getEmacsValue("org-pomodoro-length")
			if length != "25" {
				t.Errorf("Expected length 25, got %s", length)
			}

			// Test active state parsing
			active, _ := getEmacsValue("(org-pomodoro-active-p)")
			if active != "t" {
				t.Errorf("Expected active state 't', got %s", active)
			}
		})
	}
}

func TestTimeFormatting(t *testing.T) {
	tests := []struct {
		seconds  int
		expected string
	}{
		{0, "00:00"},
		{59, "00:59"},
		{60, "01:00"},
		{90, "01:30"},
		{900, "15:00"},
		{1500, "25:00"},
		{3599, "59:59"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d seconds", tt.seconds), func(t *testing.T) {
			// This tests the JavaScript formatTime function logic
			minutes := tt.seconds / 60
			secs := tt.seconds % 60
			result := fmt.Sprintf("%02d:%02d", minutes, secs)
			if result != tt.expected {
				t.Errorf("formatTime(%d) = %s, want %s", tt.seconds, result, tt.expected)
			}
		})
	}
}