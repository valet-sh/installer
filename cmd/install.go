package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
    Use:  "install",
    Short: "Install valet-sh",
    Long: `Install valet-sh`,
    SilenceUsage: true,
    RunE: func(cmd *cobra.Command, args []string) error {
        return install()
    },
}

func init() {
}

func install() error {
    fmt.Println("install")
    return nil
}
