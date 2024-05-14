package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var build, format bool

var swag = &cobra.Command{
	Use:   "swag",
	Short: "Manage swagger files",
	Long:  "You should have swag CLI installed. https://github.com/swaggo/swag",
	Run: func(cmd *cobra.Command, args []string) {
		if format {
			err := exec.Command("/bin/sh", "-c", "swag fmt").Run()
			if err != nil {
				fmt.Printf("Error executing command: %s\n", err)
				return
			}

			fmt.Println("Swagger files formatted successfully")
		}

		if build {
			out, err := exec.Command("/bin/sh", "-c", "swag init").Output()
			if err != nil {
				fmt.Printf("Error executing command: %s\n", err)
				return
			}

			fmt.Printf("Swagger files created successfully: %s\n", string(out))
		}
	},
}

func init() {
	swag.Flags().BoolVarP(&build, "init", "i", false, "Creates swagger doc files.")
	swag.Flags().BoolVarP(&format, "format", "f", false, "Format swag doc comments.")

	rootCmd.AddCommand(swag)
}
