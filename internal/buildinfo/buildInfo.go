package buildinfo

import "fmt"

var (
	GitCommitHash string
	Version       string
)

func VersionString() string {
	return fmt.Sprintf("%s %s", Version, GitCommitHash)
}
