package cmd

import (
    "fmt"
    "log"
    "errors"
    "context"
    "runtime"


    "github.com/spf13/cobra"
    "github.com/creativeprojects/go-selfupdate"
)

var selfUpgradeCmd = &cobra.Command{
    Use:  "self-upgrade",
    Short: "Update valet-sh-installer to the latest version",
    Long: `Update valet-sh-installer to the latest version`,
    SilenceUsage: true,
    RunE: func(cmd *cobra.Command, args []string) error {
        return selfUpgrade(rootCmd.Version)
    },
}

func init() {
}

func selfUpgrade(version string) error {
    latest, found, err := selfupdate.DetectLatest(context.Background(), selfupdate.ParseSlug("valet-sh/valet-sh-installer"))
    fmt.Println(latest)
    if err != nil {
        return fmt.Errorf("error occurred while detecting version: %w", err)
    }
    if !found {
        return fmt.Errorf("latest version for %s/%s could not be found from github repository", runtime.GOOS, runtime.GOARCH)
    }

    if latest.LessOrEqual(version) {
        log.Printf("Current version (%s) is the latest", version)
        return nil
    }

    exe, err := selfupdate.ExecutablePath()
    if err != nil {
        return errors.New("could not locate executable path")
    }
    if err := selfupdate.UpdateTo(context.Background(), latest.AssetURL, latest.AssetName, exe); err != nil {
        return fmt.Errorf("error occurred while updating binary: %w", err)
    }
    log.Printf("Successfully updated to version %s", latest.Version())
    return nil
}
