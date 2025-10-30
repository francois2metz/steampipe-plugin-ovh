package ovh

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ovh/go-ovh/ovh"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Region represents a cloud region with its services
type Region struct {
	Name               string      `json:"name"`
	ContinentCode      string      `json:"continentCode"`
	DatacenterLocation string      `json:"datacenterLocation"`
	IpCountries        []string    `json:"ipCountries"`
	Services           []Component `json:"services"`
	Status             string      `json:"status"`
	Type               string      `json:"type"`
}

// Component represents a service component in a region
type Component struct {
	Endpoint string `json:"endpoint"`
	Name     string `json:"name"`
	Status   string `json:"status"`
}

func connect(ctx context.Context, d *plugin.QueryData) (*ovh.Client, error) {
	// get ovh client from cache
	cacheKey := "ovh"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*ovh.Client), nil
	}

	applicationKey := ""
	applicationSecret := ""
	consumerKey := ""
	endpoint := ""

	ovhConfig := GetConfig(d.Connection)

	if ovhConfig.ApplicationKey != nil {
		applicationKey = *ovhConfig.ApplicationKey
	}
	if ovhConfig.ApplicationSecret != nil {
		applicationSecret = *ovhConfig.ApplicationSecret
	}
	if ovhConfig.ConsumerKey != nil {
		consumerKey = *ovhConfig.ConsumerKey
	}
	if ovhConfig.Endpoint != nil {
		endpoint = *ovhConfig.Endpoint
	}

	if applicationKey == "" {
		return nil, errors.New("'application_key' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}
	if applicationSecret == "" {
		return nil, errors.New("'application_secret' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}
	if consumerKey == "" {
		return nil, errors.New("'consumer_key' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}
	if endpoint == "" {
		return nil, errors.New("'endpoint' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}

	client, _ := ovh.NewClient(
		endpoint,
		applicationKey,
		applicationSecret,
		consumerKey,
	)

	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, client)

	return client, nil
}

// RegionMatrix returns a function that generates a matrix of regions to query for OVH cloud resources.
// It accepts service names as parameters to filter regions based on available services.
// The returned function respects the regions configuration in the connection config and queries
// all available regions from the OVH API if no regions are configured or if wildcards are used.
//
// Example usage:
//   - For S3 storage: RegionMatrix("storage-s3-standard", "storage-s3-high-perf")
//   - For Archive storage: RegionMatrix("storage-s3-coldarchive")
func RegionMatrix(serviceNames ...string) func(context.Context, *plugin.QueryData) []map[string]interface{} {
	return func(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {
		// Get the OVH config
		ovhConfig := GetConfig(d.Connection)

		// Get project_id from quals if available
		projectId := ""
		if d.EqualsQuals["project_id"] != nil {
			projectId = d.EqualsQuals["project_id"].GetStringValue()
		}

		// If no project_id is available, we can't query regions
		// Return empty matrix - the table will need to handle this
		if projectId == "" {
			plugin.Logger(ctx).Warn("ovh.RegionMatrix", "no_project_id", "project_id is required to query regions")
			return []map[string]interface{}{}
		}

		// Get configured regions from config
		configuredRegions := ovhConfig.Regions
		if len(configuredRegions) == 0 {
			configuredRegions = []string{"*"}
		}

		// Check if we need to query all regions (wildcard or no specific regions)
		needsAllRegions := false
		for _, region := range configuredRegions {
			if region == "*" || strings.Contains(region, "*") || strings.Contains(region, "?") {
				needsAllRegions = true
				break
			}
		}

		var regions []string

		if needsAllRegions {
			// Query all available regions from OVH API
			client, err := connect(ctx, d)
			if err != nil {
				plugin.Logger(ctx).Error("ovh.RegionMatrix", "connection_error", err)
				return []map[string]interface{}{}
			}

			// Check cache for region names list
			regionsCacheKey := fmt.Sprintf("ovh_regions_%s", projectId)
			var regionNames []string

			if cachedData, ok := d.ConnectionManager.Cache.Get(regionsCacheKey); ok {
				regionNames = cachedData.([]string)
				plugin.Logger(ctx).Debug("ovh.RegionMatrix", "cache_hit", "regions_list", "project_id", projectId)
			} else {
				err = client.Get(fmt.Sprintf("/cloud/project/%s/region", projectId), &regionNames)
				if err != nil {
					plugin.Logger(ctx).Error("ovh.RegionMatrix", "api_error", err)
					return []map[string]interface{}{}
				}
				// Cache the region names list
				d.ConnectionManager.Cache.Set(regionsCacheKey, regionNames)
				plugin.Logger(ctx).Debug("ovh.RegionMatrix", "cache_set", "regions_list", "project_id", projectId, "count", len(regionNames))
			}

			// Filter regions to only include those with the specified services
			var filteredRegions []string
			for _, regionName := range regionNames {
				// If no service names specified, include all regions
				if len(serviceNames) == 0 {
					filteredRegions = append(filteredRegions, regionName)
					continue
				}

				// Fetch detailed region info to check for services
				regionDetailsCacheKey := fmt.Sprintf("ovh_region_details_%s_%s", projectId, regionName)
				var regionInfo Region

				if cachedData, ok := d.ConnectionManager.Cache.Get(regionDetailsCacheKey); ok {
					regionInfo = cachedData.(Region)
					plugin.Logger(ctx).Debug("ovh.RegionMatrix", "cache_hit", "region_details", "project_id", projectId, "region", regionName)
				} else {
					err = client.Get(fmt.Sprintf("/cloud/project/%s/region/%s", projectId, regionName), &regionInfo)
					if err != nil {
						plugin.Logger(ctx).Warn("ovh.RegionMatrix", "region_info_error", err, "region", regionName)
						continue
					}
					// Cache the region details
					d.ConnectionManager.Cache.Set(regionDetailsCacheKey, regionInfo)
					plugin.Logger(ctx).Debug("ovh.RegionMatrix", "cache_set", "region_details", "project_id", projectId, "region", regionName)
				}

				// Check if region has any of the specified services
				hasService := false
				for _, service := range regionInfo.Services {
					for _, serviceName := range serviceNames {
						if service.Name == serviceName {
							hasService = true
							break
						}
					}
					if hasService {
						break
					}
				}

				// Only include regions with at least one of the specified services
				if hasService {
					filteredRegions = append(filteredRegions, regionName)
				}
			}

			// If wildcard is "*", use all filtered regions
			if len(configuredRegions) == 1 && configuredRegions[0] == "*" {
				regions = filteredRegions
			} else {
				// Match filtered regions against patterns
				for _, regionName := range filteredRegions {
					for _, pattern := range configuredRegions {
						if matchRegionPattern(regionName, pattern) {
							regions = append(regions, regionName)
							break
						}
					}
				}
			}
		} else {
			// Use configured regions directly
			regions = configuredRegions
		}

		// Build matrix with regions in their original case
		matrix := make([]map[string]interface{}, 0, len(regions))
		for _, region := range regions {
			matrix = append(matrix, map[string]interface{}{
				"region": region,
			})
		}

		plugin.Logger(ctx).Debug("ovh.RegionMatrix", "project_id", projectId, "configured_regions", configuredRegions, "resolved_regions", regions, "matrix_size", len(matrix))

		return matrix
	}
}

// matchRegionPattern checks if a region name matches a pattern with wildcards
// Supports * (any characters) and ? (single character)
func matchRegionPattern(region, pattern string) bool {
	// Simple wildcard matching - case-sensitive to preserve OVH region names
	// If no wildcards, do exact match
	if !strings.Contains(pattern, "*") && !strings.Contains(pattern, "?") {
		return region == pattern
	}

	// Convert pattern to regex-like matching
	// For simplicity, we'll handle basic cases:
	// - "GRA*" matches "GRA1", "GRA7", etc.
	// - "*" matches everything
	if pattern == "*" {
		return true
	}

	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(region, prefix)
	}

	if strings.HasPrefix(pattern, "*") {
		suffix := strings.TrimPrefix(pattern, "*")
		return strings.HasSuffix(region, suffix)
	}

	// For more complex patterns, do a simple contains check
	// This is a simplified implementation
	patternParts := strings.Split(pattern, "*")
	pos := 0
	for i, part := range patternParts {
		if part == "" {
			continue
		}
		idx := strings.Index(region[pos:], part)
		if idx == -1 {
			return false
		}
		if i == 0 && idx != 0 {
			return false
		}
		pos += idx + len(part)
	}

	return true
}
