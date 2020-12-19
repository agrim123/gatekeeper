package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/agrim123/gatekeeper/internal/constants"
	"github.com/agrim123/gatekeeper/internal/gatekeeper"
	"github.com/agrim123/gatekeeper/internal/pkg/store"
	"github.com/agrim123/gatekeeper/pkg/config"
	"github.com/agrim123/gatekeeper/pkg/logger"
	"github.com/agrim123/gatekeeper/pkg/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gatekeeper",
	Short: "Gatekeeper is an authentication and authorization oriented tool allowing non-root users to ssh to a machine without giving them access to private keys.",
	Long:  ``,
}

var runPlanCmd = &cobra.Command{
	Use:   "run-plan",
	Short: "Runs the specified plan with given option",
	Long: `Runs plans defined in plan.json. Also takes an option as second argument.
			For example: gatekeeper run-plan user_service deploy`,
	Run: func(cmd *cobra.Command, args []string) {
		plan := ""
		option := ""
		switch len(args) {
		case 0:
		case 1:
			plan = args[0]
		default:
			plan = args[0]
			option = args[1]
		}

		ctx := utils.AttachOptionToCtx(
			utils.AttachPlanToCtx(
				utils.AttachExecutingUserToCtx(
					context.Background(),
				),
				plan,
			),
			option,
		)

		gatekeeper.NewGatekeeper(ctx).Run(plan, option)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all commands the user can run",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := utils.AttachExecutingUserToCtx(context.Background())
		store.Init()
		allowedCmds := store.Store.GetAllowedCommands(ctx.Value(constants.UserContextKey).(string))
		if len(allowedCmds) == 0 {
			logger.Info("No allowed commands")
			return
		}

		for plan, options := range allowedCmds {
			fmt.Println(plan)
			for _, opt := range options {
				fmt.Println("    " + opt)
			}
		}
	},
}

var selfCmd = &cobra.Command{
	Use:   "self",
	Short: "Gatekeeper management commands",
}

var selfUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates gatekeeper code",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Error("TODO")
	},
}

func init() {
	rootCmd.AddCommand(runPlanCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(selfCmd)
	selfCmd.AddCommand(selfUpdateCmd)
}

// Execute is the entrypoint of gatekeeper
func Execute() {
	// load configs to memory
	config.Init()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
