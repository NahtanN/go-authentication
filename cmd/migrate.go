package cmd

import "github.com/spf13/cobra"

var migrate = &cobra.Command{
	Use:   "migrate",
	Short: "Manage migration files",
}

func init() {
	rootCmd.AddCommand(migrate)
}
