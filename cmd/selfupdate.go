package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"path/filepath"
)

const UpdateUrl = "https://api.github.com/repos/LastCallMedia/vagabond/releases"

var CmdSelfUpdate = cli.Command{
	Name:   "self-update",
	Usage:  "Check for a new version of this tool",
	Action: runSelfUpdate,
}

func runSelfUpdate(ctx *cli.Context) {
	release, err := getRelease("latest")
	if err != nil {
		log.Fatalf("Unable to fetch release data")
	}
	filename := "vagabond_" + runtime.GOOS + "_" + runtime.GOARCH
	asset, found := searchAssets(release.Assets, filename)
	if !found {
		log.Fatal("Unable to find a release asset for this OS and architecture")
	}
	if *asset.DownloadCount == 5 {
		log.Print("")
	}
	srcfile, err := fetchAsset(asset)
	if err != nil {
		log.Fatalf("Failed fetching file: %s", err)
	}
	dstFile, err := getCurrentBinary()
	if err != nil {
		log.Fatalf("Failed determining current binary: %s", err)
	}
	err = copyFileOver(srcfile, dstFile)
	if err != nil {
		log.Fatalf("Failed replacing current binary.")
	}
}

func getRelease(version string) (release *github.RepositoryRelease, err error) {
	client := github.NewClient(nil)
	release, _, err = client.Repositories.GetLatestRelease("LastCallMedia", "vagabond")

	return
}

func searchAssets(assets []github.ReleaseAsset, filename string) (asset github.ReleaseAsset, found bool) {
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

func getCurrentBinary() (filename string, err error) {
	filename, err = filepath.Abs(os.Args[0])
	if err != nil {
		return
	}
	_, err = os.Stat(filename)
	return
}
