package internal

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestNewTempoo_Success(t *testing.T) {
	// Set up environment variables
	originalEmail := os.Getenv("JIRA_EMAIL")
	originalToken := os.Getenv("JIRA_API_TOKEN")

	defer func() {
		// Restore original environment
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
	}()

	// Test successful creation
	testEmail := "test@example.com"
	testToken := "test-api-token"

	os.Setenv("JIRA_EMAIL", testEmail)
	os.Setenv("JIRA_API_TOKEN", testToken)

	tempoo, err := NewTempoo()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if tempoo == nil {
		t.Fatal("Expected Tempoo instance, got nil")
	}

	// Verify the struct type
	if reflect.TypeOf(tempoo).String() != "*internal.Tempoo" {
		t.Errorf("Expected *internal.Tempoo, got %T", tempoo)
	}

	// Test that the struct has the expected fields using reflection
	tempooValue := reflect.ValueOf(tempoo).Elem()

	// Check email field
	emailField := tempooValue.FieldByName("email")
	if !emailField.IsValid() {
		t.Error("email field should be valid")
	}
	if emailField.String() != testEmail {
		t.Errorf("Expected email %s, got %s", testEmail, emailField.String())
	}

	// Check apiToken field
	tokenField := tempooValue.FieldByName("apiToken")
	if !tokenField.IsValid() {
		t.Error("apiToken field should be valid")
	}
	if tokenField.String() != testToken {
		t.Errorf("Expected apiToken %s, got %s", testToken, tokenField.String())
	}

	// Check client field
	clientField := tempooValue.FieldByName("client")
	if !clientField.IsValid() {
		t.Error("client field should be valid")
	}
	if clientField.IsNil() {
		t.Error("client field should not be nil")
	}

	// Verify it's a resty client
	if clientField.Type().String() != "*resty.Client" {
		t.Errorf("Expected *resty.Client, got %s", clientField.Type().String())
	}
}

func TestNewTempoo_MissingJiraEmail(t *testing.T) {
	// Save original environment
	originalEmail := os.Getenv("JIRA_EMAIL")
	originalToken := os.Getenv("JIRA_API_TOKEN")

	defer func() {
		// Restore original environment
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
	}()

	// Unset JIRA_EMAIL but set JIRA_API_TOKEN
	os.Unsetenv("JIRA_EMAIL")
	os.Setenv("JIRA_API_TOKEN", "test-token")

	tempoo, err := NewTempoo()

	if tempoo != nil {
		t.Error("Expected nil Tempoo instance when JIRA_EMAIL is missing")
	}

	if err == nil {
		t.Error("Expected error when JIRA_EMAIL is missing")
	}

	// Check error type and message
	if tempooErr, ok := err.(*TempooError); ok {
		expectedMsg := "JIRA_EMAIL environment variable is not set"
		if tempooErr.Message != expectedMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedMsg, tempooErr.Message)
		}
	} else {
		t.Errorf("Expected TempooError, got %T", err)
	}
}

func TestNewTempoo_MissingJiraAPIToken(t *testing.T) {
	// Save original environment
	originalEmail := os.Getenv("JIRA_EMAIL")
	originalToken := os.Getenv("JIRA_API_TOKEN")

	defer func() {
		// Restore original environment
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
	}()

	// Set JIRA_EMAIL but unset JIRA_API_TOKEN
	os.Setenv("JIRA_EMAIL", "test@example.com")
	os.Unsetenv("JIRA_API_TOKEN")

	tempoo, err := NewTempoo()

	if tempoo != nil {
		t.Error("Expected nil Tempoo instance when JIRA_API_TOKEN is missing")
	}

	if err == nil {
		t.Error("Expected error when JIRA_API_TOKEN is missing")
	}

	// Check error type and message
	if tempooErr, ok := err.(*TempooError); ok {
		expectedMsg := "JIRA_API_TOKEN environment variable is not set"
		if tempooErr.Message != expectedMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedMsg, tempooErr.Message)
		}
	} else {
		t.Errorf("Expected TempooError, got %T", err)
	}
}

