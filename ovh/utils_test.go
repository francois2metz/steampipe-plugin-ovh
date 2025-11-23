package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/ovh/go-ovh/ovh"
	"github.com/turbot/steampipe-plugin-sdk/v5/connection"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/context_key"
)

// loadFixtureRegions loads the list of regions from fixtures/ovh_regions.json
func loadFixtureRegions(t *testing.T) []string {
	data, err := os.ReadFile(filepath.Join("..", "fixtures", "ovh_regions.json"))
	if err != nil {
		t.Fatalf("Failed to read regions fixture: %v", err)
	}

	var regions []string
	if err := json.Unmarshal(data, &regions); err != nil {
		t.Fatalf("Failed to parse regions fixture: %v", err)
	}

	return regions
}

// loadFixtureRegionDetails loads region details from fixtures/ovh_region_*.json files
func loadFixtureRegionDetails(t *testing.T, regionNames []string) map[string]Region {
	details := make(map[string]Region)

	for _, regionName := range regionNames {
		filename := filepath.Join("..", "fixtures", fmt.Sprintf("ovh_region_%s.json", regionName))
		data, err := os.ReadFile(filename)
		if err != nil {
			t.Fatalf("Failed to read region fixture %s: %v", filename, err)
		}

		var region Region
		if err := json.Unmarshal(data, &region); err != nil {
			t.Fatalf("Failed to parse region fixture %s: %v", filename, err)
		}

		details[regionName] = region
	}

	return details
}

// mockOVHServer creates a test HTTP server that mocks OVH API responses
func mockOVHServer(t *testing.T, regions []string, regionDetails map[string]Region) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock authentication endpoint
		if r.URL.Path == "/auth/time" {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(1234567890); err != nil {
				t.Errorf("Failed to encode time response: %v", err)
			}
			return
		}

		// Mock region list endpoint
		if r.URL.Path == "/cloud/project/test-project/region" {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(regions); err != nil {
				t.Errorf("Failed to encode regions response: %v", err)
			}
			return
		}

		// Mock region details endpoint
		for regionName, regionInfo := range regionDetails {
			if r.URL.Path == fmt.Sprintf("/cloud/project/test-project/region/%s", regionName) {
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(regionInfo); err != nil {
					t.Errorf("Failed to encode region details response: %v", err)
				}
				return
			}
		}

		// Default 404 for unknown endpoints
		http.NotFound(w, r)
	}))
}

// createMockQueryData creates a mock QueryData with a mocked OVH client
func createMockQueryData(t *testing.T, client *ovh.Client, projectId string, configRegions []string) *plugin.QueryData {
	// Create a mock connection with config
	conn := &plugin.Connection{
		Name: "test",
		Config: ovhConfig{
			ApplicationKey:    stringPtr("test-key"),
			ApplicationSecret: stringPtr("test-secret"),
			ConsumerKey:       stringPtr("test-consumer"),
			Endpoint:          stringPtr("ovh-eu"),
			Regions:           configRegions,
		},
	}

	// Create connection cache
	connectionCache, err := connection.NewConnectionCache("test", 10000)
	if err != nil {
		t.Fatalf("Failed to create connection cache: %v", err)
	}

	// Create cache wrapper
	cache := connection.NewCache(connectionCache)

	// Cache the OVH client
	cache.Set("ovh", client)

	queryData := &plugin.QueryData{
		Connection:        conn,
		ConnectionManager: &connection.Manager{Cache: cache},
		ConnectionCache:   connectionCache,
		EqualsQuals:       map[string]*proto.QualValue{},
	}

	// Set project_id qual if provided
	if projectId != "" {
		queryData.EqualsQuals["project_id"] = &proto.QualValue{
			Value: &proto.QualValue_StringValue{StringValue: projectId},
		}
	}

	return queryData
}

// stringPtr returns a pointer to a string
func stringPtr(s string) *string {
	return &s
}

