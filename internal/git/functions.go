package git

import (
    "os"
    "os/exec"
    "regexp"
    "sort"
    "strings"

    "golang.org/x/mod/semver"
)

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
    cmd := exec.Command("git", "branch", "--list", "--all", branchName)
    cmd.Dir = repoPath
    output, err := cmd.Output()
    if err != nil {
        return false, err
    }

    return len(strings.TrimSpace(string(output))) > 0, nil
}
