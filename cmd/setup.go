package cmd

import (
    "fmt"
    "os"
    "strings"
    goruntime "runtime"

    "github.com/spf13/cobra"

    "github.com/valet-sh/valet-sh-installer/constants"
    "github.com/valet-sh/valet-sh-installer/internal/prechecks"
    "github.com/valet-sh/valet-sh-installer/internal/runtime"
    "github.com/valet-sh/valet-sh-installer/internal/setup"
    "github.com/valet-sh/valet-sh-installer/internal/utils"
)

var setupCmd = &cobra.Command{
    Use:  "setup",
    Short: "Setup valet-sh and the runtime",
    Long: `Setup valet-sh and the runtime`,
    SilenceUsage: true,
    RunE: func(cmd *cobra.Command, args []string) error {
        return setupVsh()
    },
}

func init() {
}

func setupVsh() error {
    vshUser, err := prechecks.GetCurrentUser()
    if err != nil {
        return err
    }

    var vshGroup string
    if goruntime.GOOS == "linux" {
        vshGroup = vshUser
    } else if goruntime.GOOS == "darwin" {
        vshGroup = "admin"
    }

    arch := runtime.GetArchitecture()
    homebrewPrefix := constants.HomebrewPrefix
    if arch == "darwin" && strings.HasPrefix(arch, "arm") {
        homebrewPrefix = "/opt/homebrew"
    }

    setupLogFile, err := setup.PrepareSetupLogFile()
    if err != nil {
        return err
    }
    defer setupLogFile.Close()

    if goruntime.GOOS == "linux" {
        fmt.Println("Setting up valet-sh on Linux")
        return setupLinux(vshUser, vshGroup, setupLogFile)
    } else if goruntime.GOOS == "darwin" {
        isMacARM := strings.HasPrefix(arch, "arm")
        if isMacARM {
            fmt.Println("Setting up valet-sh on macOS (Apple Silicon)")
        } else {
            fmt.Println("Setting up valet-sh on macOS (Intel)")
        }
        return setupMacOS(vshUser, vshGroup, homebrewPrefix, isMacARM, setupLogFile)
    }

    return nil
}

func setupLinux(vshUser, vshGroup string, logFile *os.File) error {
    if err := utils.RequestSudoAccess(); err != nil {
        return err
    }

    shouldInstall := true
    if utils.PathExists(constants.VshBasePath) || utils.PathExists(constants.VshVenvPath) {
        fmt.Println("You already have valet-sh installed, do you want to reinstall? (y/n)")
        var response string
        fmt.Scanln(&response)
        if response != "y" {
            shouldInstall = false
        } else {
            if err := setup.RemoveVshAnsibleFactsFile(); err != nil {
                return err
            }

            if err := setup.RemoveVshRepository(); err != nil {
                return err
            }

            if err := setup.RemoveVshVenv(); err != nil {
                return err
            }
        }
    }

    if shouldInstall {
        if err := setup.InstallLinuxDependencies(logFile); err != nil {
            return err
        }

        if err := setup.PrepareVshDirectory(vshUser, vshGroup, logFile); err != nil {
            return err
        }

        if err := setup.SetupRepository(logFile); err != nil {
            return err
        }

        if err := setup.CreateSymlinks(vshUser, logFile); err != nil {
            return err
        }

        if err := runUpdate(); err != nil {
            return err
        }

        fmt.Println("Valet-sh setup complete")
    }
    return nil
}

func setupMacOS(vshUser, vshGroup, homebrewPrefix string, isMacARM bool, logFile *os.File) error {
    fmt.Println("Setting up valet-sh on macOS")

    if err := utils.RequestSudoAccess(); err != nil {
        return err
    }

    shouldInstall := true
    if utils.PathExists(constants.VshBasePath) || utils.PathExists(constants.VshVenvPath) {
        fmt.Println("You already have valet-sh installed, do you want to reinstall? (y/n)")
        var response string
        fmt.Scanln(&response)
        if response != "y" {
            shouldInstall = false
        } else {
            if err := setup.RemoveVshAnsibleFactsFile(); err != nil {
                return err
            }

            if err := setup.RemoveVshRepository(); err != nil {
                return err
            }

            if err := setup.RemoveVshVenv(); err != nil {
                return err
            }
        }
    }

    if shouldInstall {

        if isMacARM {
            if err := setup.InstallMacARMDependencies(logFile); err != nil {
                return err
            }
        }

        if err := setup.InstallMacOSDependencies(homebrewPrefix, logFile); err != nil {
            return err
        }

        if err := setup.PrepareVshDirectory(vshUser, vshGroup, logFile); err != nil {
            return err
        }

        if err := setup.SetupRepository(logFile); err != nil {
            return err
        }

        if err := setup.CreateSymlinks(vshUser, logFile); err != nil {
            return err
        }

        fmt.Println("Valet-sh setup complete")
    }
    return nil
}