// TestMatchRegionPattern tests the matchRegionPattern function
func TestMatchRegionPattern(t *testing.T) {
	tests := []struct {
		name     string
		region   string
		pattern  string
		expected bool
	}{
		// Exact matches
		{
			name:     "exact match",
			region:   "GRA9",
			pattern:  "GRA9",
			expected: true,
		},
		{
			name:     "exact match case sensitive",
			region:   "GRA9",
			pattern:  "gra9",
			expected: false,
		},
		{
			name:     "no match",
			region:   "GRA9",
			pattern:  "SBG5",
			expected: false,
		},
		// Wildcard * matches
		{
			name:     "wildcard all",
			region:   "GRA9",
			pattern:  "*",
			expected: true,
		},
		{
			name:     "prefix wildcard",
			region:   "GRA9",
			pattern:  "GRA*",
			expected: true,
		},
		{
			name:     "prefix wildcard no match",
			region:   "SBG5",
			pattern:  "GRA*",
			expected: false,
		},
		{
			name:     "suffix wildcard",
			region:   "GRA9",
			pattern:  "*9",
			expected: true,
		},
		{
			name:     "suffix wildcard no match",
			region:   "GRA9",
			pattern:  "*5",
			expected: false,
		},
		{
			name:     "middle wildcard",
			region:   "GRA9",
			pattern:  "G*9",
			expected: true,
		},
		// Edge cases
		{
			name:     "empty pattern",
			region:   "GRA9",
			pattern:  "",
			expected: false,
		},
		{
			name:     "empty region",
			region:   "",
			pattern:  "*",
			expected: true,
		},
		{
			name:     "pattern longer than region",
			region:   "GRA",
			pattern:  "GRA7",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchRegionPattern(tt.region, tt.pattern)
			if result != tt.expected {
				t.Errorf("matchRegionPattern(%q, %q) = %v, want %v", tt.region, tt.pattern, result, tt.expected)
			}
		})
	}
}

// TestRegionMatrix tests the RegionMatrix function with various scenarios
func TestRegionMatrix(t *testing.T) {
	// Load fixture data from JSON files
	mockRegions := loadFixtureRegions(t)
	mockRegionDetails := loadFixtureRegionDetails(t, mockRegions)

	tests := []struct {
		name            string
		serviceNames    []string
		configRegions   []string
		projectId       string
		expectedCount   int
		expectedRegions []string
		description     string
	}{
		{
			name:            "no service names - all regions",
			serviceNames:    []string{},
			configRegions:   []string{"*"},
			projectId:       "test-project",
			expectedCount:   17,
			expectedRegions: []string{"BHS", "BHS5", "CA-EAST-TOR", "DE", "DE1", "EU-WEST-PAR", "GRA", "GRA9", "RBX", "RBX-A", "RBX-ARCHIVE", "SBG", "SBG5", "UK", "UK1", "WAW", "WAW1"},
			description:     "Should return all regions when no service filter is specified",
		},
		{
			name:            "filter by S3 standard storage",
			serviceNames:    []string{"storage-s3-standard"},
			configRegions:   []string{"*"},
			projectId:       "test-project",
			expectedCount:   9,
			expectedRegions: []string{"BHS", "CA-EAST-TOR", "DE", "EU-WEST-PAR", "GRA", "RBX", "SBG", "UK", "WAW"},
			description:     "Should return only regions with S3 standard storage",
		},
		{
			name:            "filter by multiple S3 services",
			serviceNames:    []string{"storage-s3-standard", "storage-s3-high-perf"},
			configRegions:   []string{"*"},
			projectId:       "test-project",
			expectedCount:   9,
			expectedRegions: []string{"BHS", "CA-EAST-TOR", "DE", "EU-WEST-PAR", "GRA", "RBX", "SBG", "UK", "WAW"},
			description:     "Should return regions with any of the specified services",
		},
		{
			name:            "filter by coldarchive",
			serviceNames:    []string{"storage-s3-coldarchive"},
			configRegions:   []string{"*"},
			projectId:       "test-project",
			expectedCount:   1,
			expectedRegions: []string{"RBX-ARCHIVE"},
			description:     "Should return only regions with coldarchive service",
		},
		{
			name:            "wildcard pattern GRA*",
			serviceNames:    []string{},
			configRegions:   []string{"GRA*"},
			projectId:       "test-project",
			expectedCount:   2,
			expectedRegions: []string{"GRA", "GRA9"},
			description:     "Should match regions starting with GRA",
		},
		{
			name:            "wildcard pattern with service filter",
			serviceNames:    []string{"instance"},
			configRegions:   []string{"*5"},
			projectId:       "test-project",
			expectedCount:   2,
			expectedRegions: []string{"SBG5", "BHS5"},
			description:     "Should match regions ending with 5 that have instance service",
		},
		{
			name:            "specific regions in config",
			serviceNames:    []string{},
			configRegions:   []string{"GRA", "SBG5", "UK1"},
			projectId:       "test-project",
			expectedCount:   3,
			expectedRegions: []string{"GRA", "SBG5", "UK1"},
			description:     "Should return only specified regions",
		},
		{
			name:            "multiple wildcard patterns",
			serviceNames:    []string{},
			configRegions:   []string{"GRA*", "SBG*"},
			projectId:       "test-project",
			expectedCount:   4,
			expectedRegions: []string{"GRA", "GRA9", "SBG", "SBG5"},
			description:     "Should match multiple wildcard patterns",
		},
		{
			name:            "no matching regions",
			serviceNames:    []string{"non-existent-service"},
			configRegions:   []string{"*"},
			projectId:       "test-project",
			expectedCount:   0,
			expectedRegions: []string{},
			description:     "Should return empty when no regions have the service",
		},
		{
			name:            "no project_id",
			serviceNames:    []string{},
			configRegions:   []string{"*"},
			projectId:       "",
			expectedCount:   0,
			expectedRegions: []string{},
			description:     "Should return empty matrix when project_id is missing",
		},
		{
			name:            "service filter with wildcard pattern RBX",
			serviceNames:    []string{"storage-s3-coldarchive"},
			configRegions:   []string{"RBX*"},
			projectId:       "test-project",
			expectedCount:   1,
			expectedRegions: []string{"RBX-ARCHIVE"},
			description:     "Should combine service filtering with pattern matching",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := mockOVHServer(t, mockRegions, mockRegionDetails)
			defer server.Close()

			// Create OVH client pointing to mock server
			client, err := ovh.NewClient(
				server.URL,
				"test-key",
				"test-secret",
				"test-consumer",
			)
			if err != nil {
				t.Fatalf("Failed to create OVH client: %v", err)
			}

			// Create mock query data
			queryData := createMockQueryData(t, client, tt.projectId, tt.configRegions)

			// Create context with logger
			logger := hclog.New(&hclog.LoggerOptions{
				Name:   "test",
				Level:  hclog.Error,
				Output: io.Discard,
			})
			ctx := context.WithValue(context.Background(), context_key.Logger, logger)

			// Call RegionMatrix
			matrixFunc := RegionMatrix(tt.serviceNames...)
			result := matrixFunc(ctx, queryData)

			// Verify count
			if len(result) != tt.expectedCount {
				t.Errorf("Expected %d regions, got %d. Description: %s", tt.expectedCount, len(result), tt.description)
			}

			// Verify expected regions are present
			resultRegions := make(map[string]bool)
			for _, item := range result {
				if region, ok := item["region"].(string); ok {
					resultRegions[region] = true
				}
			}

			for _, expectedRegion := range tt.expectedRegions {
				if !resultRegions[expectedRegion] {
					t.Errorf("Expected region %q not found in result. Description: %s", expectedRegion, tt.description)
				}
			}

			// Verify no unexpected regions
			if len(resultRegions) != len(tt.expectedRegions) {
				t.Errorf("Result contains unexpected regions. Expected: %v, Got: %v", tt.expectedRegions, resultRegions)
			}
		})
	}
}

