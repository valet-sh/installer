package cmd

import (
	"context"
	"errors"
	"fmt"
	"runtime"

	"github.com/gookit/color"
	"github.com/valet-sh/valet-sh-installer/internal/utils"

	"github.com/creativeprojects/go-selfupdate"
	"github.com/spf13/cobra"
)

var selfUpgradeCmd = &cobra.Command{
	Use:           "self-upgrade",
	Short:         "Update valet-sh-installer to the latest version",
	Long:          `Update valet-sh-installer to the latest version`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := selfUpgrade(rootCmd.Version)
		if err != nil {
			color.Error.Printf("Error: %s\n", err.Error())
			return err
		}
		return nil
	},
}

func init() {
}

func selfUpgrade(version string) error {
	latest, found, err := selfupdate.DetectLatest(context.Background(), selfupdate.ParseSlug("valet-sh/installer"))

	utils.Println(latest)

	if err != nil {
		return fmt.Errorf("error occurred while detecting version: %w", err)
	}
	if !found {
		return fmt.Errorf("latest version for %s/%s could not be found from github repository", runtime.GOOS, runtime.GOARCH)
	}

	if latest.LessOrEqual(version) {
		color.Info.Printf("valet-sh installer: Current version (%s) is the latest\n", version)
		return nil
	}

	exe, err := selfupdate.ExecutablePath()
	if err != nil {
		return errors.New("could not locate executable path")
	}
	if err := selfupdate.UpdateTo(context.Background(), latest.AssetURL, latest.AssetName, exe); err != nil {
		return fmt.Errorf("error occurred while updating binary: %w", err)
	}
	color.Info.Printf("valet-sh installer: Successfully updated to version %s\n", latest.Version())
	return nil
}
