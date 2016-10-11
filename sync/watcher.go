package sync

import (
	"github.com/koyachi/go-term-ansicolor/ansicolor"
	"log"
)

// Watches the specific folder for changes
// When a change happens it triggers a notification on a channel
func WatchFolderRecursive(folder string, debug bool) <-chan string {
	done := make(chan string)
	go func() {
		watcher, err := NewRecursiveWatcher(folder)
		if err != nil {
			log.Fatal(err)
		}
		watcher.Run(debug)
		defer watcher.fsw.Close()

		for {
			select {
			case file := <-watcher.Files:
				done <- file
				log.Println(ansicolor.IntenseBlack("FS Changed"), ansicolor.Underline(file))
			case folder := <-watcher.Folders:
				log.Println(ansicolor.Yellow("Watching path"), ansicolor.Yellow(folder))
			}
		}
	}()

	return done
}