// TestRegionMatrixCaching tests that the RegionMatrix function properly caches API responses
func TestRegionMatrixCaching(t *testing.T) {
	mockRegions := []string{"GRA9", "SBG5"}
	mockRegionDetails := map[string]Region{
		"GRA9": {
			Name:     "GRA9",
			Services: []Component{{Name: "storage-s3-standard", Status: "UP"}},
			Status:   "UP",
		},
		"SBG5": {
			Name:     "SBG5",
			Services: []Component{{Name: "storage-s3-standard", Status: "UP"}},
			Status:   "UP",
		},
	}

	// Track API calls
	apiCalls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/auth/time" {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(1234567890); err != nil {
				t.Errorf("Failed to encode time response: %v", err)
			}
			return
		}

		if r.URL.Path == "/cloud/project/test-project/region" {
			apiCalls++
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(mockRegions); err != nil {
				t.Errorf("Failed to encode regions response: %v", err)
			}
			return
		}

		for regionName, regionInfo := range mockRegionDetails {
			if r.URL.Path == fmt.Sprintf("/cloud/project/test-project/region/%s", regionName) {
				apiCalls++
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(regionInfo); err != nil {
					t.Errorf("Failed to encode region details response: %v", err)
				}
				return
			}
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	client, err := ovh.NewClient(server.URL, "test-key", "test-secret", "test-consumer")
	if err != nil {
		t.Fatalf("Failed to create OVH client: %v", err)
	}

	queryData := createMockQueryData(t, client, "test-project", []string{"*"})

	// Create context with logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "test",
		Level:  hclog.Error,
		Output: io.Discard,
	})
	ctx := context.WithValue(context.Background(), context_key.Logger, logger)

	// First call - should hit API
	matrixFunc := RegionMatrix("storage-s3-standard")
	result1 := matrixFunc(ctx, queryData)

	firstCallCount := apiCalls

	// Second call - should use cache
	result2 := matrixFunc(ctx, queryData)

	secondCallCount := apiCalls

	// Verify results are the same
	if len(result1) != len(result2) {
		t.Errorf("Results differ between calls: first=%d, second=%d", len(result1), len(result2))
	}

	// Verify caching worked (no additional API calls)
	if secondCallCount != firstCallCount {
		t.Errorf("Expected caching to prevent additional API calls. First call: %d API calls, Second call: %d API calls", firstCallCount, secondCallCount)
	}

	// Verify we got the expected regions
	if len(result1) != 2 {
		t.Errorf("Expected 2 regions, got %d", len(result1))
	}
}

