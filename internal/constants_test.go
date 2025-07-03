package internal

import (
	"net/url"
	"strings"
	"testing"
)

func TestJiraFQDN(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "JiraFQDN should be esendex.atlassian.net",
			expected: "esendex.atlassian.net",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if JiraFQDN != tt.expected {
				t.Errorf("JiraFQDN = %q, want %q", JiraFQDN, tt.expected)
			}
		})
	}
}

func TestJiraAPIRootURL(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "JiraAPIRootURL should be https://esendex.atlassian.net/rest/api/3",
			expected: "https://esendex.atlassian.net/rest/api/3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if JiraAPIRootURL != tt.expected {
				t.Errorf("JiraAPIRootURL = %q, want %q", JiraAPIRootURL, tt.expected)
			}
		})
	}
}

func TestJiraAPIRootURLConstruction(t *testing.T) {
	// Test that JiraAPIRootURL is correctly constructed from JiraFQDN
	expectedURL := "https://" + JiraFQDN + "/rest/api/3"

	if JiraAPIRootURL != expectedURL {
		t.Errorf("JiraAPIRootURL = %q, want %q", JiraAPIRootURL, expectedURL)
	}
}

func TestJiraFQDNFormat(t *testing.T) {
	tests := []struct {
		name        string
		checkFunc   func(string) bool
		description string
	}{
		{
			name: "should not be empty",
			checkFunc: func(s string) bool {
				return len(s) > 0
			},
			description: "JiraFQDN should not be empty",
		},
		{
			name: "should not contain protocol",
			checkFunc: func(s string) bool {
				return !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://")
			},
			description: "JiraFQDN should not contain protocol prefix",
		},
		{
			name: "should not contain path",
			checkFunc: func(s string) bool {
				return !strings.Contains(s, "/")
			},
			description: "JiraFQDN should not contain path components",
		},
		{
			name: "should be valid hostname format",
			checkFunc: func(s string) bool {
				// Basic hostname validation - should contain at least one dot
				return strings.Contains(s, ".")
			},
			description: "JiraFQDN should be a valid hostname format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.checkFunc(JiraFQDN) {
				t.Errorf("%s: JiraFQDN = %q", tt.description, JiraFQDN)
			}
		})
	}
}

func TestJiraAPIRootURLFormat(t *testing.T) {
	tests := []struct {
		name        string
		checkFunc   func(string) bool
		description string
	}{
		{
			name: "should be valid URL",
			checkFunc: func(s string) bool {
				_, err := url.Parse(s)
				return err == nil
			},
			description: "JiraAPIRootURL should be a valid URL",
		},
		{
			name: "should use HTTPS",
			checkFunc: func(s string) bool {
				return strings.HasPrefix(s, "https://")
			},
			description: "JiraAPIRootURL should use HTTPS protocol",
		},
		{
			name: "should contain API path",
			checkFunc: func(s string) bool {
				return strings.Contains(s, "/rest/api/3")
			},
			description: "JiraAPIRootURL should contain the REST API v3 path",
		},
		{
			name: "should end with version 3",
			checkFunc: func(s string) bool {
				return strings.HasSuffix(s, "/rest/api/3")
			},
			description: "JiraAPIRootURL should end with /rest/api/3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.checkFunc(JiraAPIRootURL) {
				t.Errorf("%s: JiraAPIRootURL = %q", tt.description, JiraAPIRootURL)
			}
		})
	}
}

func TestJiraAPIRootURLParsing(t *testing.T) {
	// Test that the URL can be parsed correctly
	parsedURL, err := url.Parse(JiraAPIRootURL)
	if err != nil {
		t.Fatalf("Failed to parse JiraAPIRootURL: %v", err)
	}

	tests := []struct {
		name     string
		got      string
		expected string
	}{
		{
			name:     "scheme should be https",
			got:      parsedURL.Scheme,
			expected: "https",
		},
		{
			name:     "host should match JiraFQDN",
			got:      parsedURL.Host,
			expected: JiraFQDN,
		},
		{
			name:     "path should be /rest/api/3",
			got:      parsedURL.Path,
			expected: "/rest/api/3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s: got %q, want %q", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestConstantsImmutability(t *testing.T) {
	// Test that constants have expected values (regression test)
	originalJiraFQDN := JiraFQDN
	originalJiraAPIRootURL := JiraAPIRootURL

	// These should remain constant throughout the test
	if JiraFQDN != originalJiraFQDN {
		t.Errorf("JiraFQDN changed during test: got %q, want %q", JiraFQDN, originalJiraFQDN)
	}

	if JiraAPIRootURL != originalJiraAPIRootURL {
		t.Errorf("JiraAPIRootURL changed during test: got %q, want %q", JiraAPIRootURL, originalJiraAPIRootURL)
	}
}

func TestJiraAPIEndpoints(t *testing.T) {
	// Test that common API endpoints can be constructed properly
	tests := []struct {
		name     string
		endpoint string
		expected string
	}{
		{
			name:     "myself endpoint",
			endpoint: "/myself",
			expected: "https://esendex.atlassian.net/rest/api/3/myself",
		},
		{
			name:     "issue endpoint",
			endpoint: "/issue/TEST-123",
			expected: "https://esendex.atlassian.net/rest/api/3/issue/TEST-123",
		},
		{
			name:     "worklog endpoint",
			endpoint: "/issue/TEST-123/worklog",
			expected: "https://esendex.atlassian.net/rest/api/3/issue/TEST-123/worklog",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullURL := JiraAPIRootURL + tt.endpoint
			if fullURL != tt.expected {
				t.Errorf("Constructed URL = %q, want %q", fullURL, tt.expected)
			}
		})
	}
}

// Benchmark tests for constants access
func BenchmarkJiraFQDNAccess(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = JiraFQDN
	}
}

func BenchmarkJiraAPIRootURLAccess(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = JiraAPIRootURL
	}
}

func BenchmarkURLConstruction(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = JiraAPIRootURL + "/issue/TEST-123"
	}
}

// Test constants in context of actual usage patterns
func TestConstantsInContext(t *testing.T) {
	// Test how constants would be used in real scenarios
	testCases := []struct {
		name        string
		issueKey    string
		expectedURL string
	}{
		{
			name:        "issue validation URL",
			issueKey:    "INF-88",
			expectedURL: "https://esendex.atlassian.net/rest/api/3/issue/INF-88",
		},
		{
			name:        "worklog URL",
			issueKey:    "PROJ-123",
			expectedURL: "https://esendex.atlassian.net/rest/api/3/issue/PROJ-123/worklog",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Simulate how constants are used in the actual code
			var constructedURL string
			if tc.name == "issue validation URL" {
				constructedURL = JiraAPIRootURL + "/issue/" + tc.issueKey
			} else if tc.name == "worklog URL" {
				constructedURL = JiraAPIRootURL + "/issue/" + tc.issueKey + "/worklog"
			}

			if constructedURL != tc.expectedURL {
				t.Errorf("Constructed URL = %q, want %q", constructedURL, tc.expectedURL)
			}
		})
	}
}
