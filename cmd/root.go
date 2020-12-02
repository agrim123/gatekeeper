package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/agrim123/gatekeeper/internal/gatekeeper"
	"github.com/agrim123/gatekeeper/pkg/logger"
	"github.com/agrim123/gatekeeper/pkg/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gatekeeper",
	Short: "Gatekeeper is your authentication and authorization oriented deployment managment tool.",
	Long:  ``,
}

var runPlanCmd = &cobra.Command{
	Use:   "run-plan",
	Short: "Runs the specified plan with given option",
	Long: `Runs plans defined in plan.json. Also takes an option as second argument.
			For example: gatekeeper run-plan user_service deploy`,
	Run: func(cmd *cobra.Command, args []string) {
		g := gatekeeper.NewGatekeeper(utils.AttachExecutingUserToCtx(context.Background()))

		switch len(args) {
		case 0:
			logger.Fatal("Invalid arguments")
		case 1:
			g.Run(args[0], "")
		default:
			g.Run(args[0], args[1])
		}
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
