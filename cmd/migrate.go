package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	create string
	up     bool
	down   bool
)

var migrate = &cobra.Command{
	Use:   "migrate",
	Short: "Manage migration files",
	Run: func(cmd *cobra.Command, args []string) {
		// load .env variables
		err := godotenv.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading .env file")
			os.Exit(1)
		}

		if create != "" {
			createMigration(create)
			return
		}

		if up {
			runMigrations("up")
			return
		}

		if down {
			runMigrations("down")
			return
		}
	},
}

func init() {
	migrate.Flags().
		StringVarP(&create, "create", "c", "", "Create a new migration file.")
	migrate.Flags().
		BoolVarP(&up, "up", "u", false, "Apply all up migrations.")
	migrate.Flags().
		BoolVarP(&down, "down", "d", false, "Apply all down migrations.")

	rootCmd.AddCommand(migrate)
}

func createMigration(fileName string) {
	migrationsPath := os.Getenv("DATABASE_MIGRATIONS_PATH")

	command := fmt.Sprintf(
		"migrate create -dir %s -ext sql %s",
		migrationsPath,
		fileName,
	)

	err := exec.Command("/bin/sh", "-c", command).
		Run()
	if err != nil {
		fmt.Printf("Error executing command: %s\n", err)
		return
	}

	fmt.Printf("Migration file `%s` created successfully.\n", fileName)
}

func runMigrations(direction string) {
	if direction != "up" && direction != "down" {
		fmt.Println("Migrations should be 'up' or 'down'")
		return
	}

	databaseUrl := os.Getenv("DATABASE_URL")
	migrationsPath := os.Getenv("DATABASE_MIGRATIONS_PATH")

	command := fmt.Sprintf(
		"migrate -path ./%s -database '%s' %s",
		migrationsPath,
		databaseUrl,
		direction,
	)

	cmd := exec.Command("/bin/sh", "-c", command)

	if direction == "down" {
		input := bytes.NewBufferString("y")
		cmd.Stdin = input
	}

	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing command: %s\n", err)
		return
	}

	fmt.Printf("Migration %s successfully.\n%s", direction, string(out))
}
