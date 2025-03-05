package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var selfUpgradeCmd = &cobra.Command{
    Use:  "self-upgrade",
    Short: "Update valet-sh-installer to the latest version",
    Long: `Update valet-sh-installer to the latest version`,
    SilenceUsage: true,
    RunE: func(cmd *cobra.Command, args []string) error {
        return selfUpgrade()
    },
}

func init() {
}

func selfUpgrade() error {
    fmt.Println("self-upgrade")
    return nil
}
