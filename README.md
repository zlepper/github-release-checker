# GitHub Release Checker
Checks github for releases for repositories. 

## A release
A release for GitHub Release Checker looks like this:
```go
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
}
```

## Get all releases
To get all the releases for a repository:

```go
package main

import (
    ghc "github.com/zlepper/github-release-checker"
    "fmt"
)


func main() {
    releases, err := ghc.GetReleases("zlepper", "gfs")
    if err != nil {
        panic(err)
    }
    
    fmt.Println(releases)
}
```


## Get latest release for a given platform
Say we want to check if a release of a product is available, than we are already 
running, we can do the folowing:

```go
package main

import (
    ghc "github.com/zlepper/github-release-checker"
    "fmt"
)

const (
    productVersion string = "0.0.2"
    fileNamePattern string = "gfs-windows-x64"
    preRelease bool = true
)

func main() {
    release, err := ghc.GetLatestReleaseForPlatform("zlepper", "gfs", fileNamePattern, preRelease)
    if err != nil {
        panic(err)
    }
    
    newer, err := ghc.IsNewer(release, productVersion)
    if err != nil {
        panic(err)
    }
    
    if newer {
        fmt.Println("A new version is available")
    } else {
        fmt.Println("Running latest version already")
    }
}

```