package internal

import (
	"reflect"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
)

func TestJiraResponseType(t *testing.T) {
	// Test that JiraResponse is a map[string]interface{}
	var response JiraResponse
	response = make(map[string]interface{})

	// Test basic operations
	response["key"] = "value"
	response["number"] = 123
	response["boolean"] = true
	response["nested"] = map[string]interface{}{
		"inner": "value",
	}

	// Verify type
	if reflect.TypeOf(response).String() != "internal.JiraResponse" {
		t.Errorf("JiraResponse should be of type internal.JiraResponse, got %T", response)
	}

	// Verify underlying type
	if reflect.TypeOf(response).Kind() != reflect.Map {
		t.Errorf("JiraResponse should be a map type, got %v", reflect.TypeOf(response).Kind())
	}

	// Test accessing values
	if response["key"] != "value" {
		t.Errorf("Expected 'value', got %v", response["key"])
	}

	if response["number"] != 123 {
		t.Errorf("Expected 123, got %v", response["number"])
	}

	if response["boolean"] != true {
		t.Errorf("Expected true, got %v", response["boolean"])
	}
}

func TestWorklogDataType(t *testing.T) {
	// Test that WorklogData is a map[string]interface{}
	var worklog WorklogData
	worklog = make(map[string]interface{})

	// Test typical worklog data
	worklog["timeSpent"] = "2h"
	worklog["started"] = "2023-12-15T08:30:00.000+0000"
	worklog["author"] = map[string]interface{}{
		"accountId":   "user123",
		"displayName": "John Doe",
	}
	worklog["id"] = "12345"
	worklog["timeSpentSeconds"] = 7200

	// Verify type
	if reflect.TypeOf(worklog).String() != "internal.WorklogData" {
		t.Errorf("WorklogData should be of type internal.WorklogData, got %T", worklog)
	}

	// Verify underlying type
	if reflect.TypeOf(worklog).Kind() != reflect.Map {
		t.Errorf("WorklogData should be a map type, got %v", reflect.TypeOf(worklog).Kind())
	}

	// Test accessing values
	if worklog["timeSpent"] != "2h" {
		t.Errorf("Expected '2h', got %v", worklog["timeSpent"])
	}

	if worklog["timeSpentSeconds"] != 7200 {
		t.Errorf("Expected 7200, got %v", worklog["timeSpentSeconds"])
	}
}

func TestTempooStruct(t *testing.T) {
	// Test Tempoo struct creation and field access
	client := resty.New()
	tempoo := &Tempoo{
		email:    "test@example.com",
		apiToken: "test-token",
		client:   client,
	}

	// Test that struct fields are properly set (note: fields are unexported)
	// We can't directly access them, but we can test the struct exists and has the right type
	if reflect.TypeOf(tempoo).String() != "*internal.Tempoo" {
		t.Errorf("Tempoo should be of type *internal.Tempoo, got %T", tempoo)
	}

	// Test that the struct has the expected fields using reflection
	tempooType := reflect.TypeOf(tempoo).Elem()

	expectedFields := []struct {
		name string
		typ  string
	}{
		{"email", "string"},
		{"apiToken", "string"},
		{"client", "*resty.Client"},
	}

	if tempooType.NumField() != len(expectedFields) {
		t.Errorf("Expected %d fields, got %d", len(expectedFields), tempooType.NumField())
	}

	for i, expected := range expectedFields {
		field := tempooType.Field(i)
		if field.Name != expected.name {
			t.Errorf("Field %d: expected name %s, got %s", i, expected.name, field.Name)
		}
		if field.Type.String() != expected.typ {
			t.Errorf("Field %d (%s): expected type %s, got %s", i, expected.name, expected.typ, field.Type.String())
		}
	}
}

func TestTempooStructZeroValue(t *testing.T) {
	// Test zero value of Tempoo struct
	var tempoo Tempoo

	// Check that it's the zero value
	if reflect.ValueOf(tempoo).IsZero() != true {
		t.Error("Zero value Tempoo should be zero")
	}

	// Test pointer to zero value
	tempooPtr := &Tempoo{}

	// The pointer itself is not zero, but the struct it points to should be
	if reflect.ValueOf(*tempooPtr).IsZero() != true {
		t.Error("Zero value Tempoo struct should be zero")
	}
}

