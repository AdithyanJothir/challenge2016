package distributor

import (
	"CHALLENGE2016/internal/models/region"
	"fmt"
	"strings"
)

type Distributor struct {
	ID string

	Name string

	AuthorizedRegions map[string]*region.Region

	ExcludedRegions map[string]*region.Region

	ParentDistributors map[string]*Distributor

	Children map[string]*Distributor

	EffectiveAuthorizedRegions map[string]*region.Region
}

func NewDistributor(id, name string) *Distributor {
	return &Distributor{
		ID:                         id,
		Name:                       name,
		AuthorizedRegions:          make(map[string]*region.Region),
		ExcludedRegions:            make(map[string]*region.Region),
		ParentDistributors:         make(map[string]*Distributor),
		Children:                   make(map[string]*Distributor),
		EffectiveAuthorizedRegions: make(map[string]*region.Region),
	}
}

// AddDistributor creates a distributor with an initial set of authorized regions.
func AddDistributor(id, name string, authorizedRegions map[string]*region.Region) *Distributor {
	d := NewDistributor(id, name)
	for key, r := range authorizedRegions {
		d.AuthorizedRegions[key] = r
	}
	d.updateEffectiveAuthorizedRegions()
	return d
}

// allowedByAllParents checks if a candidate region is allowed by all parent distributors.
func allowedByAllParents(candidate *region.Region, parents map[string]*Distributor) bool {
	for _, parent := range parents {
		// Check parent's exclusions.
		for _, ex := range parent.ExcludedRegions {
			if isSubregion(candidate, ex) {
				return false
			}
		}
		// Check that candidate is allowed by at least one of parent's effective regions.
		allowed := false
		for _, pr := range parent.EffectiveAuthorizedRegions {
			if isSubregion(candidate, pr) {
				allowed = true
				break
			}
		}
		if !allowed {
			return false
		}
	}
	return true
}

// AddRegion adds a region to the distributor's AuthorizedRegions.
func (d *Distributor) AddRegion(r *region.Region) error {
	if !allowedByAllParents(r, d.ParentDistributors) {
		return fmt.Errorf("cannot add region %s: not authorized by all parent distributors", r.FullPath)
	}
	d.AuthorizedRegions[r.FullPath] = r
	d.updateEffectiveAuthorizedRegions()
	d.propagateEffectiveUpdate()
	return nil
}

func (d *Distributor) RemoveRegion(r *region.Region) {
	delete(d.AuthorizedRegions, r.FullPath)
	d.updateEffectiveAuthorizedRegions()
	d.propagateEffectiveUpdate()
}

func (d *Distributor) ExcludeRegion(r *region.Region) {
	d.ExcludedRegions[r.FullPath] = r
	d.updateEffectiveAuthorizedRegions()
	d.propagateEffectiveUpdate()
}

func (child *Distributor) AddParent(parent *Distributor) {
	if parent == nil {
		return
	}
	child.ParentDistributors[parent.ID] = parent
	parent.Children[child.ID] = child
	child.updateEffectiveAuthorizedRegions()
	child.propagateEffectiveUpdate()
}

// updateEffectiveAuthorizedRegions recomputes effective regions for the distributor.
func (d *Distributor) updateEffectiveAuthorizedRegions() {
	effective := make(map[string]*region.Region)
	if len(d.ParentDistributors) == 0 {
		effective = copyMap(d.AuthorizedRegions)
	} else {
		for _, r := range d.AuthorizedRegions {
			if allowedByAllParents(r, d.ParentDistributors) {
				effective[r.FullPath] = r
			}
		}
	}
	effective = subtract(effective, d.ExcludedRegions)
	d.EffectiveAuthorizedRegions = effective
}

// propagateEffectiveUpdate recursively updates effective regions for all child distributors.
func (d *Distributor) propagateEffectiveUpdate() {
	for _, child := range d.Children {
		child.updateEffectiveAuthorizedRegions()
		child.propagateEffectiveUpdate()
	}
}

func (d *Distributor) HasPermission(r *region.Region) bool {
	// Check if r is excluded at this distributor.
	for _, ex := range d.ExcludedRegions {
		if isSubregion(r, ex) {
			return false
		}
	}
	// Check if r is allowed in the effective regions.
	for _, er := range d.EffectiveAuthorizedRegions {
		if isSubregion(r, er) {
			return true
		}
	}
	return false
}

// isSubregion returns true if child is equal to or a subregion of parent.
// It splits the FullPath strings into tokens and checks that parent's tokens form a prefix of child's tokens.
func isSubregion(child, parent *region.Region) bool {
	if parent == nil {
		return false
	}
	if parent.FullPath == "" {
		return true
	}
	childParts := strings.Split(child.FullPath, "-")
	parentParts := strings.Split(parent.FullPath, "-")
	if len(childParts) < len(parentParts) {
		return false
	}
	for i, token := range parentParts {
		if childParts[i] != token {
			return false
		}
	}
	return true
}

func copyMap(src map[string]*region.Region) map[string]*region.Region {
	newMap := make(map[string]*region.Region)
	for k, v := range src {
		newMap[k] = v
	}
	return newMap
}

func subtract(a, b map[string]*region.Region) map[string]*region.Region {
	result := copyMap(a)
	for k := range b {
		delete(result, k)
	}
	return result
}
