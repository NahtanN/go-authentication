package cmd

import (
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
		if create != "" {
			createMigration(create)
			return
		}

		if up {
			return
		}

		if down {
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
	// migrationsDir := "internal/storage/database/migrations"

	/*err := exec.Command("/bin/sh", "-c", "cd internal/storage/database/migrations && migrate create -ext sql teeessssssstttt").*/
	/*Run()*/
	/*if err != nil {*/
	/*fmt.Printf("Error executing command: %s\n", err)*/
	/*return*/
	/*}*/

	//fmt.Println(string(out))
	/*if fileName != "" {*/
	/*create := "touch"*/

	/*out, err := exec.Command(create, fileName).Output()*/
	/*if err != nil {*/
	/*fmt.Printf("Error executing command: %s\n", err)*/
	/*return*/
	/*}*/

	/*} else {*/
	/*fmt.Println("file does not set")*/
	/*}*/
}

func upMigrations() {}

func downMigrations() {}
