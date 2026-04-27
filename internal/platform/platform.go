// Package platform provides version resolution for expanding asset
// pattern placeholders into a concrete download URL.
package platform

import "strings"

// Platform represents a target OS and architecture combination.
type Platform struct {
	OS   string
	Arch string
}

// Predefined platform constants for common targets.
var (
	DarwinARM64 = Platform{OS: "darwin", Arch: "arm64"}
	DarwinAMD64 = Platform{OS: "darwin", Arch: "amd64"}
	LinuxARM64  = Platform{OS: "linux", Arch: "arm64"}
	LinuxAMD64  = Platform{OS: "linux", Arch: "amd64"}
)

// Resolve expands an asset pattern by substituting the {version} placeholder.
// The version string has the given prefix stripped before substitution.
func Resolve(pattern, version, prefix string) string {
	ver := strings.TrimPrefix(version, prefix)

	return strings.ReplaceAll(pattern, "{version}", ver)
}
