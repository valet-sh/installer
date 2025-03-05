package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
    "github.com/charmbracelet/huh"

    "github.com/valet-sh/valet-sh-installer/constants"
    // "github.com/valet-sh/valet-sh-installer/internal/git"
    // "github.com/valet-sh/valet-sh-installer/internal/runtime"
)

var setChannelCmd = &cobra.Command{
    Use:  "set-channel",
    Short: "Set the channel to update from",
    Long: `Set the channel to update from`,
    SilenceUsage: true,
    RunE: func(cmd *cobra.Command, args []string) error {
        return setChannel()
    },
}

var (
    branchFlag string
)

func init() {
    setChannelCmd.Flags().StringVarP(&branchFlag, "branch", "b", "", "Branch to use (stable or next)")
}

func setChannel() error {
    repoPath := constants.ValetBasePath

    if err := checkIfRepoExists(repoPath); err != nil {
        return err
    }

    if branchFlag != "" {
        return processBranchSelection(branchFlag)
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

    err := form.Run()
    if err != nil {
        return err
    }

    return processBranchSelection(selectedBranch)
}

func processBranchSelection(branch string) error {
    switch branch {
    case "stable":
        return useStableChannel()
    case "next":
        return useNextChannel()
    default:
        return fmt.Errorf("invalid branch: %s, must be 'stable' or 'next'", branch)
    }
}

func useStableChannel() error {
    fmt.Println("Using stable channel")
    return nil
}

func useNextChannel() error {
    fmt.Println("Using next channel")
    return nil
}
