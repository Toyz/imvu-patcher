package core

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/parnurzeal/gorequest"
)

type Versions struct {
	Hash string
	URL  string
	Size int
}

func VersionList() []Versions {
	_, body, _ := gorequest.New().Get("https://api.imvu.com/desktop_update/darwin/RELEASES").End()

	lines := strings.Split(body, "\n")

	var versions []Versions

	for _, line := range lines {
		items := strings.Split(line, " ")

		if items[0] == "" {
			continue
		}

		versions = append(versions, Versions{
			Hash: items[0],
			URL:  items[1],
			Size: toInt(items[2]),
		})
	}

	return versions
}

func (v Versions) Download() ZipZile {
	name := filepath.Base(v.URL)
	s, _ := filepath.Abs(".")

	if fileExists(path.Join(s, name)) {
		return ZipZile{
			Filename:      name,
			Path:          s,
			FileNameNoExt: noExtension(name),
		}
	}

	client := grab.NewClient()
	req, _ := grab.NewRequest(".", v.URL)

	// start download
	fmt.Printf("Downloading %v...\n", req.URL())
	resp := client.Do(req)
	fmt.Printf("  %v\n", resp.HTTPResponse.Status)

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			fmt.Printf("  transferred %v / %v bytes (%.2f%%)\n",
				resp.BytesComplete(),
				resp.Size,
				100*resp.Progress())

		case <-resp.Done:
			// download is complete
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		os.Exit(1)
	}

	return ZipZile{
		Filename:      name,
		Path:          s,
		FileNameNoExt: noExtension(name),
	}
}

func (v Versions) String() string {
	return fmt.Sprintf("%s %s %d", v.Hash, v.URL, v.Size)
}
