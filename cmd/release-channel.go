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

    var selectedReleaseChannel string
    currentReleaseChannel := getCurrentReleaseChannel()
    fmt.Printf("Current release channel: %s\n", currentReleaseChannel)
    form := huh.NewForm(
        huh.NewGroup(
            huh.NewSelect[string]().
                Title("Select release channel to update from").
                Options(
                    currentMarker("2.x (stable)", "2.x", currentReleaseChannel),
                    currentMarker("3.x (preview)", "3.x", currentReleaseChannel),
                    currentMarker("next (development)", "next", currentReleaseChannel),
                ).
                Value(&selectedReleaseChannel),
        ),
    )

    err := form.Run()
    if err != nil {
        return err
    }

    return processBranchSelection(selectedReleaseChannel)
}

func processBranchSelection(branch string) error {
    switch branch {
    case "2.x":
        return useStableChannel()
    case "next":
        return useNextChannel()
    case "3.x":
        return usePreviewChannel()
    default:
        return fmt.Errorf("invalid branch: %s, must be 'stable' or 'next'", branch)
    }
}

func currentMarker(label, value, currentReleaseChannel string) huh.Option[string] {
    if value == currentReleaseChannel {
        return huh.NewOption(label+" - current", value)
    }
    return huh.NewOption(label, value)
}

func getCurrentReleaseChannel() string {
    releaseChannelFilePath := filepath.Join(constants.ValetEtcPath, constants.ReleaseChannelFile)
    releaseChannel, err := os.ReadFile(releaseChannelFilePath)
    if err != nil {
        return "stable"
    }
    return string(releaseChannel)
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

    releaseChannelFilePath := filepath.Join(constants.ValetEtcPath, constants.ReleaseChannelFile)
    err := os.WriteFile(releaseChannelFilePath, []byte("2.x"), 0644)
    if err != nil {
        return fmt.Errorf("failed to switch to stable channel: %w", err)
    }
    fmt.Println("Successfully switched to stable channel")

    return runUpdate()
}

func usePreviewChannel() error {
    fmt.Println("Switching to preview channel")

    if err := ensureEtcDirectoryExists(); err != nil {
        return err
    }

    releaseChannelFilePath := filepath.Join(constants.ValetEtcPath, constants.ReleaseChannelFile)
    err := os.WriteFile(releaseChannelFilePath, []byte("3.x"), 0644)
    if err != nil {
        return fmt.Errorf("failed to switch to preview channel: %w", err)
    }
    fmt.Println("Successfully switched to preview channel")

    return runUpdate()
}

func useNextChannel() error {
    fmt.Println("Switching to next channel")

    if err := ensureEtcDirectoryExists(); err != nil {
        return err
    }

    releaseChannelFilePath := filepath.Join(constants.ValetEtcPath, constants.ReleaseChannelFile)
    err := os.WriteFile(releaseChannelFilePath, []byte("next"), 0644)
    if err != nil {
        return fmt.Errorf("failed to switch to next channel: %w", err)
    }
    fmt.Println("Successfully switched to next channel")

    return runUpdate()
}
