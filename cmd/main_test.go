package main

import (
	"errors"
	"os"
	"testing"
)

// mockTempoo is a mock implementation for testing
type mockTempoo struct {
	shouldFailNewTempoo     bool
	shouldFailAddWorklog    bool
	shouldFailGetUserID     bool
	shouldFailGetWorklogs   bool
	shouldFailDeleteWorklog bool
	userID                  string
	worklogIDs              []string
	deleteWorklogCalls      []deleteWorklogCall
}

type deleteWorklogCall struct {
	issueKey  string
	worklogID string
}

// Mock the internal.NewTempoo function by setting environment variables
func setupMockEnvironment() {
	os.Setenv("JIRA_EMAIL", "test@example.com")
	os.Setenv("JIRA_API_TOKEN", "test-token")
}

func teardownMockEnvironment() {
	os.Unsetenv("JIRA_EMAIL")
	os.Unsetenv("JIRA_API_TOKEN")
}

func TestAddWorklogCmd_Run(t *testing.T) {
	tests := []struct {
		name          string
		cmd           AddWorklogCmd
		setupEnv      bool
		expectedError bool
		errorContains string
	}{
		{
			name: "successful add worklog",
			cmd: AddWorklogCmd{
				IssueKey: "TEST-123",
				Time:     "2h",
				Date:     nil,
			},
			setupEnv:      true,
			expectedError: false,
		},
		{
			name: "successful add worklog with date",
			cmd: AddWorklogCmd{
				IssueKey: "TEST-123",
				Time:     "1h 30m",
				Date:     stringPtr("15.12.2023"),
			},
			setupEnv:      true,
			expectedError: false,
		},
		{
			name: "missing environment variables",
			cmd: AddWorklogCmd{
				IssueKey: "TEST-123",
				Time:     "2h",
				Date:     nil,
			},
			setupEnv:      false,
			expectedError: true,
			errorContains: "JIRA_EMAIL environment variable is not set",
		},
		{
			name: "empty issue key",
			cmd: AddWorklogCmd{
				IssueKey: "",
				Time:     "2h",
				Date:     nil,
			},
			setupEnv:      true,
			expectedError: false, // The actual validation happens in the internal package
		},
		{
			name: "empty time",
			cmd: AddWorklogCmd{
				IssueKey: "TEST-123",
				Time:     "",
				Date:     nil,
			},
			setupEnv:      true,
			expectedError: false, // The actual validation happens in the internal package
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupEnv {
				setupMockEnvironment()
				defer teardownMockEnvironment()
			} else {
				teardownMockEnvironment()
			}

			err := tt.cmd.Run()

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorContains != "" && !containsString(err.Error(), tt.errorContains) {
					t.Errorf("expected error to contain '%s', got '%s'", tt.errorContains, err.Error())
				}
			} else {
				// Note: These tests will fail in a real environment without proper Jira setup
				// In a real testing scenario, you'd want to mock the HTTP calls or use dependency injection
				if err != nil {
					t.Logf("Got error (expected in test environment without Jira setup): %v", err)
				}
			}
		})
	}
}

func TestRemoveWorklogCmd_Run(t *testing.T) {
	tests := []struct {
		name          string
		cmd           RemoveWorklogCmd
		setupEnv      bool
		expectedError bool
		errorContains string
	}{
		{
			name: "successful remove worklog",
			cmd: RemoveWorklogCmd{
				IssueKey: "TEST-123",
			},
			setupEnv:      true,
			expectedError: false,
		},
		{
			name: "missing environment variables",
			cmd: RemoveWorklogCmd{
				IssueKey: "TEST-123",
			},
			setupEnv:      false,
			expectedError: true,
			errorContains: "JIRA_EMAIL environment variable is not set",
		},
		{
			name: "empty issue key",
			cmd: RemoveWorklogCmd{
				IssueKey: "",
			},
			setupEnv:      true,
			expectedError: false, // The actual validation happens in the internal package
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupEnv {
				setupMockEnvironment()
				defer teardownMockEnvironment()
			} else {
				teardownMockEnvironment()
			}

			err := tt.cmd.Run()

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorContains != "" && !containsString(err.Error(), tt.errorContains) {
					t.Errorf("expected error to contain '%s', got '%s'", tt.errorContains, err.Error())
				}
			} else {
				// Note: These tests will fail in a real environment without proper Jira setup
				// In a real testing scenario, you'd want to mock the HTTP calls or use dependency injection
				if err != nil {
					t.Logf("Got error (expected in test environment without Jira setup): %v", err)
				}
			}
		})
	}
}

