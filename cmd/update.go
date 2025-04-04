package cmd

import (
	"fmt"
	"github.com/gookit/color"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/valet-sh/valet-sh-installer/constants"
	"github.com/valet-sh/valet-sh-installer/internal/git"
	"github.com/valet-sh/valet-sh-installer/internal/runtime"
	"github.com/valet-sh/valet-sh-installer/internal/utils"
)

var updateCmd = &cobra.Command{
	Use:           "update",
	Short:         "Update valet-sh and the runtime to the latest version",
	Long:          `Update valet-sh and the runtime to the latest version`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := runUpdate()
		if err != nil {
			color.Error.Prompt(err.Error())
			return err
		}
		return nil
	},
}

func init() {
}

func runUpdate() error {
	repoPath := constants.VshBasePath
	if err := checkIfRepoExists(repoPath); err != nil {
		return err
	}

	releaseChannel := getCurrentReleaseChannel()

	if releaseChannel == "next" {
		utils.Println("Using next channel (development) for update")

		return updateNextBranch(repoPath)
	} else if strings.HasSuffix(releaseChannel, ".x") {
		majorVersion := strings.Split(releaseChannel, ".")[0]
		utils.Printf("Using %s channel for update\n", releaseChannel)

		return updateVersionBranch(repoPath, releaseChannel, majorVersion)
	} else {
		return fmt.Errorf("invalid release channel: %s", releaseChannel)
	}
}

func checkIfRepoExists(repoPath string) error {
	_, err := os.Stat(repoPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("valet-sh not found in %s", repoPath)
	}
	return nil
}

func updateNextBranch(repoPath string) error {
	utils.Println("Updating valet-sh to the latest version on the next branch")

	if err := git.CheckoutBranch(repoPath, "next"); err != nil {
		return fmt.Errorf("failed to checkout next branch: %w", err)
	}

	if err := git.PullLatest(repoPath); err != nil {
		return fmt.Errorf("failed to pull latest changes: %w", err)
	}

	color.Info.Println("valet-sh: Successfully pulled latest changes from next branch")

	return runtimeUpdate()
}

func updateVersionBranch(repoPath string, branchName string, majorVersion string) error {
	utils.Printf("Updating valet-sh to the latest version on the %s branch\n", branchName)

	if err := git.FetchTags(repoPath); err != nil {
		return fmt.Errorf("failed to fetch tags: %w", err)
	}

	currentRelease, err := git.GetCurrentReleaseTag(repoPath)
	if err != nil {
		return fmt.Errorf("failed to get current release: %w", err)
	}

	tags, err := git.GetAllTags(repoPath)
	if err != nil {
		return fmt.Errorf("failed to get all tags: %w", err)
	}

	if len(tags) == 0 {
		return fmt.Errorf("no valid releases found")
	}

	semverRegex := buildSemverRegex(majorVersion)

	validVersions := git.FilterTagsSemver(tags, semverRegex)

	if len(validVersions) == 0 {
		utils.Printf("No valid releases found for %s channel\n", branchName)

		utils.Printf("Do you want to switch to the %s branch for testing without a release? (y/n)", branchName)

		var response string
		fmt.Scanln(&response)

		if response != "y" {
			return nil
		} else {
			branchExists, err := git.DoesBranchExist(repoPath, branchName)
			if err != nil {
				return fmt.Errorf("error checking if branch %s exists: %w", branchName, err)
			}

			if !branchExists {
				return fmt.Errorf("release channel %s does not exist - please select a valid release channel", branchName)
			}

			utils.Printf("Switching to %s branch for testing without a release\n", branchName)

			if err := git.CheckoutBranch(repoPath, branchName); err != nil {
				return fmt.Errorf("failed to checkout %s branch: %w", branchName, err)
			}

			if err := git.PullLatest(repoPath); err != nil {
				return fmt.Errorf("failed to pull latest changes: %w", err)
			}

			utils.Printf("Successfully switched to %s branch for testing\n", branchName)

			return runtimeUpdate()
		}
	}

	latestVersion := validVersions[0]
	latestTag := latestVersion
	if !strings.HasPrefix(tags[0], "v") && strings.HasPrefix(latestVersion, "v") {
		latestTag = latestVersion[1:]
	}

	utils.Printf("Latest version available: %s | current: %s\n", latestTag, currentRelease)

	compareResult := utils.CompareVersions(currentRelease, latestTag)

	if compareResult < 0 {
		utils.Printf("Updating valet-sh from version %s to %s\n", currentRelease, latestTag)
		if err := git.CheckoutBranch(repoPath, latestTag); err != nil {
			return fmt.Errorf("failed to checkout version %s: %w", latestTag, err)
		}

		color.Info.Printf("valet-sh successfully updated to version %s\n\n", latestTag)
	} else {
		color.Info.Printf("Already on latest version %s\n", currentRelease)
	}

	return runtimeUpdate()
}

