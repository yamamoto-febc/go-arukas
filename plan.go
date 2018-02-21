package arukas

import "fmt"

const (
	// RegionJPTokyo represents the "jp-tokyo" region
	RegionJPTokyo = "jp-tokyo"
)

const (
	// PlanFree represents the "free" region
	PlanFree = "free"
	// PlanHobby represents the "hobby" region
	PlanHobby = "hobby"
	// PlanStandard1 represents the "standard-1" region
	PlanStandard1 = "standard-1"
	// PlanStandard2 represents the "standard-2" region
	PlanStandard2 = "standard-2"
)

var (
	// ValidRegions is a list of valid regions
	ValidRegions = []string{RegionJPTokyo}
	// ValidPlans is a list of valid
	ValidPlans = []string{PlanFree, PlanHobby, PlanStandard1, PlanStandard2}
)

// PlanID generates Plan ID from region name and plan name
func PlanID(region, plan string) string {
	return fmt.Sprintf("%s/%s", region, plan)
}
