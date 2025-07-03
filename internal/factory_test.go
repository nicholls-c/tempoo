// internal/factory_test.go
package internal

import (
	"os"
	"reflect"
	"testing"
)

func TestNewTempooFactory_Success(t *testing.T) {
	// Set up environment variables for successful Tempoo creation
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

	// Set valid environment variables
	testEmail := "test@example.com"
	testToken := "test-api-token"
	os.Setenv("JIRA_EMAIL", testEmail)
	os.Setenv("JIRA_API_TOKEN", testToken)

	// Test factory creation
	factory, err := NewTempooFactory()

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if factory == nil {
		t.Fatal("Expected TempooFactory instance, got nil")
	}

	// Verify factory type
	if reflect.TypeOf(factory).String() != "*internal.TempooFactory" {
		t.Errorf("Expected *internal.TempooFactory, got %T", factory)
	}

	// Verify factory has an instance
	if factory.instance == nil {
		t.Error("Expected factory to have a non-nil Tempoo instance")
	}

	// Verify the instance is of correct type
	if reflect.TypeOf(factory.instance).String() != "*internal.Tempoo" {
		t.Errorf("Expected factory.instance to be *internal.Tempoo, got %T", factory.instance)
	}
}

func TestNewTempooFactory_FailsWhenTempooCreationFails(t *testing.T) {
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

	// Unset environment variables to cause NewTempoo() to fail
	os.Unsetenv("JIRA_EMAIL")
	os.Unsetenv("JIRA_API_TOKEN")

	// Test factory creation
	factory, err := NewTempooFactory()

	// Assertions
	if err == nil {
		t.Error("Expected error when Tempoo creation fails, got nil")
	}

	if factory != nil {
		t.Error("Expected nil factory when Tempoo creation fails, got non-nil")
	}

	// Verify error type
	if _, ok := err.(*TempooError); !ok {
		t.Errorf("Expected TempooError, got %T", err)
	}
}

func TestNewTempooFactory_FailsWithMissingJiraEmail(t *testing.T) {
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

	// Set only API token, missing email
	os.Unsetenv("JIRA_EMAIL")
	os.Setenv("JIRA_API_TOKEN", "test-token")

	// Test factory creation
	factory, err := NewTempooFactory()

	// Assertions
	if err == nil {
		t.Error("Expected error when JIRA_EMAIL is missing, got nil")
	}

	if factory != nil {
		t.Error("Expected nil factory when JIRA_EMAIL is missing, got non-nil")
	}

	// Verify specific error message
	if tempooErr, ok := err.(*TempooError); ok {
		expectedMsg := "JIRA_EMAIL environment variable is not set"
		if tempooErr.Message != expectedMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedMsg, tempooErr.Message)
		}
	} else {
		t.Errorf("Expected TempooError, got %T", err)
	}
}

func TestNewTempooFactory_FailsWithMissingJiraAPIToken(t *testing.T) {
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

	// Set only email, missing API token
	os.Setenv("JIRA_EMAIL", "test@example.com")
	os.Unsetenv("JIRA_API_TOKEN")

	// Test factory creation
	factory, err := NewTempooFactory()

	// Assertions
	if err == nil {
		t.Error("Expected error when JIRA_API_TOKEN is missing, got nil")
	}

	if factory != nil {
		t.Error("Expected nil factory when JIRA_API_TOKEN is missing, got non-nil")
	}

	// Verify specific error message
	if tempooErr, ok := err.(*TempooError); ok {
		expectedMsg := "JIRA_API_TOKEN environment variable is not set"
		if tempooErr.Message != expectedMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedMsg, tempooErr.Message)
		}
	} else {
		t.Errorf("Expected TempooError, got %T", err)
	}
}

func TestTempooFactory_GetClient(t *testing.T) {
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

	// Set valid environment variables
	testEmail := "test@example.com"
	testToken := "test-api-token"
	os.Setenv("JIRA_EMAIL", testEmail)
	os.Setenv("JIRA_API_TOKEN", testToken)

	// Create factory
	factory, err := NewTempooFactory()
	if err != nil {
		t.Fatalf("Failed to create factory: %v", err)
	}

	// Test GetClient
	client := factory.GetClient()

	// Assertions
	if client == nil {
		t.Error("Expected non-nil client from GetClient()")
	}

	if reflect.TypeOf(client).String() != "*internal.Tempoo" {
		t.Errorf("Expected *internal.Tempoo, got %T", client)
	}

	// Verify it returns the same instance
	client2 := factory.GetClient()
	if client != client2 {
		t.Error("Expected GetClient() to return the same instance")
	}
}

func TestTempooFactory_GetClient_WithNilInstance(t *testing.T) {
	// Create factory with nil instance (simulate edge case)
	factory := &TempooFactory{instance: nil}

	// Test GetClient
	client := factory.GetClient()

	// Assertions
	if client != nil {
		t.Error("Expected nil client when factory instance is nil")
	}
}

