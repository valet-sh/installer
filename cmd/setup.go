package cmd

import (
	"fmt"
	"os"
	goruntime "runtime"
	"strings"

	"github.com/gookit/color"

	"github.com/spf13/cobra"

	"github.com/valet-sh/valet-sh-installer/constants"
	"github.com/valet-sh/valet-sh-installer/internal/prechecks"
	"github.com/valet-sh/valet-sh-installer/internal/runtime"
	"github.com/valet-sh/valet-sh-installer/internal/setup"
	"github.com/valet-sh/valet-sh-installer/internal/utils"
)

var setupCmd = &cobra.Command{
	Use:           "setup",
	Short:         "Setup valet-sh and the runtime",
	Long:          `Setup valet-sh and the runtime`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := setupVsh()
		if err != nil {
			color.Error.Printf("Error: %s\n", err.Error())
			return err
		}
		return nil
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
	if goruntime.GOOS == "darwin" && strings.HasPrefix(arch, "arm") {
		homebrewPrefix = "/opt/homebrew"
	}

	setupLogFile, err := setup.PrepareSetupLogFile()
	if err != nil {
		return err
	}

	if goruntime.GOOS == "linux" {
		color.Info.Println("Setting up valet-sh on Linux\n")
		return setupLinux(vshUser, vshGroup, setupLogFile)
	} else if goruntime.GOOS == "darwin" {
		isMacARM := strings.HasPrefix(arch, "arm")
		if isMacARM {
			color.Info.Println("Setting up valet-sh on macOS (Apple Silicon)\n")
		} else {
			color.Info.Println("Setting up valet-sh on macOS (Intel)\n")
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
		fmt.Println("You already have valet-sh installed, do you want to reinstall? (y/N)")
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

		color.Greenln("Valet-sh setup complete")
	} else {
		color.Warn.Println("Setup cancelled")
	}
	return nil
}

func setupMacOS(vshUser, vshGroup, homebrewPrefix string, isMacARM bool, logFile *os.File) error {
	if err := utils.RequestSudoAccess(); err != nil {
		return err
	}

	shouldInstall := true
	if utils.PathExists(constants.VshBasePath) || utils.PathExists(constants.VshVenvPath) {
		fmt.Println("You already have valet-sh installed, do you want to reinstall? (y/N)")
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
			if err := setup.InstallMacARMDependencies(homebrewPrefix, logFile); err != nil {
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

		if err := runUpdate(); err != nil {
			return err
		}

		color.Greenln("Valet-sh setup complete")
	} else {
		color.Warn.Println("Setup cancelled")
	}
	return nil
}
