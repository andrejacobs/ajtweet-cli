package buildinfo

import "fmt"

var (
	GitCommitHash string
	Version       string
)

func VersionString() string {
	version := "v0.0.0"
	if Version != "" {
		version = Version
	}
	return fmt.Sprintf("%s %s", version, GitCommitHash)
}