func TestNewTempoo_EmptyJiraEmail(t *testing.T) {
	// Save original environment
	originalEmail := os.Getenv("JIRA_EMAIL")
	originalToken := os.Getenv("JIRA_API_TOKEN")

	defer func() {
		// Restore original environment
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
	}()

	// Set empty JIRA_EMAIL
	os.Setenv("JIRA_EMAIL", "")
	os.Setenv("JIRA_API_TOKEN", "test-token")

	tempoo, err := NewTempoo()

	if tempoo != nil {
		t.Error("Expected nil Tempoo instance when JIRA_EMAIL is empty")
	}

	if err == nil {
		t.Error("Expected error when JIRA_EMAIL is empty")
	}

	// Check error type and message
	if tempooErr, ok := err.(*TempooError); ok {
		expectedMsg := "JIRA_EMAIL environment variable is not set"
		if tempooErr.Message != expectedMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedMsg, tempooErr.Message)
		}
	} else {
		t.Errorf("Expected TempooError, got %T", err)
	}
}

func TestNewTempoo_EmptyJiraAPIToken(t *testing.T) {
	// Save original environment
	originalEmail := os.Getenv("JIRA_EMAIL")
	originalToken := os.Getenv("JIRA_API_TOKEN")

	defer func() {
		// Restore original environment
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
	}()

	// Set empty JIRA_API_TOKEN
	os.Setenv("JIRA_EMAIL", "test@example.com")
	os.Setenv("JIRA_API_TOKEN", "")

	tempoo, err := NewTempoo()

	if tempoo != nil {
		t.Error("Expected nil Tempoo instance when JIRA_API_TOKEN is empty")
	}

	if err == nil {
		t.Error("Expected error when JIRA_API_TOKEN is empty")
	}

	// Check error type and message
	if tempooErr, ok := err.(*TempooError); ok {
		expectedMsg := "JIRA_API_TOKEN environment variable is not set"
		if tempooErr.Message != expectedMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedMsg, tempooErr.Message)
		}
	} else {
		t.Errorf("Expected TempooError, got %T", err)
	}
}

func TestNewTempoo_BothEnvironmentVariablesMissing(t *testing.T) {
	// Save original environment
	originalEmail := os.Getenv("JIRA_EMAIL")
	originalToken := os.Getenv("JIRA_API_TOKEN")

	defer func() {
		// Restore original environment
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
	}()

	// Unset both environment variables
	os.Unsetenv("JIRA_EMAIL")
	os.Unsetenv("JIRA_API_TOKEN")

	tempoo, err := NewTempoo()

	if tempoo != nil {
		t.Error("Expected nil Tempoo instance when both environment variables are missing")
	}

	if err == nil {
		t.Error("Expected error when both environment variables are missing")
	}

	// Should fail on JIRA_EMAIL first
	if tempooErr, ok := err.(*TempooError); ok {
		expectedMsg := "JIRA_EMAIL environment variable is not set"
		if tempooErr.Message != expectedMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedMsg, tempooErr.Message)
		}
	} else {
		t.Errorf("Expected TempooError, got %T", err)
	}
}

func TestNewTempoo_ClientConfiguration(t *testing.T) {
	// Save original environment
	originalEmail := os.Getenv("JIRA_EMAIL")
	originalToken := os.Getenv("JIRA_API_TOKEN")

	defer func() {
		// Restore original environment
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
	}()

	testEmail := "test@example.com"
	testToken := "test-api-token"

	os.Setenv("JIRA_EMAIL", testEmail)
	os.Setenv("JIRA_API_TOKEN", testToken)

	tempoo, err := NewTempoo()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Access the client through reflection to verify its configuration
	tempooValue := reflect.ValueOf(tempoo).Elem()
	clientField := tempooValue.FieldByName("client")

	if !clientField.IsValid() || clientField.IsNil() {
		t.Fatal("client field should be valid and not nil")
	}

	// We can't access the client directly since it's unexported,
	// but we can verify it's the correct type
	if clientField.Type().String() != "*resty.Client" {
		t.Errorf("Expected client to be *resty.Client, got %s", clientField.Type().String())
	}
}