func TestTypeAliasCompatibility(t *testing.T) {
	// Test that type aliases are compatible with their underlying types

	// JiraResponse compatibility
	var jiraResp JiraResponse
	var mapResp map[string]interface{}

	jiraResp = make(JiraResponse)
	mapResp = jiraResp // Should be assignable

	if mapResp == nil {
		t.Error("JiraResponse should be assignable to map[string]interface{}")
	}

	// WorklogData compatibility
	var worklogData WorklogData
	var mapData map[string]interface{}

	worklogData = make(WorklogData)
	mapData = worklogData // Should be assignable

	if mapData == nil {
		t.Error("WorklogData should be assignable to map[string]interface{}")
	}
}

func TestJiraResponseOperations(t *testing.T) {
	// Test common operations on JiraResponse
	response := make(JiraResponse)

	// Test setting and getting values
	response["accountId"] = "user123"
	response["displayName"] = "John Doe"
	response["active"] = true
	response["worklogs"] = []interface{}{
		map[string]interface{}{
			"id":        "worklog1",
			"timeSpent": "2h",
		},
		map[string]interface{}{
			"id":        "worklog2",
			"timeSpent": "1h",
		},
	}

	// Test length
	if len(response) != 4 {
		t.Errorf("Expected 4 items in response, got %d", len(response))
	}

	// Test key existence
	if _, exists := response["accountId"]; !exists {
		t.Error("accountId key should exist")
	}

	if _, exists := response["nonexistent"]; exists {
		t.Error("nonexistent key should not exist")
	}

	// Test deletion
	delete(response, "active")
	if len(response) != 3 {
		t.Errorf("Expected 3 items after deletion, got %d", len(response))
	}

	// Test type assertions
	if accountId, ok := response["accountId"].(string); !ok || accountId != "user123" {
		t.Errorf("Expected accountId to be 'user123', got %v", accountId)
	}

	if worklogs, ok := response["worklogs"].([]interface{}); !ok || len(worklogs) != 2 {
		t.Errorf("Expected worklogs to be []interface{} with 2 items, got %v", worklogs)
	}
}

func TestWorklogDataOperations(t *testing.T) {
	// Test common operations on WorklogData
	worklog := make(WorklogData)

	// Test typical worklog structure
	worklog["id"] = "12345"
	worklog["timeSpent"] = "2h 30m"
	worklog["timeSpentSeconds"] = 9000
	worklog["started"] = "2023-12-15T08:30:00.000+0000"
	worklog["author"] = map[string]interface{}{
		"accountId":   "user123",
		"displayName": "John Doe",
		"active":      true,
	}
	worklog["comment"] = map[string]interface{}{
		"content": []interface{}{
			map[string]interface{}{
				"type": "paragraph",
				"content": []interface{}{
					map[string]interface{}{
						"type": "text",
						"text": "Working on feature",
					},
				},
			},
		},
	}

	// Test accessing nested data
	if author, ok := worklog["author"].(map[string]interface{}); ok {
		if displayName, ok := author["displayName"].(string); !ok || displayName != "John Doe" {
			t.Errorf("Expected author displayName to be 'John Doe', got %v", displayName)
		}
	} else {
		t.Error("Expected author to be a map[string]interface{}")
	}

	// Test numeric values
	if timeSpentSeconds, ok := worklog["timeSpentSeconds"].(int); !ok || timeSpentSeconds != 9000 {
		t.Errorf("Expected timeSpentSeconds to be 9000, got %v", timeSpentSeconds)
	}
}

func TestTempooStructMethods(t *testing.T) {
	// Test that Tempoo struct can have methods (by checking method set)
	tempoo := &Tempoo{}
	tempooType := reflect.TypeOf(tempoo)

	// Check that it has the expected methods
	expectedMethods := []string{
		"AddWorklog",
		"DeleteWorklog",
		"GetUserAccountID",
		"GetWorklogs",
		"ListWorklogs",
	}

	methodCount := tempooType.NumMethod()
	if methodCount < len(expectedMethods) {
		t.Errorf("Expected at least %d methods, got %d", len(expectedMethods), methodCount)
	}

	// Check for specific methods
	for _, methodName := range expectedMethods {
		if method, found := tempooType.MethodByName(methodName); !found {
			t.Errorf("Expected method %s not found", methodName)
		} else {
			// Verify it's actually a method (has at least one parameter - the receiver)
			if method.Type.NumIn() < 1 {
				t.Errorf("Method %s should have at least one parameter (receiver)", methodName)
			}
		}
	}
}

