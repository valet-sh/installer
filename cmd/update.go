package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/charmbracelet/huh"

    "github.com/valet-sh/valet-sh-updater/constants"
    "github.com/valet-sh/valet-sh-updater/internal/git"
    "github.com/valet-sh/valet-sh-updater/internal/runtime"
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

var (
    branchFlag string
)

func init() {
    updateCmd.Flags().StringVarP(&branchFlag, "branch", "b", "", "Branch to use (stable or next)")
}

func runUpdate() error {
    branch, err := determineBranch()
    if err != nil {
     return err
    }

    repoPath := constants.ValetBasePath
    if err := checkIfRepoExists(repoPath); err != nil {
        return err
    }

    if branch == "next" {
        return updateNextBranch(repoPath)
    } else {
        return nil
    }

    fmt.Println("Just a test - branch:", branch)
    return nil
}

func determineBranch() (string, error) {
    if branchFlag == "next" || branchFlag == "stable" {
        return branchFlag, nil
    }

    _, err := os.Stat(constants.NextBranchFile)
    if err == nil {
        return "next", nil
    }

    var selectedBranch string
    form := huh.NewForm(
        huh.NewGroup(
            huh.NewSelect[string]().
                Title("Select branch to update from").
                Options(
                    huh.NewOption("Stable (production use)", "stable"),
                    huh.NewOption("Next (development use)", "next"),
                ).
                Value(&selectedBranch),
        ),
    )

    err = form.Run()
    if err != nil {
        return "stable", err
    }

    return selectedBranch, nil
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
