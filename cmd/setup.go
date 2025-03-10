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
    Short: "Setup valet-sh",
    Long: `Setup valet-sh`,
    SilenceUsage: true,
    RunE: func(cmd *cobra.Command, args []string) error {
        return setupVsh()
    },
}

func init() {
}

func setupVsh() error {
    fmt.Println("Setting up valet-sh")

    vshUser, err := prechecks.GetCurrentUser()
    if err != nil {
        return err
    }

    var vshGroup string
    if goruntime.GOOS == "linux" {
        vshGroup = "vshUser"
    } else if goruntime.GOOS == "darwin" {
        vshGroup = "admin"
    }

    arch := runtime.GetArchitecture()
    homebrewPrefix := constants.HomebrewPrefix
    if arch == "darwin" && strings.HasPrefix(arch, "arm") {
        homebrewPrefix = "/opt/homebrew"
    }
    fmt.Println("Homebrew prefix:", homebrewPrefix)

    setupLogFile, err := setup.PrepareSetupLogFile()
    if err != nil {
        return err
    }
    defer setupLogFile.Close()

    if err := utils.RequestSudoAccess(); err != nil {
        return err
    }

    if goruntime.GOOS == "linux" {
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
    fmt.Println("Setting up valet-sh on Linux")

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

        fmt.Println("Valet-sh setup complete")
    }
    return nil
}

func setupMacOS(vshUser, vshGroup, homebrewPrefix string, isMacARM bool, logFile *os.File) error {
    fmt.Println("Setting up valet-sh on macOS")

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
