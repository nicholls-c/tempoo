package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestMain(m *testing.M) {
	// Set up test environment
	log.SetHandler(discard.New()) // Disable logging during tests

	// Store original environment variables
	originalEmail := os.Getenv("JIRA_EMAIL")
	originalToken := os.Getenv("JIRA_API_TOKEN")

	// Set test environment variables
	os.Setenv("JIRA_EMAIL", "test@example.com")
	os.Setenv("JIRA_API_TOKEN", "test-token")

	// Run tests
	code := m.Run()

	// Restore original environment variables
	if originalEmail != "" {
		os.Setenv("JIRA_EMAIL", originalEmail)
	} else {
		os.Unsetenv("JIRA_EMAIL")
	}

	if originalToken != "" {
		os.Setenv("JIRA_API_TOKEN", originalToken)
	} else {
		os.Unsetenv("JIRA_API_TOKEN")
	}

	os.Exit(code)
}

func TestVersionCmd_Run(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Test version command
	cmd := &VersionCmd{}
	err := cmd.Run()

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	// Assertions
	assert.NoError(t, err)
	assert.Contains(t, output, version)
}

func TestAddWorklogCmd_Run_MissingIssueKey(t *testing.T) {
	var stderr bytes.Buffer

	cmd := &AddWorklogCmd{
		IssueKey: "",
		Hours:    "1.5",
		Date:     nil,
	}

	// Create a mock context
	parser := kong.Must(&CLI)
	ctx, err := kong.Trace(parser, []string{"add-worklog"})
	require.NoError(t, err)
	ctx.Stderr = &stderr

	err = cmd.Run(ctx)
	assert.NoError(t, err) // Should not error, just print usage
	assert.Contains(t, stderr.String(), "Usage:")
}

func TestAddWorklogCmd_Run_MissingHours(t *testing.T) {
	var stderr bytes.Buffer

	cmd := &AddWorklogCmd{
		IssueKey: "TEST-123",
		Hours:    "",
		Date:     nil,
	}

	// Create a mock context
	parser := kong.Must(&CLI)
	ctx, err := kong.Trace(parser, []string{"add-worklog"})
	require.NoError(t, err)
	ctx.Stderr = &stderr

	err = cmd.Run(ctx)
	assert.NoError(t, err) // Should not error, just print usage
	assert.Contains(t, stderr.String(), "Usage:")
}

func TestAddWorklogCmd_Run_ValidInput(t *testing.T) {
	// Reset factory for clean test
	tempooFactory = nil

	cmd := &AddWorklogCmd{
		IssueKey: "TEST-123",
		Hours:    "1.5",
		Date:     nil,
	}

	// Create a mock context
	parser := kong.Must(&CLI)
	ctx, err := kong.Trace(parser, []string{"add-worklog"})
	require.NoError(t, err)

	err = cmd.Run(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to add worklog")
}

func TestRemoveWorklogsCmd_Run_MissingIssueKey(t *testing.T) {
	var stderr bytes.Buffer

	cmd := &RemoveWorklogsCmd{
		IssueKey: "",
	}

	// Create a mock context
	parser := kong.Must(&CLI)
	ctx, err := kong.Trace(parser, []string{"remove-worklogs"})
	require.NoError(t, err)
	ctx.Stderr = &stderr

	err = cmd.Run(ctx)
	assert.NoError(t, err) // Should not error, just print usage
	assert.Contains(t, stderr.String(), "Usage:")
}

func TestListWorklogsCmd_Run_MissingIssueKey(t *testing.T) {
	var stderr bytes.Buffer

	cmd := &ListWorklogsCmd{
		IssueKey: "",
	}

	// Create a mock context
	parser := kong.Must(&CLI)
	ctx, err := kong.Trace(parser, []string{"list-worklogs"})
	require.NoError(t, err)
	ctx.Stderr = &stderr

	err = cmd.Run(ctx)
	assert.NoError(t, err) // Should not error, just print usage
	assert.Contains(t, stderr.String(), "Usage:")
}

func TestGetFactory_Success(t *testing.T) {
	// Reset factory for clean test
	tempooFactory = nil

	// Set valid environment variables
	os.Setenv("JIRA_EMAIL", "test@example.com")
	os.Setenv("JIRA_API_TOKEN", "test-token")

	factory, err := getFactory()
	assert.NoError(t, err)
	assert.NotNil(t, factory)

	// Test singleton behavior - should return same instance
	factory2, err := getFactory()
	assert.NoError(t, err)
	assert.Equal(t, factory, factory2)
}

func TestGetFactory_MissingCredentials(t *testing.T) {
	// Reset factory for clean test
	tempooFactory = nil

	// Unset environment variables
	os.Unsetenv("JIRA_EMAIL")
	os.Unsetenv("JIRA_API_TOKEN")

	factory, err := getFactory()
	assert.Error(t, err)
	assert.Nil(t, factory)
	assert.Contains(t, err.Error(), "failed to initialize Tempoo factory")
}

func TestCLI_Structure(t *testing.T) {
	// Test that CLI struct has all expected commands
	parser := kong.Must(&CLI,
		kong.Name("tempoo"),
		kong.Description("tem💩, because life is too short.\n\nA CLI tool for managing Jira worklogs."),
		kong.UsageOnError(),
	)

	// Test parsing different commands
	testCases := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "add-worklog command",
			args:     []string{"add-worklog", "-i", "TEST-123", "-t", "1.5"},
			expected: "add-worklog",
		},
		{
			name:     "remove-worklogs command",
			args:     []string{"remove-worklogs", "-i", "TEST-123"},
			expected: "remove-worklogs",
		},
		{
			name:     "list-worklogs command",
			args:     []string{"list-worklogs", "-i", "TEST-123"},
			expected: "list-worklogs",
		},
		{
			name:     "version command",
			args:     []string{"version"},
			expected: "version",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, err := kong.Trace(parser, tc.args)
			assert.NoError(t, err)
			assert.Contains(t, ctx.Command(), tc.expected)
		})
	}
}

func TestVersion_GlobalVariable(t *testing.T) {
	// Test that version variable exists and has expected format
	assert.NotEmpty(t, version)
	assert.True(t, strings.Contains(version, "dev") || strings.Contains(version, "."))
}

// Benchmark tests
func BenchmarkVersionCmd_Run(b *testing.B) {
	cmd := &VersionCmd{}

	// Redirect stdout to discard for benchmarking
	oldStdout := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() { os.Stdout = oldStdout }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cmd.Run()
	}
}

func BenchmarkGetFactory(b *testing.B) {
	// Set up environment
	os.Setenv("JIRA_EMAIL", "test@example.com")
	os.Setenv("JIRA_API_TOKEN", "test-token")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tempooFactory = nil // Reset for each iteration
		_, _ = getFactory()
	}
}

// Integration test for CLI parsing
func TestCLI_Integration(t *testing.T) {
	testCases := []struct {
		name        string
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid add-worklog with all flags",
			args:        []string{"add-worklog", "--issue-key", "TEST-123", "--hours", "1.5", "--date", "01.01.2024"},
			expectError: false,
		},
		{
			name:        "valid add-worklog with short flags",
			args:        []string{"add-worklog", "-i", "TEST-123", "-t", "2.5"},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := kong.Must(&CLI)
			_, err := kong.Trace(parser, tc.args)

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, strings.ToLower(err.Error()), strings.ToLower(tc.errorMsg))
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
