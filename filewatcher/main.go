package main

/**
mkdir filewatcher && cd filewatcher
go mod init filewatcher

go get -u github.com/fsnotify/fsnotify

mkdir watched_directory

*/

/**
fmt: For printing output.
log: For logging errors.
fsnotify: The file system watcher library.
*/

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
)

/**
fsnotify.NewWatcher(): Creates a new watcher instance.
defer watcher.Close(): Ensures the watcher is closed when the program exits.
*/

func main() {
	// Initialize the watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Error creating watcher: %v", err)
	}
	defer watcher.Close()

	// Directory to monitor
	directory := "./watched_directory"

	// watcher.Add(directory): Adds the specified directory to the watch list.
	err = watcher.Add(directory)
	if err != nil {
		log.Fatalf("Error adding directory: %v", err)
	}
	fmt.Printf("Watching directory: %s\n", directory)

	// Create a channel to receive events
	done := make(chan bool)

	/**
	watcher.Events: A channel receiving file system events.
	watcher.Errors: A channel receiving errors.
	The select block listens for both event and error channels.

	event.Op: The type of operation (e.g., create, delete, write).
	event.Name: The name of the file affected by the operation.
	*/
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fmt.Printf("EVENT: %s\n", event)
				if event.Op&fsnotify.Create == fsnotify.Create {
					fmt.Printf("File created: %s\n", event.Name)
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					fmt.Printf("File deleted: %s\n", event.Name)
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Printf("File modified: %s\n", event.Name)
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					fmt.Printf("File renamed: %s\n", event.Name)
				}
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					fmt.Printf("File permissions changed: %s\n", event.Name)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Printf("ERROR: %v\n", err)
			}
		}
	}()

	// Wait forever
	<-done
}
