// Package platform provides OS/architecture resolution for expanding asset
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

// resolveOS returns the OS name to use. If osMap contains a mapping for
// p.OS, the mapped value is returned; otherwise p.OS is used directly.
func resolveOS(p Platform, osMap map[string]string) string {
	if osMap != nil {
		if v, ok := osMap[p.OS]; ok {
			return v
		}
	}

	return p.OS
}

// resolveArch returns the arch name to use. If archMap contains a mapping for
// p.Arch, the mapped value is returned; otherwise p.Arch is used directly.
func resolveArch(p Platform, archMap map[string]string) string {
	if archMap != nil {
		if v, ok := archMap[p.Arch]; ok {
			return v
		}
	}

	return p.Arch
}

// Resolve expands an asset pattern into a single concrete asset name
// by substituting {version}, {os}, and {arch} placeholders.
//
// The version string has a leading "v" stripped before substitution.
// osMap and archMap allow per-tool overrides of the default GOOS/GOARCH values.
func Resolve(asset, version, prefix string, p Platform, osMap, archMap map[string]string) string {
	ver := strings.TrimPrefix(version, prefix)
	r := strings.ReplaceAll(asset, "{version}", ver)
	r = strings.ReplaceAll(r, "{os}", resolveOS(p, osMap))
	r = strings.ReplaceAll(r, "{arch}", resolveArch(p, archMap))

	return r
}
