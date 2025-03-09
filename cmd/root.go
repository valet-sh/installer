package cmd

import (
    "github.com/spf13/cobra"

    "github.com/valet-sh/valet-sh-installer/internal/prechecks"
)

func preflightChecks(cmd *cobra.Command, args []string) error {
    if cmd.Name() == "setup" || cmd.Name() == "self-upgrade" {
        return nil
    }

    if err := prechecks.CheckForValet(); err != nil {
        return err
    }

    if err := prechecks.CheckForEtcDirectory(); err != nil {
        return prechecks.CheckForEtcDirectory()
    }

    if err := prechecks.CheckForValetReleaseChannelFile(); err != nil {
        return prechecks.CheckForValetReleaseChannelFile()
    }

    return nil
}

var rootCmd = &cobra.Command{
    Use:   "valet-sh-installer",
    Short: "A CLI tool to update Valet-sh",
    Long: `A CLI tool to update Valet-sh`,
    Version: "0.0.1",
    SilenceErrors: true,
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        return preflightChecks(cmd, args)
    },
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    rootCmd.AddCommand(setReleaseChannelCmd)
    rootCmd.AddCommand(setupCmd)
    rootCmd.AddCommand(updateCmd)
    rootCmd.AddCommand(selfUpgradeCmd)
}
