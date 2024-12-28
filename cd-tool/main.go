package main

/**
fmt: Used for printing messages to the console.
log: Provides logging functionality for errors or informational messages.
os: Enables interaction with the operating system (e.g., environment variables, directory paths).
os/exec: Allows running external commands.
path/filepath: Assists in working with file paths.
time: Used for time-related operations (e.g., logging deployment time).
github.com/fsnotify/fsnotify: A library for monitoring filesystem changes.
*/
import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Config holds configuration details for the deployment

/*
*
The Config struct stores:

RepoPath: Path to the directory being monitored (e.g., a Git repository).
BuildCmd: Command to build the application (e.g., go build).
DeployCmd: Command to deploy the application
(e.g., running the built executable).
*/
type Config struct {
	RepoPath  string
	BuildCmd  string
	DeployCmd string
}

func main() {
	// Step 1: Define the configuration
	config := Config{
		RepoPath:  "C:/Users/ethan/GoProjects/test-repo", // Replace with your repo path
		BuildCmd:  "echo Building application...",        // "go build -o app",    Build command
		DeployCmd: "echo Deploying application...",       // "./app", Deployment command
	}

	// Step 2: Start watching the repository
	fmt.Println("Starting Continuous Deployment Tool...")
	watchRepo(config) // Calls watchRepo to monitor the specified directory for changes.
}

// fsnotify.NewWatcher() sets up a system to monitor changes in the filesystem.
func watchRepo(config Config) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	/**
	Listens for events: Monitors for file write or create operations.
	Triggers deployment: If a file is modified or created, calls deploy.
	*/
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fmt.Println("Detected change:", event)

				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					fmt.Println("Change detected, deploying...")
					deploy(config)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error:", err)
			}
		}
	}()

	/**
	Adds the directory to the watcher: Watches the RepoPath for changes.
	Blocks execution: Keeps the tool running with <-done.
	*/
	err = watcher.Add(config.RepoPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Watching for changes in:", config.RepoPath)
	<-done
}

// Runs git pull to fetch the latest changes from the repository.
func deploy(config Config) {
	// Step 3: Pull the latest changes
	fmt.Println("Pulling latest changes...")
	if err := runCommand("git", "pull", "origin", "main"); err != nil {
		log.Println("Error pulling changes:", err)
		return
	}

	// Step 4: Build the application
	fmt.Println("Building application...")
	if err := runCommand("sh", "-c", config.BuildCmd); err != nil {
		log.Println("Error building application:", err)
		return
	}

	// Step 5: Deploy the application
	fmt.Println("Deploying application...")
	if err := runCommand("sh", "-c", config.DeployCmd); err != nil {
		log.Println("Error deploying application:", err)
		return
	}

	fmt.Println("Deployment completed successfully at", time.Now())
}

/*
*
Executes a command: Uses exec.Command to run the specified command (name)
with arguments (args).
Redirects output: Sends the command's stdout and stderr to the tool's console.
Sets working directory: Defaults to the current directory.
*/
func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = filepath.Dir(".")
	return cmd.Run()
}

/**
Create a local directory for testing purposes
mkdir C:\Users\ethan\GoProjects\test-repo
cd C:\Users\ethan\GoProjects\test-repo

echo "package main\n\nfunc main() {\n\tprintln(\"Hello, World!\")\n}" > main.go

*/