func TestNewTempoo_WithSpecialCharacters(t *testing.T) {
	// Save original environment
	originalEmail := os.Getenv("JIRA_EMAIL")
	originalToken := os.Getenv("JIRA_API_TOKEN")

	defer func() {
		// Restore original environment
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
	}()

	// Test with email and token containing special characters
	testEmail := "test+user@example-domain.com"
	testToken := "token-with-special-chars_123!@#$%"

	os.Setenv("JIRA_EMAIL", testEmail)
	os.Setenv("JIRA_API_TOKEN", testToken)

	tempoo, err := NewTempoo()

	if err != nil {
		t.Errorf("Expected no error with special characters, got %v", err)
	}

	if tempoo == nil {
		t.Error("Expected Tempoo instance with special characters")
	}

	// Verify the values are preserved correctly
	tempooValue := reflect.ValueOf(tempoo).Elem()

	emailField := tempooValue.FieldByName("email")
	if emailField.String() != testEmail {
		t.Errorf("Expected email %s, got %s", testEmail, emailField.String())
	}

	tokenField := tempooValue.FieldByName("apiToken")
	if tokenField.String() != testToken {
		t.Errorf("Expected apiToken %s, got %s", testToken, tokenField.String())
	}
}

func TestNewTempoo_MultipleCallsIndependent(t *testing.T) {
	// Save original environment
	originalEmail := os.Getenv("JIRA_EMAIL")
	originalToken := os.Getenv("JIRA_API_TOKEN")

	defer func() {
		// Restore original environment
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
	}()

	testEmail := "test@example.com"
	testToken := "test-api-token"

	os.Setenv("JIRA_EMAIL", testEmail)
	os.Setenv("JIRA_API_TOKEN", testToken)

	// Create multiple instances
	tempoo1, err1 := NewTempoo()
	tempoo2, err2 := NewTempoo()

	if err1 != nil || err2 != nil {
		t.Errorf("Expected no errors, got %v, %v", err1, err2)
	}

	if tempoo1 == nil || tempoo2 == nil {
		t.Error("Expected both Tempoo instances to be non-nil")
	}

	// Verify they are different instances
	if tempoo1 == tempoo2 {
		t.Error("Expected different instances, got the same pointer")
	}

	// Verify they have the same configuration but different clients
	tempoo1Value := reflect.ValueOf(tempoo1).Elem()
	tempoo2Value := reflect.ValueOf(tempoo2).Elem()

	client1Field := tempoo1Value.FieldByName("client")
	client2Field := tempoo2Value.FieldByName("client")

	// We can't directly compare the client instances since they're unexported,
	// but we can verify they are both valid and not nil
	if client1Field.IsNil() || client2Field.IsNil() {
		t.Error("Both client instances should be non-nil")
	}

	// We can verify they have the same type
	if client1Field.Type() != client2Field.Type() {
		t.Error("Both clients should have the same type")
	}
}

// Benchmark tests
func BenchmarkNewTempoo(b *testing.B) {
	// Set up environment
	os.Setenv("JIRA_EMAIL", "benchmark@example.com")
	os.Setenv("JIRA_API_TOKEN", "benchmark-token")

	defer func() {
		os.Unsetenv("JIRA_EMAIL")
		os.Unsetenv("JIRA_API_TOKEN")
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tempoo, err := NewTempoo()
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
		_ = tempoo
	}
}

