// Package version provides build time version information.
package version

import (
	"fmt"
	"runtime"
)

// The following fields are populated at buildtime with bazel's linkstamp
// feature. This is equivalent to using golang directly with -ldflags -X.
var (
	buildVersion   string
	buildTime      string
	buildGitCommit string
	buildGitBranch string
)

// BuildInfo describes version information about the binary build.
type BuildInfo struct {
	Version   string `json:"version"`
	GitCommit string `json:"commit"`
	GitBranch string `json:"branch"`
	BuildTime string `json:"build_time"`
	GoVersion string `json:"go_version"`
}

func (b *BuildInfo) String() string {
	return fmt.Sprintf("%v %v", b.Version, b.GitCommit)
}

var (
	// Info exports the build version information.
	Info BuildInfo
)

func init() {
	Info.Version = buildVersion
	Info.GitCommit = buildGitCommit
	Info.GitBranch = buildGitBranch
	Info.BuildTime = buildTime
	Info.GoVersion = runtime.Version()
}

// Version returns a multi-line version information
func Version() string {
	return fmt.Sprintf("Version: \033[32m%v\033[0m\nGitCommit: \033[33m%v\033[0m\nGitBranch: \033[36m%v\033[0m\nGoVersion: %v\nBuildTime: %v\n",
		Info.Version,
		Info.GitCommit,
		Info.GitBranch,
		Info.GoVersion,
		Info.BuildTime)
}
