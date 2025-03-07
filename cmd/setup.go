package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
    Use:  "setup",
    Short: "Setup valet-sh",
    Long: `Setup valet-sh`,
    SilenceUsage: true,
    RunE: func(cmd *cobra.Command, args []string) error {
        return setup()
    },
}

func init() {
}

func setup() error {
    fmt.Println("setup")
    return nil
}
