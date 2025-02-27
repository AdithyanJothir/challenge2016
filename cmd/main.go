package main

import (
	"CHALLENGE2016/internal/models/distributor"
	"CHALLENGE2016/internal/models/region"
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// initializeRegions reads the CSV file and adds regions to the RegionManager
func initializeRegions(rm *region.RegionManager, csvFile string) error {
	file, err := os.Open(csvFile)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV file: %w", err)
	}

	for i, line := range lines {
		if i == 0 {
			continue
		}
		if len(line) < 6 {
			continue
		}
		cityName := line[3]
		provinceName := line[4]
		countryName := line[5]

		regionPath := fmt.Sprintf("%s-%s-%s", countryName, provinceName, cityName)
		rm.AddRegion(regionPath)
	}
	return nil
}

func main() {
	// Initialize RegionManager and load regions from CSV
	rm := region.NewRegionManager()
	err := initializeRegions(rm, "../internal/assets/cities.csv")
	if err != nil {
		fmt.Println("Error initializing regions:", err)
		return
	}

	distributors := make(map[string]*distributor.Distributor)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n----- Menu -----")
		fmt.Println("1. Create Distributor")
		fmt.Println("2. Add Region to Distributor")
		fmt.Println("3. Exclude Region from Distributor")
		fmt.Println("4. Add Parent to Distributor")
		fmt.Println("5. Check Distributor Permission")
		fmt.Println("6. List Distributors")
		fmt.Println("7. Exit")
		fmt.Print("Enter option: ")

		if !scanner.Scan() {
			break
		}
		option := strings.TrimSpace(scanner.Text())

		switch option {
		case "1":
			// Create a new distributor
			fmt.Print("Enter distributor ID (unique string): ")
			scanner.Scan()
			id := strings.TrimSpace(scanner.Text())

			fmt.Print("Enter distributor Name: ")
			scanner.Scan()
			name := strings.TrimSpace(scanner.Text())

			// Input authorized regions as comma-separated region paths
			fmt.Print("Enter authorized region paths (comma separated, e.g. 'India-Tamil Nadu-Keelakarai, India-Jammu and Kashmir-Punch'): ")
			scanner.Scan()
			regionsInput := scanner.Text()
			regionPaths := strings.Split(regionsInput, ",")

			authRegions := make(map[string]*region.Region)
			for _, rp := range regionPaths {
				rp = strings.TrimSpace(rp)
				if rp == "" {
					continue
				}
				if r, exists := rm.GetRegion(rp); exists {
					authRegions[rp] = r
				} else {
					fmt.Printf("Region '%s' not found. Skipping.\n", rp)
				}
			}

			dist := distributor.AddDistributor(id, name, authRegions)
			distributors[id] = dist
			fmt.Println("Distributor created.")

		case "2":
			// Add a region to a distributor.
			fmt.Print("Enter distributor ID: ")
			scanner.Scan()
			id := strings.TrimSpace(scanner.Text())
			dist, ok := distributors[id]
			if !ok {
				fmt.Println("Distributor not found.")
				break
			}

			fmt.Print("Enter region path to add (e.g. 'India-Tamil Nadu-Keelakarai'): ")
			scanner.Scan()
			rp := strings.TrimSpace(scanner.Text())
			if r, exists := rm.GetRegion(rp); exists {
				if err := dist.AddRegion(r); err != nil {
					fmt.Println("Error:", err)
				} else {
					fmt.Println("Region added to distributor.")
				}
			} else {
				fmt.Println("Region not found.")
			}

		case "3":
			// Exclude a region from a distributor
			fmt.Print("Enter distributor ID: ")
			scanner.Scan()
			id := strings.TrimSpace(scanner.Text())
			dist, ok := distributors[id]
			if !ok {
				fmt.Println("Distributor not found.")
				break
			}

			fmt.Print("Enter region path to exclude: ")
			scanner.Scan()
			rp := strings.TrimSpace(scanner.Text())
			if r, exists := rm.GetRegion(rp); exists {
				dist.ExcludeRegion(r)
				fmt.Println("Region excluded from distributor.")
			} else {
				fmt.Println("Region not found.")
			}

		case "4":
			// Add a parent distributor
			fmt.Print("Enter child distributor ID: ")
			scanner.Scan()
			childID := strings.TrimSpace(scanner.Text())
			childDist, ok := distributors[childID]
			if !ok {
				fmt.Println("Child distributor not found.")
				break
			}

			fmt.Print("Enter parent distributor ID: ")
			scanner.Scan()
			parentID := strings.TrimSpace(scanner.Text())
			parentDist, ok := distributors[parentID]
			if !ok {
				fmt.Println("Parent distributor not found.")
				break
			}

			childDist.AddParent(parentDist)
			fmt.Println("Parent added to distributor.")

		case "5":

			// Check if a distributor has permission for a given region.
			fmt.Print("Enter distributor ID: ")
			scanner.Scan()
			id := strings.TrimSpace(scanner.Text())
			dist, ok := distributors[id]
			if !ok {
				fmt.Println("Distributor not found.")
				break
			}

			fmt.Print("Enter region path to check (e.g. 'India-Tamil Nadu-Keelakarai'): ")
			scanner.Scan()
			rp := strings.TrimSpace(scanner.Text())
			if r, exists := rm.GetRegion(rp); exists {
				if dist.HasPermission(r) {
					fmt.Println("Permission granted!")
				} else {
					fmt.Println("Permission denied!")
				}
			} else {
				fmt.Println("Region not found.")
			}

		case "6":
			// List all distributors.
			fmt.Println("Listing Distributors:")
			for id, dist := range distributors {
				fmt.Println("--------------------------------")
				fmt.Printf("ID: %s, Name: %s\n", id, dist.Name)
				fmt.Print("\n  Excluded Regions: ")
				for regID := range dist.ExcludedRegions {
					fmt.Printf("%s ", regID)
				}
				fmt.Print("\n  Effective Authorized Regions: ")
				for regID := range dist.EffectiveAuthorizedRegions {
					fmt.Printf("%s ", regID)
				}
				fmt.Println("\n--------------------------------")
			}

		case "7":
			fmt.Println("Exiting.")
			return

		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
