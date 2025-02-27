package region

import (
	"strings"
)

type Region struct {
	Name     string
	Parent   *Region
	Children map[string]*Region
	FullPath string // Pre-computed full path, e.g., "USA-California-LosAngeles"
}

type RegionManager struct {
	Root      *Region
	regionMap map[string]*Region
}

func NewRegionManager() *RegionManager {
	root := &Region{
		Name:     "",
		Children: make(map[string]*Region),
		FullPath: "", // Root has an empty path
	}
	return &RegionManager{
		Root:      root,
		regionMap: make(map[string]*Region),
	}
}

// AddRegion adds a region given a full path (e.g., "USA-California-LosAngeles")
func (rm *RegionManager) AddRegion(path string) *Region {
	parts := strings.Split(path, "-")
	current := rm.Root
	fullPath := ""

	for _, part := range parts {
		if fullPath == "" {
			fullPath = part
		} else {
			fullPath = fullPath + "-" + part
		}

		if child, exists := current.Children[part]; exists {
			current = child
		} else {
			newRegion := &Region{
				Name:     part,
				Parent:   current,
				Children: make(map[string]*Region),
				FullPath: fullPath,
			}
			current.Children[part] = newRegion
			current = newRegion

			rm.regionMap[fullPath] = newRegion
		}
	}

	return current
}

// GetRegion retrieves a region by its full path in O(1) time i wanted this to be read optimized hence the ugly code :P
func (rm *RegionManager) GetRegion(path string) (*Region, bool) {
	region, exists := rm.regionMap[path]
	return region, exists
}
