package sync

import (
	"github.com/koyachi/go-term-ansicolor/ansicolor"
	"gopkg.in/fsnotify.v1"
	"log"
)

type FSWatcher struct {
	*fsnotify.Watcher
	Files   chan string
	Folders chan string
	Changes <-chan string
}

// Watches the specific folder for changes
// When a change happens it triggers a notification on a channel
func WatchFolderRecursive(folder string) <-chan string {
	done := make(chan string)
	go func() {
		watcher, err := NewRecursiveWatcher(folder)
		if err != nil {
			log.Fatal(err)
		}
		watcher.Run(false)
		defer watcher.Close()

		for {
			select {
			case event := <-watcher.Events:
				done <- event.Name
				log.Println(ansicolor.IntenseBlack("FS Changed"), ansicolor.Underline(event.Name))
			case folder := <-watcher.Folders:
				log.Println(ansicolor.Yellow("Watching path"), ansicolor.Yellow(folder))
			}
		}
	}()

	return done
}
