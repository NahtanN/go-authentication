package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nahtann/go-lab/cmd/rest_api"
)

var (
	port     int
	rootPath string
	swagger  bool
)

var restApi = &cobra.Command{
	Use:   "rest-api",
	Short: "Run rest api",
	Run: func(cmd *cobra.Command, args []string) {
		formattedPort := fmt.Sprintf(":%d", port)
		formattedRootPath := fmt.Sprintf("/%s", rootPath)

		rest_api.RunServer(formattedPort, formattedRootPath, swagger)
	},
}

func init() {
	restApi.Flags().IntVarP(&port, "port", "p", 3333, "define server port.")
	restApi.Flags().StringVarP(&rootPath, "root-path", "r", "api", "define root path for server.")
	restApi.Flags().BoolVarP(&swagger, "swagger", "s", false, "Toggle swagger docs. Default false")

	rootCmd.AddCommand(restApi)
}