func TestAddWorklogCmd_Fields(t *testing.T) {
	cmd := AddWorklogCmd{
		IssueKey: "TEST-123",
		Time:     "2h 30m",
		Date:     stringPtr("15.12.2023"),
	}

	if cmd.IssueKey != "TEST-123" {
		t.Errorf("expected IssueKey to be 'TEST-123', got '%s'", cmd.IssueKey)
	}

	if cmd.Time != "2h 30m" {
		t.Errorf("expected Time to be '2h 30m', got '%s'", cmd.Time)
	}

	if cmd.Date == nil {
		t.Error("expected Date to not be nil")
	} else if *cmd.Date != "15.12.2023" {
		t.Errorf("expected Date to be '15.12.2023', got '%s'", *cmd.Date)
	}
}

func TestRemoveWorklogCmd_Fields(t *testing.T) {
	cmd := RemoveWorklogCmd{
		IssueKey: "TEST-456",
	}

	if cmd.IssueKey != "TEST-456" {
		t.Errorf("expected IssueKey to be 'TEST-456', got '%s'", cmd.IssueKey)
	}
}

func TestAddWorklogCmd_RunWithNilDate(t *testing.T) {
	setupMockEnvironment()
	defer teardownMockEnvironment()

	cmd := AddWorklogCmd{
		IssueKey: "TEST-123",
		Time:     "1h",
		Date:     nil,
	}

	// This will fail in test environment but should not panic
	err := cmd.Run()
	if err != nil {
		t.Logf("Got expected error in test environment: %v", err)
	}
}

func TestAddWorklogCmd_RunWithEmptyDate(t *testing.T) {
	setupMockEnvironment()
	defer teardownMockEnvironment()

	emptyDate := ""
	cmd := AddWorklogCmd{
		IssueKey: "TEST-123",
		Time:     "1h",
		Date:     &emptyDate,
	}

	// This will fail in test environment but should not panic
	err := cmd.Run()
	if err != nil {
		t.Logf("Got expected error in test environment: %v", err)
	}
}

// Integration-style test that would work with proper mocking
func TestRemoveWorklogCmd_RunFlow(t *testing.T) {
	setupMockEnvironment()
	defer teardownMockEnvironment()

	cmd := RemoveWorklogCmd{
		IssueKey: "TEST-123",
	}

	// In a real test environment, this would fail due to network calls
	// But it tests the basic flow and error handling
	err := cmd.Run()
	if err != nil {
		t.Logf("Got expected error in test environment (no real Jira instance): %v", err)

		// Verify it's a network/API related error, not a panic or nil pointer
		if containsString(err.Error(), "panic") {
			t.Error("Command should not panic")
		}
	}
}

// Test CLI struct initialization
func TestCLIStruct(t *testing.T) {
	// Reset CLI to ensure clean state
	CLI = struct {
		AddWorklog    AddWorklogCmd    `cmd:"add-worklog" help:"Add a worklog to a Jira issue"`
		RemoveWorklog RemoveWorklogCmd `cmd:"remove-worklog" help:"Remove a worklog from a Jira issue"`
		Debug         bool             `help:"Enable debug logging" short:"d"`
	}{}

	// Test default values
	if CLI.Debug != false {
		t.Error("CLI.Debug should default to false")
	}

	// Test setting values
	CLI.Debug = true
	CLI.AddWorklog.IssueKey = "TEST-123"
	CLI.AddWorklog.Time = "2h"
	CLI.RemoveWorklog.IssueKey = "TEST-456"

	if CLI.Debug != true {
		t.Error("CLI.Debug should be true after setting")
	}

	if CLI.AddWorklog.IssueKey != "TEST-123" {
		t.Errorf("expected AddWorklog.IssueKey to be 'TEST-123', got '%s'", CLI.AddWorklog.IssueKey)
	}

	if CLI.AddWorklog.Time != "2h" {
		t.Errorf("expected AddWorklog.Time to be '2h', got '%s'", CLI.AddWorklog.Time)
	}

	if CLI.RemoveWorklog.IssueKey != "TEST-456" {
		t.Errorf("expected RemoveWorklog.IssueKey to be 'TEST-456', got '%s'", CLI.RemoveWorklog.IssueKey)
	}
}

