package cmd

import (
	"github.com/LastCallMedia/vagabond/util"
	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
	"github.com/kardianos/osext"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
)

// The url to fetch releases from
const UpdateUrl = "https://api.github.com/repos/LastCallMedia/vagabond/releases"

// Replaces the running binary with a new version fetched from
// github
var CmdSelfUpdate = cli.Command{
	Name:   "self-update",
	Usage:  "Check for a new version of this tool",
	Action: runSelfUpdate,
}

func runSelfUpdate(ctx *cli.Context) {
	version := "latest"
	if len(ctx.Args()) > 0 {
		version = ctx.Args()[0]
	}
	release, err := getRelease(version)
	if err != nil {
		util.Fatal("Unable to fetch release data")
	}
	filename := "vagabond_" + runtime.GOOS + "_" + runtime.GOARCH
	asset, found := assetSearch(release.Assets, filename)
	if !found {
		util.Fatal("Unable to find a release asset for this OS and architecture")
	}
	srcfile, err := fetchAsset(asset)
	if err != nil {
		util.Fatalf("Failed fetching file: %s", err)
	}
	dstFile, err := osext.Executable()
	if err != nil {
		util.Fatalf("Failed determining current binary: %s", err)
	}
	err = copyFileOver(srcfile, dstFile)
	if err != nil {
		util.Fatalf("Failed replacing current binary.")
	}
	util.Successf("Updated to %s (%s)", *release.TagName, dstFile)
}

func getRelease(version string) (release *github.RepositoryRelease, err error) {
	client := github.NewClient(nil)
	if version == "latest" {
		release, _, err = client.Repositories.GetLatestRelease("LastCallMedia", "vagabond")
	} else {
		release, _, err = client.Repositories.GetReleaseByTag("LastCallMedia", "vagabond", version)
	}

	return
}

func assetSearch(assets []github.ReleaseAsset, filename string) (asset github.ReleaseAsset, found bool) {
	found = false
	for _, potentialAsset := range assets {
		if *potentialAsset.Name == filename {
			asset = potentialAsset
			found = true
			return
		}
	}
	return
}

func fetchAsset(asset github.ReleaseAsset) (filename string, err error) {
	file, err := ioutil.TempFile(os.TempDir(), "vagabond-release")
	filename = file.Name()
	if err != nil {
		return
	}
	res, err := http.Get(*asset.BrowserDownloadURL)
	if err != nil {
		return
	}
	defer res.Body.Close()
	_, err = io.Copy(file, res.Body)

	return
}

func copyFileOver(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
