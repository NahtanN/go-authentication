package cmd

import (
	"github.com/spf13/cobra"

	"github.com/nahtann/go-lab/cmd/rest_api"
)

var run = &cobra.Command{
	Use:   "run",
	Short: "Run some specified resource",
	Run: func(cmd *cobra.Command, args []string) {
		restApi, _ := cmd.Flags().GetBool("rest-api")

		if restApi {
			rest_api.RunServer()
		}
	},
}

func init() {
	run.Flags().Bool("rest-api", true, "start a rest api server.")

	rootCmd.AddCommand(run)
}
