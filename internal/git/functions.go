package git

import (
    "os"
    "os/exec"
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
