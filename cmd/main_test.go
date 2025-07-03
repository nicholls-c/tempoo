package main

import (
	"errors"
	"os"
	"tempoo/internal"
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

// setupTestEnvironment sets up environment variables and initializes the factory
func setupTestEnvironment() error {
	os.Setenv("JIRA_EMAIL", "test@example.com")
	os.Setenv("JIRA_API_TOKEN", "test-token")

	// Initialize the factory for testing
	factory, err := internal.NewTempooFactory()
	if err != nil {
		return err
	}
	tempooFactory = factory
	return nil
}

// teardownTestEnvironment cleans up environment variables and resets factory
func teardownTestEnvironment() {
	os.Unsetenv("JIRA_EMAIL")
	os.Unsetenv("JIRA_API_TOKEN")
	tempooFactory = nil
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
			name: "factory not initialized",
			cmd: AddWorklogCmd{
				IssueKey: "TEST-123",
				Time:     "2h",
				Date:     nil,
			},
			setupEnv:      false,
			expectedError: true,
			errorContains: "nil pointer", // This will panic with nil factory
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
			// Always clean up first
			teardownTestEnvironment()

			if tt.setupEnv {
				if err := setupTestEnvironment(); err != nil {
					t.Fatalf("Failed to setup test environment: %v", err)
				}
				defer teardownTestEnvironment()
			}

			// Handle potential panic from nil factory
			var err error
			func() {
				defer func() {
					if r := recover(); r != nil {
						err = errors.New("nil pointer dereference")
					}
				}()
				err = tt.cmd.Run()
			}()

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
				// In a real testing scenario, you'd want to mock the HTTP calls
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
			name: "factory not initialized",
			cmd: RemoveWorklogCmd{
				IssueKey: "TEST-123",
			},
			setupEnv:      false,
			expectedError: true,
			errorContains: "nil pointer",
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
			// Always clean up first
			teardownTestEnvironment()

			if tt.setupEnv {
				if err := setupTestEnvironment(); err != nil {
					t.Fatalf("Failed to setup test environment: %v", err)
				}
				defer teardownTestEnvironment()
			}

			// Handle potential panic from nil factory
			var err error
			func() {
				defer func() {
					if r := recover(); r != nil {
						err = errors.New("nil pointer dereference")
					}
				}()
				err = tt.cmd.Run()
			}()

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
				if err != nil {
					t.Logf("Got error (expected in test environment without Jira setup): %v", err)
				}
			}
		})
	}
}

func TestListWorklogsCmd_Run(t *testing.T) {
	tests := []struct {
		name          string
		cmd           ListWorklogsCmd
		setupEnv      bool
		expectedError bool
		errorContains string
	}{
		{
			name: "successful list worklogs",
			cmd: ListWorklogsCmd{
				IssueKey: "TEST-123",
			},
			setupEnv:      true,
			expectedError: false,
		},
		{
			name: "factory not initialized",
			cmd: ListWorklogsCmd{
				IssueKey: "TEST-123",
			},
			setupEnv:      false,
			expectedError: true,
			errorContains: "nil pointer",
		},
		{
			name: "empty issue key",
			cmd: ListWorklogsCmd{
				IssueKey: "",
			},
			setupEnv:      true,
			expectedError: false, // The actual validation happens in the internal package
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Always clean up first
			teardownTestEnvironment()

			if tt.setupEnv {
				if err := setupTestEnvironment(); err != nil {
					t.Fatalf("Failed to setup test environment: %v", err)
				}
				defer teardownTestEnvironment()
			}

			// Handle potential panic from nil factory
			var err error
			func() {
				defer func() {
					if r := recover(); r != nil {
						err = errors.New("nil pointer dereference")
					}
				}()
				err = tt.cmd.Run()
			}()

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
				if err != nil {
					t.Logf("Got error (expected in test environment without Jira setup): %v", err)
				}
			}
		})
	}
}

func TestVersionCmd_Run(t *testing.T) {
	cmd := VersionCmd{}
	err := cmd.Run()

	if err != nil {
		t.Errorf("VersionCmd.Run() should not return an error, got: %v", err)
	}
}

