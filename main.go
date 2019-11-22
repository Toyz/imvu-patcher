package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	c "github.com/otiai10/copy"
	"github.com/toyz/imvu-patcher/core"
)

var working_folder = ""
var versions_folder = ""
var logPath = ""

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	working_folder = path.Join(user.HomeDir, ".imvu-patcher")
	makeDir(working_folder)
	versions_folder = path.Join(working_folder, "versions")
	makeDir(versions_folder)
	logPath = path.Join(working_folder, "logs")
	makeDir(logPath)

	log.Printf("Working in directory: %s", working_folder)
	log.Printf("Running on: %s", runtime.GOOS)
	os.Chdir(working_folder)

	items := core.VersionList()
	latest := items[len(items)-1]

	versionFolder := path.Join(versions_folder, latest.Hash)
	logPath = versionFolder

	if !existsDir(versionFolder) {
		patch(latest, versionFolder)

		if runtime.GOOS == "linux" {
			log.Println("Patching paths for linux host")
			filepath.Walk(versionFolder, visit)
		}

		log.Println("Installing electron")
		cmd := exec.Command("npm", "install", "electron", "--save-dev")
		cmd.Dir = path.Join(versionFolder, "app")
		cmdReader, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
			os.Exit(1)
		}

		scanner := bufio.NewScanner(cmdReader)
		go func() {
			for scanner.Scan() {
				fmt.Printf("%s\n", scanner.Text())
			}
		}()

		cmd.Run()
	}

	ioutil.WriteFile(path.Join(versionFolder, "app", "package.json"), core.CreatePackageJson(path.Join(versionFolder, "app", "package.json")), 0755)

	log.Printf("Launching from: %s", path.Join(versionFolder, "app"))
	cmd := exec.Command("npm", "start")
	cmd.Dir = path.Join(versionFolder, "app")

	if err := cmd.Start(); err != nil {
		log.Panicln(err)
	}

	time.Sleep(5 * time.Second)
}

func patch(latest core.Versions, versionPath string) {
	zip := latest.Download()
	log.Printf("Downloaded: %s", zip.Filename)

	zip.Extract()

	c.Copy(path.Join(zip.Path, zip.FileNameNoExt, "IMVU Desktop.app", "Contents", "Resources", "app"), path.Join(versionPath, "app"))

	os.RemoveAll(path.Join(zip.Path, zip.FileNameNoExt))
	os.Remove(path.Join(zip.Path, zip.Filename))

	ioutil.WriteFile("latest", []byte(latest.String()), 0755)
}

func visit(path string, fi os.FileInfo, err error) error {

	if err != nil {
		return err
	}

	if !!fi.IsDir() {
		return nil //
	}

	matched, err := filepath.Match("*.js", fi.Name())

	if err != nil {
		panic(err)
		return err
	}

	if matched {
		read, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		newContents := ""
		newContents = strings.Replace(string(read), "/var', 'local", logPath, -1)
		newContents = strings.Replace(newContents, "/var', 'log", logPath, -1)

		err = ioutil.WriteFile(path, []byte(newContents), 0)
		if err != nil {
			panic(err)
		}

	}

	return nil
}

func makeDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}
}

func existsDir(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
