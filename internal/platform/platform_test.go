package platform

import (
	"testing"
)

func TestResolve(t *testing.T) { //nolint:funlen // table-driven test with many cases
	t.Parallel()

	tests := []struct {
		name     string
		asset    string
		version  string
		prefix   string
		platform Platform
		osMap    map[string]string
		archMap  map[string]string
		want     string
	}{
		{
			name:     "darwin arm64 basic pattern",
			asset:    "gh_{version}_{os}_{arch}.tar.gz",
			version:  "v2.87.0",
			prefix:   "v",
			platform: DarwinARM64,
			want:     "gh_2.87.0_darwin_arm64.tar.gz",
		},
		{
			name:     "linux amd64 basic pattern",
			asset:    "gh_{version}_{os}_{arch}.tar.gz",
			version:  "v2.87.0",
			prefix:   "v",
			platform: LinuxAMD64,
			want:     "gh_2.87.0_linux_amd64.tar.gz",
		},
		{
			name:     "darwin amd64 pattern",
			asset:    "tool-{version}-{os}-{arch}.zip",
			version:  "v1.0.0",
			prefix:   "v",
			platform: DarwinAMD64,
			want:     "tool-1.0.0-darwin-amd64.zip",
		},
		{
			name:     "linux arm64 pattern",
			asset:    "ripgrep-{version}-{arch}-{os}.tar.gz",
			version:  "14.1.0",
			prefix:   "v",
			platform: LinuxARM64,
			want:     "ripgrep-14.1.0-arm64-linux.tar.gz",
		},
		{
			name:     "no placeholders at all",
			asset:    "static-binary.tar.gz",
			version:  "v1.0.0",
			prefix:   "v",
			platform: LinuxAMD64,
			want:     "static-binary.tar.gz",
		},
		{
			name:     "version already without v prefix",
			asset:    "tool-{version}.tar.gz",
			version:  "1.2.3",
			prefix:   "v",
			platform: LinuxAMD64,
			want:     "tool-1.2.3.tar.gz",
		},
		{
			name:     "multiple version placeholders",
			asset:    "{version}/tool-{version}-{os}-{arch}.tar.gz",
			version:  "v3.0.0",
			prefix:   "v",
			platform: DarwinARM64,
			want:     "3.0.0/tool-3.0.0-darwin-arm64.tar.gz",
		},
		{
			name:     "empty version string",
			asset:    "tool-{version}-{os}.tar.gz",
			version:  "",
			prefix:   "v",
			platform: DarwinARM64,
			want:     "tool--darwin.tar.gz",
		},
		{
			name:     "os_map overrides OS",
			asset:    "tool_{os}_{arch}.tar.gz",
			version:  "v1.0.0",
			prefix:   "v",
			platform: DarwinARM64,
			osMap:    map[string]string{"darwin": "Darwin"},
			want:     "tool_Darwin_arm64.tar.gz",
		},
		{
			name:     "arch_map overrides arch",
			asset:    "tool_{os}_{arch}.tar.gz",
			version:  "v1.0.0",
			prefix:   "v",
			platform: LinuxAMD64,
			archMap:  map[string]string{"amd64": "x86_64"},
			want:     "tool_linux_x86_64.tar.gz",
		},
		{
			name:     "both os_map and arch_map",
			asset:    "tool_{os}_{arch}.tar.gz",
			version:  "v1.0.0",
			prefix:   "v",
			platform: DarwinAMD64,
			osMap:    map[string]string{"darwin": "macOS"},
			archMap:  map[string]string{"amd64": "x86_64"},
			want:     "tool_macOS_x86_64.tar.gz",
		},
		{
			name:     "os_map with no matching key uses default",
			asset:    "tool_{os}_{arch}.tar.gz",
			version:  "v1.0.0",
			prefix:   "v",
			platform: LinuxARM64,
			osMap:    map[string]string{"darwin": "Darwin"},
			want:     "tool_linux_arm64.tar.gz",
		},
		{
			name:     "go prefix",
			asset:    "go{version}.{os}-{arch}.tar.gz",
			version:  "go1.26.0",
			prefix:   "go",
			platform: DarwinARM64,
			want:     "go1.26.0.darwin-arm64.tar.gz",
		},
		{
			name:     "unknown platform uses raw values",
			asset:    "tool_{os}_{arch}.tar.gz",
			version:  "v1.0.0",
			prefix:   "v",
			platform: Platform{OS: "freebsd", Arch: "riscv64"},
			want:     "tool_freebsd_riscv64.tar.gz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Resolve(tt.asset, tt.version, tt.prefix, tt.platform, tt.osMap, tt.archMap)

			if got != tt.want {
				t.Errorf("Resolve() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPlatformConstants(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		platform Platform
		wantOS   string
		wantArch string
	}{
		{"DarwinARM64", DarwinARM64, "darwin", "arm64"},
		{"DarwinAMD64", DarwinAMD64, "darwin", "amd64"},
		{"LinuxARM64", LinuxARM64, "linux", "arm64"},
		{"LinuxAMD64", LinuxAMD64, "linux", "amd64"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.platform.OS != tt.wantOS {
				t.Errorf("OS = %q, want %q", tt.platform.OS, tt.wantOS)
			}

			if tt.platform.Arch != tt.wantArch {
				t.Errorf("Arch = %q, want %q", tt.platform.Arch, tt.wantArch)
			}
		})
	}
}
