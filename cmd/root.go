package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/agrim123/gatekeeper/cmd/gatekeeper"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gatekeeper",
	Short: "Gatekeeper is your authentication and authorization oriented deployment managment tool.",
	Long:  ``,
}

var runPlanCmd = &cobra.Command{
	Use:   "run-plan",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		g := gatekeeper.NewGatekeeper(ctx)

		switch len(args) {
		case 1:
			g.Runtime.SetPlan(args[0])
		default:
			g.Runtime.SetPlan(args[0])
			g.Runtime.SetOption(args[1])
		}

		g.Run()
	},
}

func init() {
	rootCmd.AddCommand(runPlanCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
