package cmd

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"
    "github.com/charmbracelet/huh"

    "github.com/valet-sh/valet-sh-installer/constants"
    // "github.com/valet-sh/valet-sh-installer/internal/git"
    // "github.com/valet-sh/valet-sh-installer/internal/runtime"
)

var setReleaseChannelCmd= &cobra.Command{
    Use:  "release-channel",
    Short: "Set the release channel to update from",
    Long: `Set the release channel to update from`,
    SilenceUsage: true,
    RunE: func(cmd *cobra.Command, args []string) error {
        if len(args) == 1 {
            return processBranchSelection(args[0])
        }
        return setReleaseChannel()
    },
}

func init() {
}

func setReleaseChannel() error {
    repoPath := constants.ValetBasePath

    if err := checkIfRepoExists(repoPath); err != nil {
        return err
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

func ensureEtcDirectoryExists() error {
    if _, err := os.Stat(constants.ValetEtcPath); os.IsNotExist(err) {
        err := os.MkdirAll(constants.ValetEtcPath, 0755)
        if err != nil {
            return fmt.Errorf("failed to create etc directory: %w", err)
        }
    }
    return nil
}

func useStableChannel() error {
    fmt.Println("Switching to stable channel")

    if err := ensureEtcDirectoryExists(); err != nil {
        return err
    }

    enableNextFilePath := filepath.Join(constants.ValetEtcPath, constants.NextBranchFile)
    if _, err := os.Stat(enableNextFilePath); err == nil {
        fmt.Println("Removing next channel file")
        err := os.Remove(enableNextFilePath)
        if err != nil {
            return fmt.Errorf("failed to disable next channel: %w", err)
        }
        fmt.Println("Successfully switched to stable channel")
    } else {
        fmt.Println("Already on stable channel")
    }

    return runUpdate()
}

func useNextChannel() error {
    fmt.Println("Switching to next channel")

    if err := ensureEtcDirectoryExists(); err != nil {
        return err
    }

    enableNextFilePath := filepath.Join(constants.ValetEtcPath, constants.NextBranchFile)
    if _, err := os.Stat(enableNextFilePath); os.IsNotExist(err) {
        fmt.Println("Creating next channel file")
        _, err := os.Create(enableNextFilePath)
        if err != nil {
            return fmt.Errorf("failed to enable next channel: %w", err)
        }
        fmt.Println("Successfully switched to next channel")
    } else if err != nil {
        return fmt.Errorf("error checking next channel file: %w", err)
    } else {
        fmt.Println("Already on next channel")
    }

    return runUpdate()
}
