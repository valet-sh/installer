package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCommand(t *testing.T) {
	cmd := rootCmd
	var output bytes.Buffer
	cmd.SetOut(&output)
	err := cmd.Execute()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expectedOutput := "A CLI tool to install/update valet-sh and the runtime"
	if !strings.Contains(output.String(), expectedOutput) {
		t.Errorf("Expected output to contain '%s', got '%s'", expectedOutput, output.String())
	}
}

func TestSetupCommand(t *testing.T) {
	cmd := rootCmd
	cmd.SetArgs([]string{"setup"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestUpdateCommand(t *testing.T) {
	cmd := rootCmd
	cmd.SetArgs([]string{"update"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSelfUpgradeCommand(t *testing.T) {
	cmd := rootCmd
	cmd.Version = "0.0.12"
	cmd.SetArgs([]string{"self-upgrade"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestReleaseChannel(t *testing.T) {
	cmd := rootCmd
	cmd.SetArgs([]string{"release-channel", "next"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