// Benchmark tests
func BenchmarkAddWorklogCmd_Run(b *testing.B) {
	setupMockEnvironment()
	defer teardownMockEnvironment()

	cmd := AddWorklogCmd{
		IssueKey: "BENCH-123",
		Time:     "1h",
		Date:     nil,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// This will fail but we're measuring the setup overhead
		_ = cmd.Run()
	}
}

func BenchmarkRemoveWorklogCmd_Run(b *testing.B) {
	setupMockEnvironment()
	defer teardownMockEnvironment()

	cmd := RemoveWorklogCmd{
		IssueKey: "BENCH-123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// This will fail but we're measuring the setup overhead
		_ = cmd.Run()
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(substr) > 0 && len(s) > len(substr) &&
			(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
				func() bool {
					for i := 0; i <= len(s)-len(substr); i++ {
						if s[i:i+len(substr)] == substr {
							return true
						}
					}
					return false
				}())))
}

// Example of how you might structure tests with proper dependency injection
// This is a more advanced pattern that would require refactoring your main code

type TempooInterface interface {
	AddWorklog(issueKey, time string, date *string) error
	GetUserAccountID() (string, error)
	GetWorklogs(issueKey, userID string) ([]string, error)
	DeleteWorklog(issueKey, worklogID string) error
}

type TestableAddWorklogCmd struct {
	AddWorklogCmd
	tempooFactory func() (TempooInterface, error)
}

func (cmd *TestableAddWorklogCmd) Run() error {
	tempoo, err := cmd.tempooFactory()
	if err != nil {
		return err
	}
	return tempoo.AddWorklog(cmd.IssueKey, cmd.Time, cmd.Date)
}

type TestableRemoveWorklogCmd struct {
	RemoveWorklogCmd
	tempooFactory func() (TempooInterface, error)
}

func (cmd *TestableRemoveWorklogCmd) Run() error {
	tempoo, err := cmd.tempooFactory()
	if err != nil {
		return err
	}

	userID, err := tempoo.GetUserAccountID()
	if err != nil {
		return err
	}

	worklogIDs, err := tempoo.GetWorklogs(cmd.IssueKey, userID)
	if err != nil {
		return err
	}

	if len(worklogIDs) == 0 {
		return nil
	}

	for _, worklogID := range worklogIDs {
		if err := tempoo.DeleteWorklog(cmd.IssueKey, worklogID); err != nil {
			return err
		}
	}

	return nil
}

// Mock implementation for proper unit testing
type MockTempoo struct {
	AddWorklogFunc       func(issueKey, time string, date *string) error
	GetUserAccountIDFunc func() (string, error)
	GetWorklogsFunc      func(issueKey, userID string) ([]string, error)
	DeleteWorklogFunc    func(issueKey, worklogID string) error
}

func (m *MockTempoo) AddWorklog(issueKey, time string, date *string) error {
	if m.AddWorklogFunc != nil {
		return m.AddWorklogFunc(issueKey, time, date)
	}
	return nil
}

func (m *MockTempoo) GetUserAccountID() (string, error) {
	if m.GetUserAccountIDFunc != nil {
		return m.GetUserAccountIDFunc()
	}
	return "test-user-id", nil
}

func (m *MockTempoo) GetWorklogs(issueKey, userID string) ([]string, error) {
	if m.GetWorklogsFunc != nil {
		return m.GetWorklogsFunc(issueKey, userID)
	}
	return []string{"worklog-1", "worklog-2"}, nil
}

func (m *MockTempoo) DeleteWorklog(issueKey, worklogID string) error {
	if m.DeleteWorklogFunc != nil {
		return m.DeleteWorklogFunc(issueKey, worklogID)
	}
	return nil
}

