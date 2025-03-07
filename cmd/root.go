package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

func preflightChecks(cmd *cobra.Command, args []string) error {
    fmt.Println("Running preflight checks now after root command")
    return nil
}

var rootCmd = &cobra.Command{
    Use:   "valet-sh-installer",
    Short: "A CLI tool to update Valet-sh",
    Long: `A CLI tool to update Valet-sh`,
    PersistentPreRunE: preflightChecks,
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    rootCmd.AddCommand(setChannelCmd)
    rootCmd.AddCommand(installCmd)
    rootCmd.AddCommand(updateCmd)
    rootCmd.AddCommand(selfUpgradeCmd)
}
