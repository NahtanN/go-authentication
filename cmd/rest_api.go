package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nahtann/go-lab/cmd/rest_api"
)

var (
	port     int
	rootPath string
)

var restApi = &cobra.Command{
	Use:   "rest-api",
	Short: "Run rest api",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		rootPath, _ := cmd.Flags().GetString("root-path")

		formattedPort := fmt.Sprintf(":%d", port)
		formattedRootPath := fmt.Sprintf("/%s", rootPath)

		rest_api.RunServer(formattedPort, formattedRootPath)
	},
}

func init() {
	restApi.Flags().IntVarP(&port, "port", "p", 3333, "define server port.")
	restApi.Flags().StringVarP(&rootPath, "root-path", "r", "api", "define root path for server.")

	rootCmd.AddCommand(restApi)
}
