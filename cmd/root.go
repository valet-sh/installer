package cmd

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "valet-sh-installer",
    Short: "A CLI tool to update Valet-sh",
    Long: `A CLI tool to update Valet-sh`,
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    rootCmd.AddCommand(setChannelCmd)
    rootCmd.AddCommand(updateCmd)
    rootCmd.AddCommand(selfUpgradeCmd)
}
