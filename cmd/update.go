package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
    //"github.com/charmbracelet/huh"

    //"github.com/valet-sh/valet-sh-updater/internal/git"
    //"github.com/valet-sh/valet-sh-updater/internal/runtime"
)

var updateCmd = &cobra.Command{
    Use:  "update",
    Short: "Update valet-sh to the latest version",
    Long: `Update valet-sh to the latest version`,
    RunE: func(cmd *cobra.Command, args []string) error {
        return runUpdate()
    },
}

var (
    branchFlag string
)

func init() {
	updateCmd.Flags().StringVarP(&branchFlag, "branch", "b", "", "Branch to use (stable or next)")
}

func runUpdate() error {
    fmt.Println("Checking for updates...")
    return nil
}
