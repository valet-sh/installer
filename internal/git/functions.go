package git

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/valet-sh/valet-sh-installer/constants"
	"golang.org/x/mod/semver"
)

const (
	apiTimeout = 3 * time.Second
)

type releaseResponse struct {
	TagName string `json:"tag_name"`
}

func CheckoutBranch(dir, branch string) error {
	cmd := exec.Command("git", "checkout", "--force", "--quiet", branch)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func PullLatest(dir string) error {
	cmd := exec.Command("git", "pull", "--quiet")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func FetchTags(dir string) error {
	cmd := exec.Command("git", "fetch", "--tags", "--quiet")
	cmd.Dir = dir
	return cmd.Run()
}

func GetCurrentReleaseTag(dir string) (string, error) {
	cmd := exec.Command("git", "describe", "--tags")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func FilterTagsSemver(tags []string, pattern string) []string {
	re := regexp.MustCompile(pattern)
	var filtered []string

	for _, tag := range tags {
		if re.MatchString(tag) {
			vTag := tag
			if !strings.HasPrefix(vTag, "v") {
				vTag = "v" + vTag
			}
			filtered = append(filtered, vTag)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return semver.Compare(filtered[i], filtered[j]) > 0
	})

	return filtered
}

func GetAllTags(dir string) ([]string, error) {
	cmd := exec.Command("git", "tag")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	tags := strings.Split(string(out), "\n")
	return tags[:len(tags)-1], nil
}

func CompareReleases(current, latest string) bool {
	return semver.Compare(current, latest) < 0
}

func CloneRepository(repoUrl, repoPath string) error {
	cmd := exec.Command("git", "clone", "--quiet", repoUrl, repoPath)
	return cmd.Run()
}

func DoesBranchExist(repoPath, branchName string) (bool, error) {
	cmd := exec.Command("git", "branch", "--list", "--all")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	branches := strings.Split(string(output), "\n")
	for _, branch := range branches {
		cleanBranch := strings.TrimSpace(branch)
		if strings.HasPrefix(cleanBranch, "*") {
			cleanBranch = strings.TrimSpace(cleanBranch[1:])
		}
		if cleanBranch == branchName ||
			cleanBranch == "remotes/origin/"+branchName ||
			strings.HasSuffix(cleanBranch, "/"+branchName) {
			return true, nil
		}
	}

	return false, nil
}

func FetchLatestCliTag(timeout time.Duration) (string, error) {
	resp, err := githubGet(constants.VshCliGithubRepoUrl, timeout)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var rel releaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return "", err
	}
	return strings.TrimPrefix(rel.TagName, ""), nil
}

func githubGet(url string, timeout time.Duration) (*http.Response, error) {
	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	return client.Do(req)
}
