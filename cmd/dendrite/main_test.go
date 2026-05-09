package main

import (
	"testing"

	"github.com/sivchari/dendrite/internal/config"
	"github.com/sivchari/dendrite/internal/platform"
)

func TestResolveCandidates(t *testing.T) {
	t.Parallel()

	tool := &config.Tool{
		Owner:         "cli",
		Repo:          "cli",
		Version:       "v2.87.0",
		VersionPrefix: "v",
		Asset: map[string]string{
			"darwin/arm64": "gh_{version}_{os}_{arch}.zip",
		},
	}

	got := resolveCandidates(tool, platform.DarwinARM64)

	if len(got) == 0 {
		t.Fatal("resolveCandidates() returned no candidates")
	}

	wantFirst := "https://github.com/cli/cli/releases/download/v2.87.0/gh_2.87.0_darwin_arm64.zip"
	if got[0] != wantFirst {
		t.Errorf("resolveCandidates()[0] = %q, want %q", got[0], wantFirst)
	}

	wantFallback := "https://github.com/cli/cli/releases/download/v2.87.0/gh_2.87.0_darwin_arm64.tar.gz"
	if !containsString(got, wantFallback) {
		t.Errorf("resolveCandidates() does not contain %q\ngot: %#v", wantFallback, got)
	}

	wantAlias := "https://github.com/cli/cli/releases/download/v2.87.0/gh_2.87.0_macOS_arm64.zip"
	if !containsString(got, wantAlias) {
		t.Errorf("resolveCandidates() does not contain %q\ngot: %#v", wantAlias, got)
	}
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}

	return false
}
