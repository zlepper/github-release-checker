package github_release_checker

import (
	"github.com/blang/semver"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetReleases(t *testing.T) {
	a := assert.New(t)

	releases, err := GetReleases("zlepper", "gfs")
	if a.NoError(err) {
		a.True(len(releases) > 2)
	}
}

func TestGetLatestReleaseForPlatform(t *testing.T) {
	a := assert.New(t)

	latestRelease, err := GetLatestReleaseForPlatform("zlepper", "gfs", "gfs-windows-x64", true)
	if a.NoError(err) {
		a.Equal("0.0.2", latestRelease.TagName)
	}
}

func TestIsNewer(t *testing.T) {
	a := assert.New(t)

	result, err := IsNewer(Release{semverTag: semver.MustParse("0.0.3")}, "0.0.2")
	if a.NoError(err) {
		a.True(result)
	}

	result, err = IsNewer(Release{semverTag: semver.MustParse("0.0.1")}, "0.0.2")
	if a.NoError(err) {
		a.False(result)
	}
}
