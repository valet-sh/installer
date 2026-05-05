package setup

import (
	"fmt"
	"os"

	"github.com/valet-sh/valet-sh-installer/internal/utils"
)

func InstallMacARMDependencies(homebrewPrefix string) error {
	utils.Println("Installing dependencies for Mac ARM")

	InstallMacOSHomebrew()
	InstallMacOSRosetta()
	InstallMacOSHomebrewPackages(homebrewPrefix)

	return nil
}

func InstallMacOSDependencies(homebrewPrefix string) error {
	utils.Println("Installing dependencies for macOS")

	InstallMacOSHomebrew()
	InstallMacOSHomebrewPackages(homebrewPrefix)

	return nil
}

func InstallMacOSHomebrew() error {
	if err := utils.RunCommand("curl", []string{"-fsSL", "https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh", "-o", "/tmp/homebrew_install.sh"}); err != nil {
		return fmt.Errorf("failed to download homebrew install script: %w", err)
	}
	if err := utils.RunCommand("/bin/bash", []string{"/tmp/homebrew_install.sh"}); err != nil {
		return fmt.Errorf("failed to install homebrew: %w", err)
	}

	return nil
}

func InstallMacOSRosetta() error {
	if err := utils.RunCommand("/usr/sbin/softwareupdate", []string{"--install-rosetta", "--agree-to-license"}); err != nil {
		return fmt.Errorf("failed to install rosetta: %w", err)
	}

	return nil
}

func InstallMacOSHomebrewPackages(homebrewPrefix string) error {
	// @FIXME
	os.Setenv("CPPFLAGS", "-I"+homebrewPrefix+"/opt/openssl/include")
	os.Setenv("LDFLAGS", "-L"+homebrewPrefix+"/opt/openssl/lib")

	utils.Println(" - installing required brew packages")
	if err := utils.RunCommand(homebrewPrefix+"/bin/brew", []string{"install", "openssl", "rust", "python@3.12"}); err != nil {
		return fmt.Errorf("failed to install required brew packages: %w", err)
	}

	utils.Println(" - initializing brew services")
	if err := utils.RunCommand(homebrewPrefix+"/bin/brew", []string{"services", "list"}); err != nil {
		return fmt.Errorf("failed to initialize brew services: %w", err)
	}

	return nil
}
