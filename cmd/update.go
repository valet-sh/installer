package cmd

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"

    "github.com/valet-sh/valet-sh-installer/constants"
    "github.com/valet-sh/valet-sh-installer/internal/git"
    "github.com/valet-sh/valet-sh-installer/internal/runtime"
)

var updateCmd = &cobra.Command{
    Use:  "update",
    Short: "Update valet-sh to the latest version",
    Long: `Update valet-sh to the latest version`,
    SilenceUsage: true,
    RunE: func(cmd *cobra.Command, args []string) error {
        return runUpdate()
    },
}

func init() {
}

func runUpdate() error {
    repoPath := constants.ValetBasePath

    if err := checkIfRepoExists(repoPath); err != nil {
        return err
    }

    nextChannelEnabled := false
    enableNextFilePath := filepath.Join(constants.ValetEtcPath, constants.NextBranchFile)

    if _, err := os.Stat(enableNextFilePath); err == nil {
        nextChannelEnabled = true
    }

    if nextChannelEnabled {
        fmt.Println("Using next channel (development) for update")
        return updateNextBranch(repoPath)
    } else {
        fmt.Println("Using stable channel for update")
        return updateStableBranch(repoPath)
    }
}

func checkIfRepoExists(repoPath string) error {
    _, err := os.Stat(repoPath)
    if os.IsNotExist(err) {
        return fmt.Errorf("valet-sh not found in %s", repoPath)
    }
    return nil
}

func updateNextBranch(repoPath string) error {
    fmt.Println("Updating valet-sh to the latest version on the next branch")

    if err := git.CheckoutBranch(repoPath, "next"); err != nil {
        return fmt.Errorf("Failed to checkout next branch: %w", err)
    }

    if err := git.PullLatest(repoPath); err != nil {
        return fmt.Errorf("Failed to pull latest changes: %w", err)
    }
    return runtime.CheckRuntime()
}

func updateStableBranch(repoPath string) error {
    fmt.Println("Updating valet-sh to the latest version on the stable branch")
    return nil
}
