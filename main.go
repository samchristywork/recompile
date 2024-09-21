package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	White = "\033[0m"
	Red   = "\033[31m"
	Grey  = "\033[90m"
	Reset = "\033[0m"
)

func runBuild() {
	cmd := exec.Command("go", "build", "./...")

	stderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(White), "Errors found:")
		fmt.Println(string(White), string(stderr), string(Reset))
	} else {
		fmt.Println(string(Grey), "Build successful, no errors found.", string(Reset))
	}
}

func watch(watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if strings.HasSuffix(event.Name, "~") {
				continue
			}

			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Remove == fsnotify.Remove {
				fmt.Printf(string(Grey)+"Detected change in file: %s"+string(Reset)+"\n", event.Name)
				runBuild()
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Error:", err)
		}
	}
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go watch(watcher)

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(Grey), "Watching for changes...", string(Reset))
	runBuild()

	<-done
}
