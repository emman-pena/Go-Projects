/**
go mod init cli-task-manager

go get -u github.com/spf13/cobra

go build -o task-manager

-----------------

Add a task:
./task-manager add "Buy groceries"

List tasks:
./task-manager list

Mark a task as done:
./task-manager done 1

Delete a task:
./task-manager delete 1

IMPORTANT
.\task-manager add "Buy groceries"
Rename-Item task-manager task-manager.exe

.\task-manager.exe add "Buy groceries"
*/

package main

/**
encoding/json: For encoding and decoding the tasks to and from JSON format.
fmt: For formatted I/O operations like printing to the console.
os: For basic operating system operations (like file reading/writing).
strconv: For converting string inputs to integer IDs.
*/
import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

/*
*
Defines a Task struct with three fields:

ID: An integer representing the task ID.
Description: A string for the task description.
Completed: A boolean indicating whether the task is completed or not.
JSON tags (json:"id", json:"description", json:"completed") are used to
specify how each field should be stored when encoding or decoding from JSON
format.
*/
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// Specifies the filename (tasks.json) where the tasks are stored.
const taskFile = "tasks.json"

// Load tasks from file
/**
Purpose: This function loads tasks from the tasks.json file.
Steps:
It tries to read the file tasks.json.
If the file doesn't exist (os.IsNotExist), it returns an empty list of tasks.
If there's an error (e.g., file permissions), it returns an error.
If successful, it unmarshals (decodes) the JSON data into a slice of
Task structs ([]Task).
*/
func loadTasks() ([]Task, error) {
	file, err := os.ReadFile(taskFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil // No tasks exist yet, return empty list
		}
		return nil, fmt.Errorf("failed to read tasks file: %w", err)
	}
	var tasks []Task
	err = json.Unmarshal(file, &tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tasks file: %w", err)
	}
	return tasks, nil
}

// Save tasks to file
func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}
	err = os.WriteFile(taskFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save tasks file: %w", err)
	}
	return nil
}

// Add a new task
/**
Loads existing tasks.
Calculates a new ID (incrementing the existing number of tasks).
Appends a new task to the list.
Saves the updated list of tasks back to tasks.json.
*/
func addTask(description string) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	id := len(tasks) + 1
	tasks = append(tasks, Task{ID: id, Description: description, Completed: false})
	err = saveTasks(tasks)
	if err != nil {
		return err
	}
	fmt.Println("Task added successfully!")
	return nil
}

// List all tasks
func listTasks() error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return nil
	}
	fmt.Println("Tasks:")
	for _, task := range tasks {
		status := "Pending"
		if task.Completed {
			status = "Done"
		}
		fmt.Printf("[%d] %s - %s\n", task.ID, task.Description, status)
	}
	return nil
}

// Mark a task as done
func markTaskDone(id int) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	found := false
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Completed = true
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("task with ID %d not found", id)
	}
	err = saveTasks(tasks)
	if err != nil {
		return err
	}
	fmt.Printf("Task %d marked as done.\n", id)
	return nil
}

// Delete a task
func deleteTask(id int) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	found := false
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("task with ID %d not found", id)
	}
	err = saveTasks(tasks)
	if err != nil {
		return err
	}
	fmt.Printf("Task %d deleted successfully.\n", id)
	return nil
}

// Main function
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cli-task-manager [add|list|done|delete] [args]")
		return
	}

	command := os.Args[1]
	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: cli-task-manager add <task description>")
			return
		}
		description := os.Args[2]
		if err := addTask(description); err != nil {
			fmt.Println("Error:", err)
		}
	case "list":
		if err := listTasks(); err != nil {
			fmt.Println("Error:", err)
		}
	case "done":
		if len(os.Args) < 3 {
			fmt.Println("Usage: cli-task-manager done <task ID>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil || id <= 0 {
			fmt.Println("Invalid task ID. Please enter a positive number.")
			return
		}
		if err := markTaskDone(id); err != nil {
			fmt.Println("Error:", err)
		}
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: cli-task-manager delete <task ID>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil || id <= 0 {
			fmt.Println("Invalid task ID. Please enter a positive number.")
			return
		}
		if err := deleteTask(id); err != nil {
			fmt.Println("Error:", err)
		}
	default:
		fmt.Println("Unknown command:", command)
		fmt.Println("Usage: cli-task-manager [add|list|done|delete] [args]")
	}
}
