package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/charmbracelet/huh"

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
    branch, err := determineBranch()
    if err != nil {
     return err
    }
    fmt.Println("Updating valet-sh to the latest version on the " + branch + " branch")
    return nil
}

func determineBranch() (string, error) {
    if branchFlag == "next" || branchFlag == "stable" {
        return branchFlag, nil
    }

    _, err := os.Stat("/usr/local/valet-sh/etc/ENABLE_NEXT")
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
