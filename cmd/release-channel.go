package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/valet-sh/valet-sh-installer/internal/migration"
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
	selectedReleaseChannel = currentReleaseChannel

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select release channel to update from").
				Options(
					currentMarker("2.x (stable)", "2.x", currentReleaseChannel),
					currentMarker("3.x", "3.x", currentReleaseChannel),
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
		return use3xChannel()
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
	utils.Println("\nSuccessfully switched to stable channel\n")

	return runUpdate()
}

func use3xChannel() error {

	utils.Println("Checking requirements for 3.x channel")

	if err := prechecks.Requirements3xCheck(); err != nil {
		fmt.Println("Requirements check failed:", err)
		os.Exit(1)
	}

	migrationConfirmed, err := prechecks.CheckMigration()
	if err != nil {
		return err
	}

	utils.Println("Switching to 3.x channel")

	if err := prechecks.CheckForEtcDirectory(); err != nil {
		return err
	}

	releaseChannelFilePath := filepath.Join(constants.VshEtcPath, constants.ReleaseChannelFile)
	previousReleaseChannel := getCurrentReleaseChannel()

	if err := os.WriteFile(releaseChannelFilePath, []byte("3.x"), 0644); err != nil {
		return fmt.Errorf("failed to switch to 3.x channel: %w", err)
	}
	utils.Printf("\nSuccessfully switched to 3.x channel\n")

	if err := runUpdate(); err != nil {
		revertReleaseChannel(releaseChannelFilePath, previousReleaseChannel)
		return fmt.Errorf("failed to update after switching to 3.x channel: %w", err)
	}

	if err := migration.Cli(); err != nil {
		revertReleaseChannel(releaseChannelFilePath, previousReleaseChannel)
		return fmt.Errorf("failed to migrate CLI after switching to 3.x channel: %w", err)
	}

	if err := migration.RunInstall(migrationConfirmed); err != nil {
		return fmt.Errorf("failed to run install on the new valet.sh CLI: %w", err)
	}

	return nil
}

func revertReleaseChannel(releaseChannelFilePath string, previousReleaseChannel string) {
	if err := os.WriteFile(releaseChannelFilePath, []byte(previousReleaseChannel), 0644); err != nil {
		utils.Printf("warning: failed to restore previous release channel %q: %s\n", previousReleaseChannel, err.Error())
	}
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

	utils.Println("\nSuccessfully switched to next channel\n")

	return runUpdate()
}
