package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/luongtruong20201/bookmark-management/internal/infrastructure"
	"github.com/luongtruong20201/bookmark-management/pkg/common"
	sqldb "github.com/luongtruong20201/bookmark-management/pkg/sql"
)

func main() {
	var (
		mode  = flag.String("mode", "up", "Migration mode: up, down, or steps")
		steps = flag.Int("steps", 0, "Number of steps to migrate (positive for up, negative for down). Required when mode is 'steps'")
	)
	flag.Parse()

	db, err := sqldb.NewClient("")
	common.HandleError(err)

	migrationPath := "file://./migrations"
	var migrationErr error

	switch *mode {
	case "up":
		migrationErr = infrastructure.MigrateSQLDB(db, migrationPath, "up", 0)
	case "down":
		migrationErr = infrastructure.MigrateSQLDB(db, migrationPath, "down", 0)
	case "steps":
		if *steps == 0 {
			fmt.Fprintf(os.Stderr, "Error: --steps must be specified and non-zero when using mode 'steps'\n")
			fmt.Fprintf(os.Stderr, "Usage: go run cmd/migrate/main.go -mode=steps -steps=1 (or -steps=-1 for down)\n")
			os.Exit(1)
		}
		migrationErr = infrastructure.MigrateSQLDB(db, migrationPath, "steps", *steps)
	default:
		fmt.Fprintf(os.Stderr, "Error: Invalid mode '%s'. Use 'up', 'down', or 'steps'\n", *mode)
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  go run cmd/migrate/main.go -mode=up\n")
		fmt.Fprintf(os.Stderr, "  go run cmd/migrate/main.go -mode=down\n")
		fmt.Fprintf(os.Stderr, "  go run cmd/migrate/main.go -mode=steps -steps=1\n")
		os.Exit(1)
	}

	if migrationErr != nil {
		fmt.Fprintf(os.Stderr, "Migration error: %v\n", migrationErr)
		os.Exit(1)
	}

	fmt.Printf("Migration completed successfully: mode=%s", *mode)
	if *mode == "steps" {
		fmt.Printf(", steps=%d", *steps)
	}
	fmt.Println()
}
