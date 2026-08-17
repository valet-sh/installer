package migration

import (
	"fmt"
	"os"

	"github.com/valet-sh/valet-sh-installer/constants"
	"github.com/valet-sh/valet-sh-installer/internal/git"
)

type Snapshot struct {
	RepoPath       string
	PreviousCommit string
	SymlinkExisted bool
	SymlinkTarget  string
	VenvBackupPath string
}

func CaptureState(repoPath string) (*Snapshot, error) {
	commit, err := git.GetCurrentCommit(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to determine current valet-sh checkout: %w", err)
	}

	snapshot := &Snapshot{
		RepoPath:       repoPath,
		PreviousCommit: commit,
	}

	if target, err := os.Readlink(constants.VshCliSymlinkPath); err == nil {
		snapshot.SymlinkExisted = true
		snapshot.SymlinkTarget = target
	}

	return snapshot, nil
}

func Rollback(snapshot *Snapshot, releaseChannelFilePath, previousReleaseChannel string) []error {
	var errs []error

	if err := os.WriteFile(releaseChannelFilePath, []byte(previousReleaseChannel), 0644); err != nil {
		errs = append(errs, fmt.Errorf("failed to restore previous release channel %q: %w", previousReleaseChannel, err))
	}

	if err := git.CheckoutBranch(snapshot.RepoPath, snapshot.PreviousCommit); err != nil {
		errs = append(errs, fmt.Errorf("failed to restore previous valet-sh checkout %q: %w", snapshot.PreviousCommit, err))
	}

	if snapshot.VenvBackupPath != "" {
		if err := restoreVenvBackup(snapshot.VenvBackupPath); err != nil {
			errs = append(errs, fmt.Errorf("failed to restore previous runtime: %w", err))
		}
	}

	if snapshot.SymlinkExisted {
		if err := UpdateSymlink(snapshot.SymlinkTarget); err != nil {
			errs = append(errs, fmt.Errorf("failed to restore previous CLI symlink: %w", err))
		}
	}

	return errs
}

func restoreVenvBackup(backupPath string) error {
	if _, err := os.Stat(backupPath); err != nil {
		return fmt.Errorf("runtime backup %q not found: %w", backupPath, err)
	}

	if err := os.RemoveAll(constants.VshVenvPath); err != nil {
		return err
	}

	if err := os.Rename(backupPath, constants.VshVenvPath); err != nil {
		return err
	}

	return nil
}
