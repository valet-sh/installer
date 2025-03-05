package main

import (
    "os"

    "github.com/valet-sh/valet-sh-installer/cmd"
)

func main() {
    if err := cmd.Execute(); err != nil {
        os.Exit(1)
    }
}
