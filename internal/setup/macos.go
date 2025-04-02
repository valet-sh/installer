package setup

import (
	"fmt"
	"os"

	"github.com/valet-sh/valet-sh-installer/internal/utils"
)

func InstallMacARMDependencies(logFile *os.File) error {
	utils.Println("Installing dependencies for Mac ARM")
	utils.Println("Install rosetta")
	if err := utils.RunCommand("/usr/sbin/softwareupdate", []string{"--install-rosetta", "--agree-to-license"}, logFile); err != nil {
		return fmt.Errorf("failed to install rosetta: %w", err)
	}
	return nil
}

func InstallMacOSDependencies(homebrewPrefix string, logFile *os.File) error {
	utils.Println("Installing dependencies for macOS")

	if err := utils.RunCommand("curl", []string{"-fsSL", "https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh", "-o", "/tmp/homebrew_install.sh"}, logFile); err != nil {
		return fmt.Errorf("failed to download homebrew install script: %w", err)
	}
	if err := utils.RunCommand("/bin/bash", []string{"/tmp/homebrew_install.sh"}, logFile); err != nil {
		return fmt.Errorf("failed to install homebrew: %w", err)
	}

	// @FIXME
	os.Setenv("CPPFLAGS", "-I"+homebrewPrefix+"/opt/openssl/include")
	os.Setenv("LDFLAGS", "-L"+homebrewPrefix+"/opt/openssl/lib")

	utils.Println(" - installing required brew packages")
	if err := utils.RunCommand(homebrewPrefix+"/bin/brew", []string{"install", "openssl", "rust", "python@3.12"}, logFile); err != nil {
		return fmt.Errorf("failed to install required brew packages: %w", err)
	}

	utils.Println(" - initializing brew services")
	if err := utils.RunCommand(homebrewPrefix+"/bin/brew", []string{"services", "list"}, logFile); err != nil {
		return fmt.Errorf("failed to initialize brew services: %w", err)
	}

	return nil
}