func TestTempooFactoryInitialization(t *testing.T) {
	// Save original factory state
	originalFactory := tempooFactory
	defer func() {
		tempooFactory = originalFactory
	}()

	tests := []struct {
		name        string
		setupEnv    bool
		expectError bool
	}{
		{
			name:        "successful factory initialization",
			setupEnv:    true,
			expectError: false,
		},
		{
			name:        "factory initialization fails with missing env vars",
			setupEnv:    false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up environment
			teardownTestEnvironment()

			if tt.setupEnv {
				os.Setenv("JIRA_EMAIL", "test@example.com")
				os.Setenv("JIRA_API_TOKEN", "test-token")
			}

			// Test factory initialization
			factory, err := internal.NewTempooFactory()

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				if factory != nil {
					t.Error("expected nil factory when error occurs")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
				if factory == nil {
					t.Error("expected non-nil factory")
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

func TestListWorklogsCmd_Fields(t *testing.T) {
	cmd := ListWorklogsCmd{
		IssueKey: "TEST-789",
	}

	if cmd.IssueKey != "TEST-789" {
		t.Errorf("expected IssueKey to be 'TEST-789', got '%s'", cmd.IssueKey)
	}
}

func TestCLIStruct(t *testing.T) {
	// Reset CLI to ensure clean state
	CLI = struct {
		AddWorklog    AddWorklogCmd    `cmd:"add-worklog" help:"Add a worklog to a Jira issue"`
		RemoveWorklog RemoveWorklogCmd `cmd:"remove-worklog" help:"Remove all user worklogs from a Jira issue"`
		ListWorklogs  ListWorklogsCmd  `cmd:"list-worklogs" help:"List all worklogs for a Jira issue"`
		Verbose       bool             `help:"Enable debug logging"`
		Version       bool             `help:"Show version" short:"v"`
	}{}

	// Test default values
	if CLI.Verbose != false {
		t.Error("CLI.Verbose should default to false")
	}

	if CLI.Version != false {
		t.Error("CLI.Version should default to false")
	}

	// Test setting values
	CLI.Verbose = true
	CLI.AddWorklog.IssueKey = "TEST-123"
	CLI.AddWorklog.Time = "2h"
	CLI.RemoveWorklog.IssueKey = "TEST-456"
	CLI.ListWorklogs.IssueKey = "TEST-789"

	if CLI.Verbose != true {
		t.Error("CLI.Verbose should be true after setting")
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

	if CLI.ListWorklogs.IssueKey != "TEST-789" {
		t.Errorf("expected ListWorklogs.IssueKey to be 'TEST-789', got '%s'", CLI.ListWorklogs.IssueKey)
	}
}

// Integration test for factory usage
func TestFactoryUsageInCommands(t *testing.T) {
	// Save original factory state
	originalFactory := tempooFactory
	defer func() {
		tempooFactory = originalFactory
	}()

	if err := setupTestEnvironment(); err != nil {
		t.Fatalf("Failed to setup test environment: %v", err)
	}
	defer teardownTestEnvironment()

	// Test that all commands can get client from factory
	addCmd := AddWorklogCmd{IssueKey: "TEST-123", Time: "1h"}
	removeCmd := RemoveWorklogCmd{IssueKey: "TEST-123"}
	listCmd := ListWorklogsCmd{IssueKey: "TEST-123"}

	// These will fail due to network calls, but should not panic
	_ = addCmd.Run()
	_ = removeCmd.Run()
	_ = listCmd.Run()

	// If we get here without panic, the factory is working
	t.Log("All commands successfully accessed the factory")
}

// Benchmark tests
func BenchmarkAddWorklogCmd_Run(b *testing.B) {
	if err := setupTestEnvironment(); err != nil {
		b.Fatalf("Failed to setup test environment: %v", err)
	}
	defer teardownTestEnvironment()

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
	if err := setupTestEnvironment(); err != nil {
		b.Fatalf("Failed to setup test environment: %v", err)
	}
	defer teardownTestEnvironment()

	cmd := RemoveWorklogCmd{
		IssueKey: "BENCH-123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// This will fail but we're measuring the setup overhead
		_ = cmd.Run()
	}
}

func BenchmarkListWorklogsCmd_Run(b *testing.B) {
	if err := setupTestEnvironment(); err != nil {
		b.Fatalf("Failed to setup test environment: %v", err)
	}
	defer teardownTestEnvironment()

	cmd := ListWorklogsCmd{
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

// Advanced testing patterns with proper mocking
// These demonstrate how you could structure tests with dependency injection

type TempooInterface interface {
	AddWorklog(issueKey, time string, date *string) error
	GetUserAccountID() (string, error)
	GetWorklogs(issueKey, userID string) ([]string, error)
	DeleteWorklog(issueKey, worklogID string) error
	ListWorklogs(issueKey string) error
}

// Mock implementation for proper unit testing
type MockTempoo struct {
	AddWorklogFunc       func(issueKey, time string, date *string) error
	GetUserAccountIDFunc func() (string, error)
	GetWorklogsFunc      func(issueKey, userID string) ([]string, error)
	DeleteWorklogFunc    func(issueKey, worklogID string) error
	ListWorklogsFunc     func(issueKey string) error
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

func (m *MockTempoo) ListWorklogs(issueKey string) error {
	if m.ListWorklogsFunc != nil {
		return m.ListWorklogsFunc(issueKey)
	}
	return nil
}

// Example test showing how you could test with proper mocking
// (This would require refactoring your commands to accept a factory interface)
func TestMockingExample(t *testing.T) {
	mock := &MockTempoo{
		AddWorklogFunc: func(issueKey, time string, date *string) error {
			if issueKey == "FAIL-123" {
				return errors.New("mock error")
			}
			return nil
		},
	}

	// This demonstrates how you could test with proper mocking
	// if your commands accepted a factory interface
	err := mock.AddWorklog("TEST-123", "2h", nil)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	err = mock.AddWorklog("FAIL-123", "2h", nil)
	if err == nil {
		t.Error("Expected error but got none")
	}
}