func buildSemverRegex(majorVersion string) string {
	return fmt.Sprintf(`^(%s)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(\-[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$`, majorVersion)
}

func runtimeUpdate() error {
	status, err := runtime.CheckRuntime()
	if err != nil {
		return fmt.Errorf("failed to check runtime: %w", err)
	}

	if status.NeedsUpdate || status.PackageChanged {

		fmt.Printf("Updating valet-sh runtime to version %s\n", status.CurrentVersion)

		utils.Println("Updating runtime")

		url := fmt.Sprintf("https://github.com/valet-sh/runtime/releases/download/%s/%s.tar.gz", status.CurrentVersion, status.CurrentPackage)

		utils.Printf("Check if runtime release '%s' exists\n", url)
		resp, err := http.Head(url)
		if err != nil {
			return fmt.Errorf("failed to check runtime release: %w", err)
		}

		// @FIXME
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("runtime release not found: %s", url)
		}

		tmpDir := constants.VshVenvTmpPath
		utils.Println("Cleaning up temporary directory")

		err = os.RemoveAll(tmpDir)
		if err != nil {
			return fmt.Errorf("failed to remove temporary directory: %w", err)
		}

		venvDir := constants.VshVenvPath
		_, err = os.Stat(venvDir)
		venvExists := !os.IsNotExist(err)

		if venvExists {
			utils.Println("Backing up current venv")
			err = os.Rename(venvDir, tmpDir)
			if err != nil {
				return fmt.Errorf("failed to backup current venv: %w", err)
			}
		}

		utils.Printf("Downloading and unpacking new runtime '%s' ", status.CurrentVersion)
		respDownload, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to download runtime release: %w", err)
		}
		defer respDownload.Body.Close()

		if respDownload.StatusCode != http.StatusOK {
			return fmt.Errorf("bad status code:  %s", respDownload.Status)
		}

		err = utils.Untar(constants.VshPath, respDownload.Body)
		if err != nil {
			if venvExists {
				err = os.RemoveAll(venvDir)
				if err != nil {
					return fmt.Errorf("failed to remove venv directory: %w", err)
				}
				err = os.Rename(tmpDir, venvDir)
				if err != nil {
					return fmt.Errorf("failed to move venv directory: %w", err)
				}
			}
			return fmt.Errorf("failed to extract runtime: %w", err)
		}

		venvVersion := status.CurrentPackage + "-" + status.CurrentVersion
		err = os.WriteFile(filepath.Join(constants.VshVenvPath, constants.VersionFileName),
			[]byte(venvVersion), 0644)
		if err != nil {
			return fmt.Errorf("failed to update version file: %w", err)
		}

		color.Info.Printf("\nRuntime '%s' updated successfully\n\n", status.CurrentVersion)

		if venvExists {
			err = os.RemoveAll(tmpDir)
			if err != nil {
				return fmt.Errorf("failed to remove temporary directory: %w", err)
			}
		}

	} else {
		color.Info.Println("valet-sh runtime: already up to date\n")
	}

	return nil
}
