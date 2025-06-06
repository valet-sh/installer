package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/valet-sh/valet-sh-installer/internal/utils"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/valet-sh/valet-sh-installer/constants"
	"github.com/valet-sh/valet-sh-installer/internal/prechecks"
)

var setReleaseChannelCmd = &cobra.Command{
	Use:           "release-channel",
	Short:         "Set the release channel to update from",
	Long:          `Set the release channel to update from`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			err := processBranchSelection(args[0])
			if err != nil {
				color.Error.Printf("Error: %s\n", err.Error())
			}
			return err
		}
		err := setReleaseChannel()
		if err != nil {
			color.Error.Printf("Error: %s\n", err.Error())
			return err
		}
		return nil
	},
}

func init() {
}

func setReleaseChannel() error {
	repoPath := constants.VshBasePath

	if err := checkIfRepoExists(repoPath); err != nil {
		return err
	}

	var selectedReleaseChannel string
	currentReleaseChannel := getCurrentReleaseChannel()
	fmt.Printf("Current release channel: %s\n\n", currentReleaseChannel)
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select release channel to update from").
				Options(
					currentMarker("2.x (stable)", "2.x", currentReleaseChannel),
					// currentMarker("3.x (preview)", "3.x", currentReleaseChannel),
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
		//case "3.x":
		//		return usePreviewChannel()
	default:
		return fmt.Errorf("invalid branch: %s, must be 'stable' or 'next'", branch)
	}
}

func currentMarker(label, value, currentReleaseChannel string) huh.Option[string] {
	if value == currentReleaseChannel {
		return huh.NewOption(label+" - current", value).Selected(true)
	}
	return huh.NewOption(label, value)
}

func getCurrentReleaseChannel() string {
	releaseChannelFilePath := filepath.Join(constants.VshEtcPath, constants.ReleaseChannelFile)
	releaseChannel, err := os.ReadFile(releaseChannelFilePath)
	if err != nil {
		return "2.x"
	}
	return strings.TrimSpace(string(releaseChannel))
}

func useStableChannel() error {
	utils.Println("Switching to stable channel")

	if err := prechecks.CheckForEtcDirectory(); err != nil {
		return err
	}

	releaseChannelFilePath := filepath.Join(constants.VshEtcPath, constants.ReleaseChannelFile)
	err := os.WriteFile(releaseChannelFilePath, []byte("2.x"), 0644)
	if err != nil {
		return fmt.Errorf("failed to switch to stable channel: %w", err)
	}
	color.Info.Println("\nSuccessfully switched to stable channel\n")

	return runUpdate()
}

func usePreviewChannel() error {
	utils.Println("Switching to preview channel")

	if err := prechecks.CheckForEtcDirectory(); err != nil {
		return err
	}

	releaseChannelFilePath := filepath.Join(constants.VshEtcPath, constants.ReleaseChannelFile)
	err := os.WriteFile(releaseChannelFilePath, []byte("3.x"), 0644)
	if err != nil {
		return fmt.Errorf("failed to switch to preview channel: %w", err)
	}
	color.Info.Println("\nSuccessfully switched to preview channel\n")

	return runUpdate()
}

func useNextChannel() error {
	utils.Println("Switching to next channel")

	if err := prechecks.CheckForEtcDirectory(); err != nil {
		return err
	}

	releaseChannelFilePath := filepath.Join(constants.VshEtcPath, constants.ReleaseChannelFile)
	err := os.WriteFile(releaseChannelFilePath, []byte("next"), 0644)
	if err != nil {
		return fmt.Errorf("failed to switch to next channel: %w", err)
	}

	color.Info.Println("\nSuccessfully switched to next channel\n")

	return runUpdate()
}