func TestTypeDefinitions(t *testing.T) {
	// Test that our types are properly defined
	tests := []struct {
		name         string
		typeInstance interface{}
		expectedType string
	}{
		{
			name:         "JiraResponse",
			typeInstance: JiraResponse{},
			expectedType: "internal.JiraResponse",
		},
		{
			name:         "WorklogData",
			typeInstance: WorklogData{},
			expectedType: "internal.WorklogData",
		},
		{
			name:         "Tempoo",
			typeInstance: Tempoo{},
			expectedType: "internal.Tempoo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualType := reflect.TypeOf(tt.typeInstance).String()
			if actualType != tt.expectedType {
				t.Errorf("Expected type %s, got %s", tt.expectedType, actualType)
			}
		})
	}
}

func TestTypeAliasUnderlyingTypes(t *testing.T) {
	// Test that type aliases have the correct underlying types
	var jiraResp JiraResponse
	var worklogData WorklogData

	// Test that they can be used as their underlying types
	jiraResp = make(map[string]interface{})
	worklogData = make(map[string]interface{})

	// Test that they support map operations
	jiraResp["test"] = "value"
	worklogData["test"] = "value"

	if jiraResp["test"] != "value" {
		t.Error("JiraResponse should support map operations")
	}

	if worklogData["test"] != "value" {
		t.Error("WorklogData should support map operations")
	}
}

// Benchmark tests
func BenchmarkJiraResponseCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		response := make(JiraResponse)
		response["key"] = "value"
		response["number"] = 123
		_ = response
	}
}

func BenchmarkWorklogDataCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		worklog := make(WorklogData)
		worklog["timeSpent"] = "2h"
		worklog["id"] = "12345"
		_ = worklog
	}
}

func BenchmarkTempooStructCreation(b *testing.B) {
	client := resty.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tempoo := &Tempoo{
			email:    "test@example.com",
			apiToken: "token",
			client:   client,
		}
		_ = tempoo
	}
}

func TestTempooStructWithRealClient(t *testing.T) {
	// Test Tempoo struct with actual resty client
	client := resty.New()
	client.SetTimeout(5 * time.Second)
	client.SetHeader("Content-Type", "application/json")

	tempoo := &Tempoo{
		email:    "test@example.com",
		apiToken: "test-token",
		client:   client,
	}

	// Test that the struct can be created without errors
	if tempoo == nil {
		t.Error("Tempoo struct should not be nil")
	}

	// Test using reflection to verify the client field is set
	tempooValue := reflect.ValueOf(tempoo).Elem()
	clientField := tempooValue.FieldByName("client")

	if !clientField.IsValid() {
		t.Error("client field should be valid")
	}

	if clientField.IsNil() {
		t.Error("client field should not be nil")
	}
}

func TestTypeConversions(t *testing.T) {
	// Test conversions between type aliases and their underlying types

	// JiraResponse to map[string]interface{}
	jiraResp := JiraResponse{"key": "value"}
	var mapResp map[string]interface{} = jiraResp

	if mapResp["key"] != "value" {
		t.Error("Conversion from JiraResponse to map should preserve data")
	}

	// map[string]interface{} to JiraResponse
	mapData := map[string]interface{}{"test": "data"}
	var jiraData JiraResponse = mapData

	if jiraData["test"] != "data" {
		t.Error("Conversion from map to JiraResponse should preserve data")
	}

	// WorklogData to map[string]interface{}
	worklogData := WorklogData{"timeSpent": "1h"}
	var worklogMap map[string]interface{} = worklogData

	if worklogMap["timeSpent"] != "1h" {
		t.Error("Conversion from WorklogData to map should preserve data")
	}

	// map[string]interface{} to WorklogData
	mapWorklog := map[string]interface{}{"id": "123"}
	var worklogType WorklogData = mapWorklog

	if worklogType["id"] != "123" {
		t.Error("Conversion from map to WorklogData should preserve data")
	}
}
