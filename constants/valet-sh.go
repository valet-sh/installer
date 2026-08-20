package constants

const (
	VshBasePath        = "/usr/local/valet-sh/valet-sh"
	VshVenvPath        = "/usr/local/valet-sh/venv"
	VshEtcPath         = "/usr/local/valet-sh/etc"
	VshVenvTmpPath     = "/usr/local/valet-sh/venv-tmp"
	VshPath            = "/usr/local/valet-sh"
	VshInstallerPath   = "/usr/local/valet-sh/installer/valet-sh-installer"
	ReleaseChannelFile = "RELEASE_CHANNEL"
	RuntimeFileName    = ".runtime_version"
	VersionFileName    = ".version"
	ValetStableVersion = "2.x"

	VshServiceFile      = VshEtcPath + "/services.yml"
	VshBundlesFile      = VshEtcPath + "/bundles.yml"
	VshMigrationFile    = VshEtcPath + "/migration.yml"
	VshUrl              = "https://valet.sh"
	HomebrewPrefix      = "/usr/local"
	VshInstallLog       = "/tmp/valet-sh-install.log"
	VshGithubRepoUrl    = "https://github.com/valet-sh/valet-sh"
	VshAnsibleFactsFile = "/tmp/ansible-facts/local"
	VshCliGithubRepoUrl = "https://api.github.com/repos/" + "valet-sh" + "/go-cli" + "/releases/latest"
	VshOldCliPath       = "/usr/local/bin/valet-sh"
	VshCliPath          = "/usr/local/valet-sh/bin"
	VshCliBinaryPath    = "/usr/local/valet-sh/bin/valet-sh"
	VshCliSymlinkPath   = "/usr/local/bin/valet.sh"

	Vsh3xMinMacOSVersion = "26.0"
	Vsh3xMinLinuxVersion = "24.04"
)
