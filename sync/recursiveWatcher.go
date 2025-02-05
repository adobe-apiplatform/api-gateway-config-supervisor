package sync

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/koyachi/go-term-ansicolor/ansicolor"
)

type RecursiveWatcher struct {
	fsw     *fsnotify.Watcher
	Files   chan string
	Folders chan string
}

func DebugMessage(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println(ansicolor.IntenseBlack(msg))
}

func DebugError(msg error) {
	fmt.Println(ansicolor.IntenseBlack(msg.Error()))
}

func NewRecursiveWatcher(path string) (*RecursiveWatcher, error) {
	folders := Subfolders(path)
	if len(folders) == 0 {
		return nil, errors.New("No folders to watch.")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	rw := &RecursiveWatcher{fsw: watcher}

	rw.Files = make(chan string, 10)
	rw.Folders = make(chan string, len(folders))

	for _, folder := range folders {
		rw.AddFolder(folder)
	}
	return rw, nil
}

func (watcher *RecursiveWatcher) AddFolder(folder string) {
	err := watcher.fsw.Add(folder)
	if err != nil {
		log.Println("Error watching: ", folder, err)
	}
	watcher.Folders <- folder
}

func (watcher *RecursiveWatcher) Run(debug bool) {
	go func() {
		for {
			select {
			case event := <-watcher.fsw.Events:
				if debug {
					log.Println("Event: ", event)
				}
				// create a file/directory
				if event.Op&fsnotify.Create == fsnotify.Create {
					fi, err := os.Stat(event.Name)
					if err != nil {
						// eg. stat .subl513.tmp : no such file or directory
						if debug {
							DebugError(err)
						}
					} else if fi.IsDir() {
						if debug {
							DebugMessage("Detected new directory %s", event.Name)
						}
						for _, folder := range Subfolders(event.Name) {
							watcher.AddFolder(folder)
						}
					} else {
						if debug {
							DebugMessage("Detected new file %s", event.Name)
						}
						watcher.Files <- event.Name // created a file
					}
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					// modified a file, assuming that you don't modify folders
					if debug {
						DebugMessage("Detected file modification %s", event.Name)
					}
					watcher.Files <- event.Name
				}

				if event.Op&fsnotify.Remove == fsnotify.Remove {
					if debug {
						DebugMessage("Detected remove event %s", event.Name)
					}
					watcher.Files <- event.Name
				}

			case err := <-watcher.fsw.Errors:
				log.Println("error", err)
			}
		}
	}()
}

// Subfolders returns a slice of subfolders (recursive), including the folder provided.
func Subfolders(path string) (paths []string) {
	filepath.Walk(path, func(newPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			name := info.Name()
			if shouldIgnoreFile(name) {
				return filepath.SkipDir
			}
			paths = append(paths, newPath)
		}
		return nil
	})
	return paths
}

// shouldIgnoreFile determines if a file should be ignored.
// File names that begin with "." or ".." or "_" are ignored by the go tool.
func shouldIgnoreFile(name string) bool {
	return strings.HasPrefix(name, ".") || strings.HasPrefix(name, "..") || strings.HasPrefix(name, "_")
}
