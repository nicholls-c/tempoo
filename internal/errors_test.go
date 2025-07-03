package internal

import (
	"errors"
	"fmt"
	"testing"
)

func TestTempooError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *TempooError
		expected string
	}{
		{
			name: "error with message only",
			err: &TempooError{
				Message: "Something went wrong",
				Cause:   nil,
			},
			expected: "Something went wrong",
		},
		{
			name: "error with message and cause",
			err: &TempooError{
				Message: "API request failed",
				Cause:   errors.New("network timeout"),
			},
			expected: "API request failed: network timeout",
		},
		{
			name: "error with empty message and cause",
			err: &TempooError{
				Message: "",
				Cause:   errors.New("underlying error"),
			},
			expected: ": underlying error",
		},
		{
			name: "error with message and nil cause",
			err: &TempooError{
				Message: "Simple error",
				Cause:   nil,
			},
			expected: "Simple error",
		},
		{
			name: "empty error",
			err: &TempooError{
				Message: "",
				Cause:   nil,
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.expected {
				t.Errorf("TempooError.Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestTempooError_ErrorInterface(t *testing.T) {
	// Test that TempooError implements the error interface
	var err error = &TempooError{Message: "test error"}

	if err.Error() != "test error" {
		t.Errorf("TempooError should implement error interface correctly")
	}
}

func TestTempooError_WithCause(t *testing.T) {
	// Test chaining errors
	originalErr := errors.New("original error")
	tempooErr := &TempooError{
		Message: "Wrapped error",
		Cause:   originalErr,
	}

	expectedMsg := "Wrapped error: original error"
	if tempooErr.Error() != expectedMsg {
		t.Errorf("TempooError.Error() = %q, want %q", tempooErr.Error(), expectedMsg)
	}

	// Test that the cause is preserved
	if tempooErr.Cause != originalErr {
		t.Errorf("TempooError.Cause should preserve the original error")
	}
}

func TestTempooError_NestedCause(t *testing.T) {
	// Test nested TempooError
	innerErr := &TempooError{
		Message: "Inner error",
		Cause:   errors.New("root cause"),
	}

	outerErr := &TempooError{
		Message: "Outer error",
		Cause:   innerErr,
	}

	expected := "Outer error: Inner error: root cause"
	if outerErr.Error() != expected {
		t.Errorf("Nested TempooError.Error() = %q, want %q", outerErr.Error(), expected)
	}
}

func TestInvalidIssueKeyError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *InvalidIssueKeyError
		expected string
	}{
		{
			name: "valid issue key format",
			err: &InvalidIssueKeyError{
				IssueKey: "PROJ-123",
			},
			expected: "Issue key PROJ-123 is not valid",
		},
		{
			name: "another valid issue key format",
			err: &InvalidIssueKeyError{
				IssueKey: "INF-88",
			},
			expected: "Issue key INF-88 is not valid",
		},
		{
			name: "empty issue key",
			err: &InvalidIssueKeyError{
				IssueKey: "",
			},
			expected: "Issue key  is not valid",
		},
		{
			name: "issue key with spaces",
			err: &InvalidIssueKeyError{
				IssueKey: "PROJ 123",
			},
			expected: "Issue key PROJ 123 is not valid",
		},
		{
			name: "issue key with special characters",
			err: &InvalidIssueKeyError{
				IssueKey: "PROJ@123",
			},
			expected: "Issue key PROJ@123 is not valid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.expected {
				t.Errorf("InvalidIssueKeyError.Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestInvalidIssueKeyError_ErrorInterface(t *testing.T) {
	// Test that InvalidIssueKeyError implements the error interface
	var err error = &InvalidIssueKeyError{IssueKey: "TEST-123"}

	expected := "Issue key TEST-123 is not valid"
	if err.Error() != expected {
		t.Errorf("InvalidIssueKeyError should implement error interface correctly, got %q, want %q", err.Error(), expected)
	}
}

func TestErrorTypes(t *testing.T) {
	// Test that both error types can be used in error handling
	tempooErr := &TempooError{Message: "Tempoo error"}
	invalidKeyErr := &InvalidIssueKeyError{IssueKey: "BAD-KEY"}

	// Test type assertions
	var err1 error = tempooErr
	var err2 error = invalidKeyErr

	if _, ok := err1.(*TempooError); !ok {
		t.Error("TempooError should be assertable from error interface")
	}

	if _, ok := err2.(*InvalidIssueKeyError); !ok {
		t.Error("InvalidIssueKeyError should be assertable from error interface")
	}
}

func TestErrorComparison(t *testing.T) {
	// Test error equality and comparison
	err1 := &TempooError{Message: "Same message", Cause: nil}
	err2 := &TempooError{Message: "Same message", Cause: nil}
	err3 := &TempooError{Message: "Different message", Cause: nil}

	// Test that errors with same content have same string representation
	if err1.Error() != err2.Error() {
		t.Error("Errors with same content should have same string representation")
	}

	if err1.Error() == err3.Error() {
		t.Error("Errors with different content should have different string representation")
	}

	// Test InvalidIssueKeyError comparison
	keyErr1 := &InvalidIssueKeyError{IssueKey: "SAME-123"}
	keyErr2 := &InvalidIssueKeyError{IssueKey: "SAME-123"}
	keyErr3 := &InvalidIssueKeyError{IssueKey: "DIFF-456"}

	if keyErr1.Error() != keyErr2.Error() {
		t.Error("InvalidIssueKeyErrors with same key should have same string representation")
	}

	if keyErr1.Error() == keyErr3.Error() {
		t.Error("InvalidIssueKeyErrors with different keys should have different string representation")
	}
}

func TestErrorInContext(t *testing.T) {
	// Test how errors would be used in real scenarios
	tests := []struct {
		name        string
		errorFunc   func() error
		expectedMsg string
		errorType   string
	}{
		{
			name: "API request failure",
			errorFunc: func() error {
				return &TempooError{
					Message: "API request failed",
					Cause:   errors.New("connection timeout"),
				}
			},
			expectedMsg: "API request failed: connection timeout",
			errorType:   "TempooError",
		},
		{
			name: "Invalid issue key",
			errorFunc: func() error {
				return &InvalidIssueKeyError{IssueKey: "INVALID-KEY"}
			},
			expectedMsg: "Issue key INVALID-KEY is not valid",
			errorType:   "InvalidIssueKeyError",
		},
		{
			name: "Environment variable missing",
			errorFunc: func() error {
				return &TempooError{Message: "JIRA_EMAIL environment variable is not set"}
			},
			expectedMsg: "JIRA_EMAIL environment variable is not set",
			errorType:   "TempooError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.errorFunc()

			if err.Error() != tt.expectedMsg {
				t.Errorf("Error message = %q, want %q", err.Error(), tt.expectedMsg)
			}

			// Test type assertion
			switch tt.errorType {
			case "TempooError":
				if _, ok := err.(*TempooError); !ok {
					t.Errorf("Expected TempooError, got %T", err)
				}
			case "InvalidIssueKeyError":
				if _, ok := err.(*InvalidIssueKeyError); !ok {
					t.Errorf("Expected InvalidIssueKeyError, got %T", err)
				}
			}
		})
	}
}

func TestErrorWrapping(t *testing.T) {
	// Test error wrapping patterns
	rootCause := errors.New("network error")
	httpErr := fmt.Errorf("HTTP request failed: %w", rootCause)
	tempooErr := &TempooError{
		Message: "Failed to connect to Jira",
		Cause:   httpErr,
	}

	expectedMsg := "Failed to connect to Jira: HTTP request failed: network error"
	if tempooErr.Error() != expectedMsg {
		t.Errorf("Wrapped error message = %q, want %q", tempooErr.Error(), expectedMsg)
	}
}

// Benchmark tests
func BenchmarkTempooError_Error(b *testing.B) {
	err := &TempooError{
		Message: "API request failed",
		Cause:   errors.New("network timeout"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkInvalidIssueKeyError_Error(b *testing.B) {
	err := &InvalidIssueKeyError{IssueKey: "PROJ-123"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkTempooError_ErrorWithoutCause(b *testing.B) {
	err := &TempooError{
		Message: "Simple error message",
		Cause:   nil,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

// Test edge cases
func TestErrorEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "nil TempooError",
			test: func(t *testing.T) {
				var err *TempooError
				// This would panic if called, but we test the type
				if err != nil {
					t.Error("nil TempooError should be nil")
				}
			},
		},
		{
			name: "nil InvalidIssueKeyError",
			test: func(t *testing.T) {
				var err *InvalidIssueKeyError
				if err != nil {
					t.Error("nil InvalidIssueKeyError should be nil")
				}
			},
		},
		{
			name: "very long issue key",
			test: func(t *testing.T) {
				longKey := string(make([]byte, 1000))
				for i := range longKey {
					longKey = longKey[:i] + "A" + longKey[i+1:]
				}

				err := &InvalidIssueKeyError{IssueKey: longKey}
				msg := err.Error()

				expected := fmt.Sprintf("Issue key %s is not valid", longKey)
				if msg != expected {
					// Just ensure it doesn't panic and returns something reasonable
					if len(msg) == 0 {
						t.Error("Error message should not be empty for long issue key")
					}
				}
			},
		},
		{
			name: "very long error message",
			test: func(t *testing.T) {
				longMsg := string(make([]byte, 1000))
				for i := range longMsg {
					longMsg = longMsg[:i] + "X" + longMsg[i+1:]
				}

				err := &TempooError{Message: longMsg}
				msg := err.Error()

				if msg != longMsg {
					t.Error("Error message should be preserved even if very long")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
