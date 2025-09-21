package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRegionList(t *testing.T) {
	regions := GetRegionList()

	// Basic validation
	assert.NotEmpty(t, regions)
	assert.Greater(t, len(regions), 20) // Ensure sufficient number of regions

	// Verify specific regions are included
	expectedRegions := []string{
		"asia-northeast1", // Tokyo, Japan
		"asia-northeast2", // Osaka, Japan
		"us-central1",     // Council Bluffs, Iowa, USA
		"europe-west1",    // St. Ghislain, Belgium
	}

	for _, expected := range expectedRegions {
		assert.Contains(t, regions, expected, "Region %s should be in the list", expected)
	}

	// Verify no duplicates
	regionMap := make(map[string]bool)
	for _, region := range regions {
		assert.False(t, regionMap[region], "Region %s should not be duplicated", region)
		regionMap[region] = true
	}

	// Verify all regions have valid format
	for _, region := range regions {
		assert.NotEmpty(t, region, "Region should not be empty")
		assert.Greater(t, len(region), 5, "Region %s should have reasonable length", region)
	}
}

func TestGetRegionList_Consistency(t *testing.T) {
	// Verify that multiple calls return the same result
	regions1 := GetRegionList()
	regions2 := GetRegionList()

	assert.Equal(t, regions1, regions2, "GetRegionList should return consistent results")
}

func TestGetRegionList_Order(t *testing.T) {
	regions := GetRegionList()

	// Verify the first few regions (ensure order doesn't change)
	expectedFirstRegions := []string{
		"asia-east1",
		"asia-east2",
		"asia-northeast1",
	}

	for i, expected := range expectedFirstRegions {
		if i < len(regions) {
			assert.Equal(t, expected, regions[i], "Region at index %d should be %s", i, expected)
		}
	}
}
