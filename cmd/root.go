package cmd

import (
    "fmt"

    "github.com/spf13/cobra"

    "github.com/valet-sh/valet-sh-installer/internal/prechecks"
)

func preflightChecks(cmd *cobra.Command, args []string) error {
    fmt.Println("Running preflight checks now after root command")

    if err := prechecks.CheckForValet(); err != nil {
        return err
    }

    if err := prechecks.CheckForEtcDirectory(); err != nil {
        return prechecks.CheckForEtcDirectory()
    }

    if err := prechecks.CheckForValetMajorReleaseFile(); err != nil {
        return prechecks.CheckForValetMajorReleaseFile()
    }

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
