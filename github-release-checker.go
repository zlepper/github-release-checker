package github_release_checker

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blang/semver"
	"log"
	"net/http"
	"regexp"
	"sort"
)

// A release made on GitHub
type Release struct {
	// Indicates if this release is a prerelease
	PreRelease bool
	// The version tag
	TagName string
	// The name of the file on release
	Filename string
	// The url to download the file
	DownloadUrl string
	// The size of the file
	Size      int64
	semverTag semver.Version
}

type githubRelease struct {
	TagName    string `json:"tag_name"`
	Prerelease bool   `json:"prerelease"`
	Assets     []struct {
		Name               string `json:"name"`
		Size               int64  `json:"size"`
		BrowserDownloadUrl string `json:"browser_download_url"`
	} `json:"assets"`
}

// A sortable array. Useful for sorting the releases
type semverReleases []Release

func (r semverReleases) Len() int           { return len(r) }
func (r semverReleases) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r semverReleases) Less(i, j int) bool { return r[i].semverTag.LT(r[j].semverTag) }

// Gets all the latest releases from GitHub
func GetReleases(username, repository string) (releases []Release, err error) {
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", username, repository))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ghReleases []githubRelease
	err = json.NewDecoder(resp.Body).Decode(&ghReleases)
	if err != nil {
		return nil, err
	}

	for _, release := range ghReleases {
		semverRelease, err := semver.Parse(release.TagName)
		if err != nil {
			log.Println("Unable to parse tagName according to semver: ", err.Error())
			continue
		}
		for _, asset := range release.Assets {
			releases = append(releases, Release{
				PreRelease:  release.Prerelease,
				TagName:     release.TagName,
				Filename:    asset.Name,
				DownloadUrl: asset.BrowserDownloadUrl,
				Size:        asset.Size,
				semverTag:   semverRelease,
			})
		}
	}

	return
}

// Gets the latest release available that matches the given regex.
// This assumes that your tags follows semver
func GetLatestReleaseForPlatform(username, repository, fileRegex string, acceptPreleases bool) (release Release, err error) {
	fileRx, err := regexp.Compile(fileRegex)
	if err != nil {
		return release, err
	}

	releases, err := GetReleases(username, repository)
	if err != nil {
		return release, err
	}

	sort.Sort(sort.Reverse(semverReleases(releases)))

	for _, r := range releases {
		fmt.Println(r.TagName, r.Filename)
		// Filename should match the regex
		if fileRx.Match([]byte(r.Filename)) {
			if r.PreRelease && acceptPreleases {
				return r, err
			}
		}
	}

	return release, errors.New("No matching releases found")
}

// Checks if the provided release is newer than the running program version
// Useful for checking if you should update
func IsNewer(release Release, programVersion string) (bool, error) {
	pv, err := semver.Parse(programVersion)
	if err != nil {
		return false, err
	}

	return pv.LT(release.semverTag), nil
}