func TestTempooFactory_MultipleFactories(t *testing.T) {
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

	// Set valid environment variables
	testEmail := "test@example.com"
	testToken := "test-api-token"
	os.Setenv("JIRA_EMAIL", testEmail)
	os.Setenv("JIRA_API_TOKEN", testToken)

	// Create multiple factories
	factory1, err1 := NewTempooFactory()
	factory2, err2 := NewTempooFactory()

	// Assertions
	if err1 != nil || err2 != nil {
		t.Errorf("Expected no errors, got %v, %v", err1, err2)
	}

	if factory1 == nil || factory2 == nil {
		t.Error("Expected both factories to be non-nil")
	}

	// Verify they are different factory instances
	if factory1 == factory2 {
		t.Error("Expected different factory instances")
	}

	// Verify they have different Tempoo instances
	client1 := factory1.GetClient()
	client2 := factory2.GetClient()

	if client1 == client2 {
		t.Error("Expected different Tempoo instances from different factories")
	}
}

func TestTempooFactory_StructFields(t *testing.T) {
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

	// Set valid environment variables
	testEmail := "test@example.com"
	testToken := "test-api-token"
	os.Setenv("JIRA_EMAIL", testEmail)
	os.Setenv("JIRA_API_TOKEN", testToken)

	// Create factory
	factory, err := NewTempooFactory()
	if err != nil {
		t.Fatalf("Failed to create factory: %v", err)
	}

	// Test struct fields using reflection
	factoryValue := reflect.ValueOf(factory).Elem()
	factoryType := factoryValue.Type()

	// Verify struct has expected fields
	if factoryType.NumField() != 1 {
		t.Errorf("Expected 1 field in TempooFactory, got %d", factoryType.NumField())
	}

	// Verify instance field
	instanceField := factoryValue.FieldByName("instance")
	if !instanceField.IsValid() {
		t.Error("Expected 'instance' field to be valid")
	}

	if instanceField.Type().String() != "*internal.Tempoo" {
		t.Errorf("Expected instance field to be *internal.Tempoo, got %s", instanceField.Type().String())
	}

	if instanceField.IsNil() {
		t.Error("Expected instance field to be non-nil")
	}
}

// Benchmark tests
func BenchmarkNewTempooFactory(b *testing.B) {
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

	// Set valid environment variables
	os.Setenv("JIRA_EMAIL", "test@example.com")
	os.Setenv("JIRA_API_TOKEN", "test-api-token")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		factory, err := NewTempooFactory()
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
		_ = factory
	}
}

func BenchmarkTempooFactory_GetClient(b *testing.B) {
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

	// Set valid environment variables
	os.Setenv("JIRA_EMAIL", "test@example.com")
	os.Setenv("JIRA_API_TOKEN", "test-api-token")

	// Create factory once
	factory, err := NewTempooFactory()
	if err != nil {
		b.Fatalf("Failed to create factory: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client := factory.GetClient()
		_ = client
	}
}

// Table-driven tests
func TestNewTempooFactory_TableDriven(t *testing.T) {
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

	tests := []struct {
		name        string
		email       string
		token       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid credentials",
			email:       "test@example.com",
			token:       "valid-token",
			expectError: false,
		},
		{
			name:        "missing email",
			email:       "",
			token:       "valid-token",
			expectError: true,
			errorMsg:    "JIRA_EMAIL environment variable is not set",
		},
		{
			name:        "missing token",
			email:       "test@example.com",
			token:       "",
			expectError: true,
			errorMsg:    "JIRA_API_TOKEN environment variable is not set",
		},
		{
			name:        "both missing",
			email:       "",
			token:       "",
			expectError: true,
			errorMsg:    "JIRA_EMAIL environment variable is not set",
		},
		{
			name:        "special characters in email",
			email:       "test+user@example-domain.com",
			token:       "valid-token",
			expectError: false,
		},
		{
			name:        "special characters in token",
			email:       "test@example.com",
			token:       "token-with-special-chars_123!@#$%",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			if tt.email != "" {
				os.Setenv("JIRA_EMAIL", tt.email)
			} else {
				os.Unsetenv("JIRA_EMAIL")
			}
			if tt.token != "" {
				os.Setenv("JIRA_API_TOKEN", tt.token)
			} else {
				os.Unsetenv("JIRA_API_TOKEN")
			}

			// Test factory creation
			factory, err := NewTempooFactory()

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if factory != nil {
					t.Error("Expected nil factory on error")
				}
				if tt.errorMsg != "" {
					if tempooErr, ok := err.(*TempooError); ok {
						if tempooErr.Message != tt.errorMsg {
							t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, tempooErr.Message)
						}
					} else {
						t.Errorf("Expected TempooError, got %T", err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if factory == nil {
					t.Error("Expected non-nil factory")
				}
				if factory != nil {
					client := factory.GetClient()
					if client == nil {
						t.Error("Expected non-nil client from factory")
					}
				}
			}
		})
	}
}
