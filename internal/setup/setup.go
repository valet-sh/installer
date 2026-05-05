package setup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/valet-sh/valet-sh-installer/constants"
	"github.com/valet-sh/valet-sh-installer/internal/git"
	"github.com/valet-sh/valet-sh-installer/internal/utils"
)

func PrepareSetupLogFile() error {
	setupLogFile, err := os.Create(constants.VshInstallLog)
	if err != nil {
		return fmt.Errorf("failed to create setup log file: %w", err)
	}

	utils.LogFile = setupLogFile

	return nil
}

func PrepareVshDirectory(vshUser, vshGroup string) error {
	if _, err := os.Stat(constants.VshPath); os.IsNotExist(err) {
		if err := utils.RunCommand("sudo", []string{"mkdir", "-p", constants.VshPath}); err != nil {
			return fmt.Errorf("failed to create install directory: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check if directory exists: %w", err)
	}

	if err := utils.RunCommand("sudo", []string{"chown", fmt.Sprintf("%s:%s", vshUser, vshGroup), constants.VshPath}); err != nil {
		return fmt.Errorf("failed to set permissions on install directory: %w", err)
	}

	return nil
}

func RemoveVshAnsibleFactsFile() error {
	if _, err := os.Stat(constants.VshAnsibleFactsFile); err == nil {
		_ = os.Remove(constants.VshAnsibleFactsFile)
	}

	return nil
}

func RemoveVshRepository() error {
	if _, err := os.Stat(constants.VshBasePath); err == nil {
		utils.Println("Removing existing repository...")
		if err := os.RemoveAll(constants.VshBasePath); err != nil {
			return fmt.Errorf("failed to remove existing repository: %w", err)
		}
	}

	return nil
}

func RemoveVshVenv() error {
	if _, err := os.Stat(constants.VshVenvPath); err == nil {
		utils.Println("Removing existing virtual environment...")
		if err := os.RemoveAll(constants.VshVenvPath); err != nil {
			return fmt.Errorf("failed to remove existing runtime: %w", err)
		}
	}

	return nil
}

func SetupVshRepository() error {
	if _, err := os.Stat(filepath.Join(constants.VshBasePath, ".git")); os.IsNotExist(err) {
		if err := git.CloneRepository(constants.VshGithubRepoUrl, constants.VshBasePath); err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}
	} else {
		if err := utils.RunCommand("git", []string{"-C", constants.VshBasePath, "pull"}); err != nil {
			return fmt.Errorf("failed to update repository: %w", err)
		}
	}

	return nil
}