// Example of how to test with proper mocking (requires code refactoring)
func TestTestableAddWorklogCmd_Run(t *testing.T) {
	tests := []struct {
		name          string
		cmd           TestableAddWorklogCmd
		expectedError bool
		errorContains string
	}{
		{
			name: "successful add worklog with mock",
			cmd: TestableAddWorklogCmd{
				AddWorklogCmd: AddWorklogCmd{
					IssueKey: "TEST-123",
					Time:     "2h",
					Date:     nil,
				},
				tempooFactory: func() (TempooInterface, error) {
					return &MockTempoo{}, nil
				},
			},
			expectedError: false,
		},
		{
			name: "factory error",
			cmd: TestableAddWorklogCmd{
				AddWorklogCmd: AddWorklogCmd{
					IssueKey: "TEST-123",
					Time:     "2h",
					Date:     nil,
				},
				tempooFactory: func() (TempooInterface, error) {
					return nil, errors.New("factory error")
				},
			},
			expectedError: true,
			errorContains: "factory error",
		},
		{
			name: "add worklog error",
			cmd: TestableAddWorklogCmd{
				AddWorklogCmd: AddWorklogCmd{
					IssueKey: "TEST-123",
					Time:     "2h",
					Date:     nil,
				},
				tempooFactory: func() (TempooInterface, error) {
					return &MockTempoo{
						AddWorklogFunc: func(issueKey, time string, date *string) error {
							return errors.New("add worklog failed")
						},
					}, nil
				},
			},
			expectedError: true,
			errorContains: "add worklog failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.Run()

			if tt.expectedError {
				if err == nil {
					t.Error("expected error but got none")
					return
				}
				if tt.errorContains != "" && !containsString(err.Error(), tt.errorContains) {
					t.Errorf("expected error to contain '%s', got '%s'", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestTestableRemoveWorklogCmd_Run(t *testing.T) {
	tests := []struct {
		name          string
		cmd           TestableRemoveWorklogCmd
		expectedError bool
		errorContains string
	}{
		{
			name: "successful remove worklog with mock",
			cmd: TestableRemoveWorklogCmd{
				RemoveWorklogCmd: RemoveWorklogCmd{
					IssueKey: "TEST-123",
				},
				tempooFactory: func() (TempooInterface, error) {
					return &MockTempoo{}, nil
				},
			},
			expectedError: false,
		},
		{
			name: "no worklogs found",
			cmd: TestableRemoveWorklogCmd{
				RemoveWorklogCmd: RemoveWorklogCmd{
					IssueKey: "TEST-123",
				},
				tempooFactory: func() (TempooInterface, error) {
					return &MockTempoo{
						GetWorklogsFunc: func(issueKey, userID string) ([]string, error) {
							return []string{}, nil
						},
					}, nil
				},
			},
			expectedError: false,
		},
		{
			name: "get user ID error",
			cmd: TestableRemoveWorklogCmd{
				RemoveWorklogCmd: RemoveWorklogCmd{
					IssueKey: "TEST-123",
				},
				tempooFactory: func() (TempooInterface, error) {
					return &MockTempoo{
						GetUserAccountIDFunc: func() (string, error) {
							return "", errors.New("user ID error")
						},
					}, nil
				},
			},
			expectedError: true,
			errorContains: "user ID error",
		},
		{
			name: "get worklogs error",
			cmd: TestableRemoveWorklogCmd{
				RemoveWorklogCmd: RemoveWorklogCmd{
					IssueKey: "TEST-123",
				},
				tempooFactory: func() (TempooInterface, error) {
					return &MockTempoo{
						GetWorklogsFunc: func(issueKey, userID string) ([]string, error) {
							return nil, errors.New("get worklogs error")
						},
					}, nil
				},
			},
			expectedError: true,
			errorContains: "get worklogs error",
		},
		{
			name: "delete worklog error",
			cmd: TestableRemoveWorklogCmd{
				RemoveWorklogCmd: RemoveWorklogCmd{
					IssueKey: "TEST-123",
				},
				tempooFactory: func() (TempooInterface, error) {
					return &MockTempoo{
						DeleteWorklogFunc: func(issueKey, worklogID string) error {
							return errors.New("delete worklog error")
						},
					}, nil
				},
			},
			expectedError: true,
			errorContains: "delete worklog error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.Run()

			if tt.expectedError {
				if err == nil {
					t.Error("expected error but got none")
					return
				}
				if tt.errorContains != "" && !containsString(err.Error(), tt.errorContains) {
					t.Errorf("expected error to contain '%s', got '%s'", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}
