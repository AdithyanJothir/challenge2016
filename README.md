
# Go Based Distribution Program

  

## Overview

  

This command-line tool, written in Go, allows you to manage distributors and their permissions based on geographic regions. With this program, you can:

  

-  **Create Distributors:** Define distributors with specific authorized regions.

-  **Add Regions:** Add new regions to a distributor’s authorized list.

-  **Exclude Regions:** Exclude certain regions from a distributor’s effective permissions.

-  **Establish Parent-Child Relationships:** Set up hierarchical relationships so that a child's permissions are always a subset of its parent's.

-  **Check Permissions:** Verify if a distributor is permitted to operate in a given region.

  

## Getting Started

  

### Prerequisites

  

-  **Go:** Ensure you have Go installed (version 1.16 or later).

-  **CSV File:** The program requires a `cities.csv` file (located in the `internal/assets` folder) with region data. The CSV should include columns like:

- City Code

- Province Code

- Country Code

- City Name

- Province Name

- Country Name

### How to Run?

1.  **Build the Program:**

```bash
cd main
go build ./main
```

1.  **Run the Program:**


```bash
 ./main
```