func BenchmarkNewTempoo_WithMissingEnv(b *testing.B) {
	// Ensure environment variables are not set
	os.Unsetenv("JIRA_EMAIL")
	os.Unsetenv("JIRA_API_TOKEN")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tempoo, err := NewTempoo()
		if err == nil {
			b.Fatal("Expected error but got none")
		}
		_ = tempoo
	}
}

// Test edge cases
func TestNewTempoo_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		token       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "email with unicode",
			email:       "test.üñíçødé@example.com",
			token:       "valid-token",
			expectError: false,
		},
		{
			name:        "token with unicode",
			email:       "test@example.com",
			token:       "tøkën-wíth-üñíçødé",
			expectError: false,
		},
		{
			name:        "whitespace in email",
			email:       " test@example.com ",
			token:       "valid-token",
			expectError: false,
		},
		{
			name:        "whitespace in token",
			email:       "test@example.com",
			token:       " valid-token ",
			expectError: false,
		},
		{
			name:        "very long email",
			email:       strings.Repeat("a", 200) + "@example.com",
			token:       "valid-token",
			expectError: false,
		},
		{
			name:        "very long token",
			email:       "test@example.com",
			token:       strings.Repeat("t", 500),
			expectError: false,
		},
	}

	// Save original environment
	originalEmail := os.Getenv("JIRA_EMAIL")
	originalToken := os.Getenv("JIRA_API_TOKEN")

	defer func() {
		// Restore original environment
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
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("JIRA_EMAIL", tt.email)
			os.Setenv("JIRA_API_TOKEN", tt.token)

			tempoo, err := NewTempoo()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
				if tempoo != nil {
					t.Error("Expected nil Tempoo instance on error")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if tempoo == nil {
					t.Error("Expected non-nil Tempoo instance")
				}

				// Verify the values are preserved
				if tempoo != nil {
					tempooValue := reflect.ValueOf(tempoo).Elem()

					emailField := tempooValue.FieldByName("email")
					if emailField.String() != tt.email {
						t.Errorf("Expected email %s, got %s", tt.email, emailField.String())
					}

					tokenField := tempooValue.FieldByName("apiToken")
					if tokenField.String() != tt.token {
						t.Errorf("Expected token %s, got %s", tt.token, tokenField.String())
					}
				}
			}
		})
	}
}

func TestNewTempoo_ErrorTypes(t *testing.T) {
	// Save original environment
	originalEmail := os.Getenv("JIRA_EMAIL")
	originalToken := os.Getenv("JIRA_API_TOKEN")

	defer func() {
		// Restore original environment
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
	}()

	// Test that errors are of the correct type
	os.Unsetenv("JIRA_EMAIL")
	os.Setenv("JIRA_API_TOKEN", "test-token")

	_, err := NewTempoo()

	if err == nil {
		t.Fatal("Expected error but got none")
	}

	// Test error interface implementation
	var errInterface error = err
	if errInterface.Error() == "" {
		t.Error("Error should implement error interface with non-empty message")
	}

	// Test specific error type
	if _, ok := err.(*TempooError); !ok {
		t.Errorf("Expected *TempooError, got %T", err)
	}
}

func TestNewTempoo_EnvironmentVariableOrder(t *testing.T) {
	// Save original environment
	originalEmail := os.Getenv("JIRA_EMAIL")
	originalToken := os.Getenv("JIRA_API_TOKEN")

	defer func() {
		// Restore original environment
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
	}()

	// Test that JIRA_EMAIL is checked first
	os.Unsetenv("JIRA_EMAIL")
	os.Unsetenv("JIRA_API_TOKEN")

	_, err := NewTempoo()

	if err == nil {
		t.Fatal("Expected error but got none")
	}

	if tempooErr, ok := err.(*TempooError); ok {
		if !strings.Contains(tempooErr.Message, "JIRA_EMAIL") {
			t.Errorf("Expected error about JIRA_EMAIL first, got: %s", tempooErr.Message)
		}
	} else {
		t.Errorf("Expected TempooError, got %T", err)
	}
}
