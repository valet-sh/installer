package setup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/valet-sh/valet-sh-installer/constants"
	"github.com/valet-sh/valet-sh-installer/internal/git"
	"github.com/valet-sh/valet-sh-installer/internal/utils"
)

func SetupRepository() error {
	if _, err := os.Stat(constants.VshBasePath); err == nil {
		utils.Println("Removing existing repository2...")
		if err := os.RemoveAll(constants.VshBasePath); err != nil {
			return fmt.Errorf("failed to remove existing repository: %w", err)
		}
	}

	if _, err := os.Stat(filepath.Join(constants.VshBasePath, ".git")); os.IsNotExist(err) {
		if err := git.CloneRepository(constants.VshGithubRepoUrl, constants.VshBasePath); err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}
	}
	return nil
}