// TestRegionMatrixEmptyRegionsList tests behavior when API returns empty regions list
func TestRegionMatrixEmptyRegionsList(t *testing.T) {
	server := mockOVHServer(t, []string{}, map[string]Region{})
	defer server.Close()

	client, err := ovh.NewClient(server.URL, "test-key", "test-secret", "test-consumer")
	if err != nil {
		t.Fatalf("Failed to create OVH client: %v", err)
	}

	queryData := createMockQueryData(t, client, "test-project", []string{"*"})

	// Create context with logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "test",
		Level:  hclog.Error,
		Output: io.Discard,
	})
	ctx := context.WithValue(context.Background(), context_key.Logger, logger)

	matrixFunc := RegionMatrix()
	result := matrixFunc(ctx, queryData)

	if len(result) != 0 {
		t.Errorf("Expected empty result for empty regions list, got %d regions", len(result))
	}
}

// TestRegionMatrixWithoutServiceFilter tests that all regions are returned when no service filter is applied
func TestRegionMatrixWithoutServiceFilter(t *testing.T) {
	mockRegions := []string{"GRA9", "SBG5", "UK1"}
	mockRegionDetails := map[string]Region{
		"GRA9": {
			Name:     "GRA9",
			Services: []Component{{Name: "storage-s3-standard", Status: "UP"}},
			Status:   "UP",
		},
		"SBG5": {
			Name:     "SBG5",
			Services: []Component{{Name: "instance", Status: "UP"}},
			Status:   "UP",
		},
		"UK1": {
			Name:     "UK1",
			Services: []Component{}, // No services
			Status:   "UP",
		},
	}

	server := mockOVHServer(t, mockRegions, mockRegionDetails)
	defer server.Close()

	client, err := ovh.NewClient(server.URL, "test-key", "test-secret", "test-consumer")
	if err != nil {
		t.Fatalf("Failed to create OVH client: %v", err)
	}

	queryData := createMockQueryData(t, client, "test-project", []string{"*"})

	// Create context with logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "test",
		Level:  hclog.Error,
		Output: io.Discard,
	})
	ctx := context.WithValue(context.Background(), context_key.Logger, logger)

	// Call without service filter
	matrixFunc := RegionMatrix()
	result := matrixFunc(ctx, queryData)

	// Should return all regions including those without services
	if len(result) != 3 {
		t.Errorf("Expected 3 regions (including region without services), got %d", len(result))
	}
}

// TestRegionMatrixSpecificRegionsNoAPICall tests that specific regions don't trigger API calls
func TestRegionMatrixSpecificRegionsNoAPICall(t *testing.T) {
	apiCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/auth/time" {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(1234567890); err != nil {
				t.Errorf("Failed to encode time response: %v", err)
			}
			return
		}
		// Any other call means API was called when it shouldn't be
		apiCalled = true
		http.NotFound(w, r)
	}))
	defer server.Close()

	client, err := ovh.NewClient(server.URL, "test-key", "test-secret", "test-consumer")
	if err != nil {
		t.Fatalf("Failed to create OVH client: %v", err)
	}

	// Use specific regions (no wildcards)
	queryData := createMockQueryData(t, client, "test-project", []string{"GRA9", "SBG5"})

	// Create context with logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "test",
		Level:  hclog.Error,
		Output: io.Discard,
	})
	ctx := context.WithValue(context.Background(), context_key.Logger, logger)

	matrixFunc := RegionMatrix()
	result := matrixFunc(ctx, queryData)

	if apiCalled {
		t.Error("API should not be called when using specific regions without wildcards")
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 regions, got %d", len(result))
	}
}
