package cmd

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "valet-sh-updater",
    Short: "A CLI tool to update Valet-sh",
    Long: `A CLI tool to update Valet-sh`,
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    rootCmd.AddCommand(updateCmd)
}
