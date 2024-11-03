package main

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	Red   = "\033[31m"
	Grey  = "\033[90m"
	Reset = "\033[0m"
)

func runBuild(command string) {
	cmd := exec.Command("sh", "-c", command)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(Red + "Errors found:")
		fmt.Println(Red + string(output) + Reset)
	} else {
		fmt.Println(Grey + "Build successful, no errors found." + Reset)
		if len(output) > 0 {
			fmt.Print(string(output))
		}
	}
}

func watch(watcher *fsnotify.Watcher, command string, ignore []string) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if strings.HasSuffix(event.Name, "~") {
				continue
			}

			ignoreFlag := false
			for _, ig := range ignore {
				for _, part := range strings.Split(event.Name, string(filepath.Separator)) {
					if part == ig {
						ignoreFlag = true
						break
					}
				}
				if ignoreFlag {
					break
				}
			}
			if ignoreFlag {
				continue
			}

			if event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Remove == fsnotify.Remove {
				fmt.Printf(Grey+"Detected change in file: %s"+Reset+"\n", event.Name)
				runBuild(command)
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
	command := flag.String("command", "go build ./...", "Command to run")

	ignore := []string{".git", ".cache"}
	flag.Func("ignore", "File or directory to ignore", func(s string) error {
		fmt.Println("Ignoring file or directory: " + s)
		ignore = append(ignore, s)
		return nil
	})

	flag.Parse()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go watch(watcher, *command, ignore)

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			for _, ig := range ignore {
				if filepath.Base(path) == ig {
					return filepath.SkipDir
				}
			}
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(Grey + "Watching for changes..." + Reset)
	runBuild(*command)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
}
