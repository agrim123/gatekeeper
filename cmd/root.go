package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/agrim123/gatekeeper/internal/constants"
	"github.com/agrim123/gatekeeper/internal/gatekeeper"
	"github.com/agrim123/gatekeeper/internal/pkg/setup"
	"github.com/agrim123/gatekeeper/internal/pkg/store"
	"github.com/agrim123/gatekeeper/pkg/config"
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
			g.Run("", "")
		case 1:
			g.Run(args[0], "")
		default:
			g.Run(args[0], args[1])
		}
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all commands the user can run",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := utils.AttachExecutingUserToCtx(context.Background())
		allowedCmds := store.GetAllowedCommands(ctx.Value(constants.UserContextKey).(string))
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
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var selfUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates gatekeeper code",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(runPlanCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(selfCmd)
	selfCmd.AddCommand(selfUpdateCmd)
}

func Execute() {
	config.Init()

	setup.Init()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
