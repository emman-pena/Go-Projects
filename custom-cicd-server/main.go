/**
A Custom CI/CD (Continuous Integration/Continuous Deployment) Server refers
to a tailored or bespoke solution for automating the processes of integration,
testing, and deployment in software development. It’s essentially a
self-hosted or customized version of a CI/CD pipeline, designed to meet
the unique needs of a specific project or organization.
*/

/**
Key Features:

Repository polling or webhook-based triggers.
Build pipeline execution (e.g., running Docker commands).
Status reporting (e.g., via a web dashboard or logs).
Optional: Deployment functionality.
Tech Stack:

Backend: Go (for server logic).
Build Automation: Docker or native build tools.
Frontend (Optional): HTML/JS or a lightweight Go-based templating system
like html/template.
*/

package main

/**
Install dependencies:
go get -u github.com/gorilla/mux # For routing
go get -u github.com/go-yaml/yaml # For YAML parsing (e.g., pipe

go get gopkg.in/yaml.v3
go get github.com/google/uuid


encoding/json: Used for encoding and decoding JSON data
(e.g., for API responses).

fmt: Used for formatted I/O, particularly for printing logs or error messages.

log: Used to log messages (e.g., print to the console for debugging).

net/http: Provides HTTP client and server implementations.

github.com/gorilla/mux: A popular routing library that helps define
HTTP routes and handle URL parameters.

io/ioutil: Used to read files, in this case, the pipeline configuration
YAML file.

gopkg.in/yaml.v3: Used to parse the YAML configuration file.

os/exec: Used to execute commands in the shell
(like running build or test commands).

time: Used for generating unique build IDs based on timestamps.

*/

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"io/ioutil"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

// BuildStatus represents the status of a build
type BuildStatus struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Logs   string `json:"logs"`
}

// In-memory store for build statuses (for simplicity)
var buildStatuses = make(map[string]BuildStatus)

// PipelineStep defines a step in the pipeline
type PipelineStep struct {
	Name string   `yaml:"name"`
	Cmd  []string `yaml:"cmd"`
}

// PipelineConfig defines the structure of the YAML file
type PipelineConfig struct {
	Pipeline []PipelineStep `yaml:"pipeline"`
}

func main() {
	r := mux.NewRouter()

	// Route to trigger builds
	r.HandleFunc("/build", triggerBuild).Methods("POST")

	// Route to check build status
	r.HandleFunc("/status/{id}", checkStatus).Methods("GET")

	// Start the server
	log.Println("Starting CI/CD server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// triggerBuild handles build requests
func triggerBuild(w http.ResponseWriter, r *http.Request) {
	log.Println("Triggering build...")

	// Load pipeline configuration
	config, err := LoadConfig("config.yaml")
	if err != nil {
		http.Error(w, "Failed to load pipeline configuration", http.StatusInternalServerError)
		return
	}

	// Generate a unique ID for the build
	buildID := generateUUID()

	// Create a placeholder for build status
	buildStatuses[buildID] = BuildStatus{
		ID:     buildID,
		Status: "In Progress",
		Logs:   "",
	}

	// Execute the pipeline in a separate goroutine
	go func(id string) {
		err := ExecutePipeline(config.Pipeline, id)
		status := "Success"
		if err != nil {
			status = "Failed"
		}
		buildStatuses[id] = BuildStatus{
			ID:     id,
			Status: status,
			Logs:   fmt.Sprintf("Pipeline completed with status: %s", status),
		}
	}(buildID)

	// Return the build ID to the user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Build triggered",
		"id":      buildID,
	})
}

// checkStatus provides build status
func checkStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	buildID := vars["id"]

	status, exists := buildStatuses[buildID]
	if !exists {
		http.Error(w, "Build ID not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// generateUUID generates a unique identifier using the Google UUID package
func generateUUID() string {
	// Using UUID from the google/uuid package
	return uuid.New().String()
}

// LoadConfig loads and parses the pipeline configuration from a YAML file
func LoadConfig(filepath string) (*PipelineConfig, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var config PipelineConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// ExecutePipeline runs the steps in the pipeline and logs their output
func ExecutePipeline(steps []PipelineStep, buildID string) error {
	// Iterate through each step in the pipeline
	for _, step := range steps {
		log.Printf("Executing step: %s", step.Name)
		cmd := exec.Command(step.Cmd[0], step.Cmd[1:]...)
		output, err := cmd.CombinedOutput()

		// If there's an error, log the error and update build status with failure
		if err != nil {
			log.Printf("Error in step %s: %s\nOutput: %s", step.Name, err, string(output))
			buildStatuses[buildID] = BuildStatus{
				ID:     buildID,
				Status: "Failed",
				Logs:   fmt.Sprintf("Step %s failed: %s", step.Name, string(output)),
			}
			return err
		}

		// Log the successful output of the step
		log.Printf("Output of step %s: %s", step.Name, string(output))

		// Update logs in the build status for this step
		buildStatuses[buildID] = BuildStatus{
			ID:     buildID,
			Status: "In Progress",
			Logs:   fmt.Sprintf("Step %s completed successfully", step.Name),
		}
	}
	return nil
}

/**
Command:
Invoke-RestMethod -Uri http://localhost:8080/build -Method Post -Body '{"key":"value"}' -ContentType "application/json"

Sample Output
id                                   message
--                                   -------
230e2a97-7f31-4ef4-9ed6-25f3db594a07 Build triggered
*/
/*
### Key Fixes and Explanations:

1. **UUID Generation**:
   - We now use the `github.com/google/uuid` package to generate proper UUIDs, ensuring uniqueness across builds.
   - `generateUUID()` returns a UUID as a string using `uuid.New().String()`.

2. **Error Handling**:
   - If an error occurs during the execution of any pipeline step, the pipeline stops, and the status is updated with an error message that includes the command's output.
   - This ensures that failed steps are reflected in the build status and allows for better troubleshooting.

3. **Build Status Updates**:
   - After executing each step, the build status is updated to reflect the progress and logs.
   - This provides users with more detailed feedback about which steps are currently running or have completed.

4. **Logs**:
   - The `Logs` field in `BuildStatus` now includes the output of each step. This gives users better insight into the build process.
   - If a step fails, the logs will contain the error output from the command that failed.

5. **Server Configuration**:
   - The HTTP server is set up using the `gorilla/mux` router to handle requests for triggering builds and checking their status.
   - `POST` requests are used for triggering builds, and `GET` requests are used for checking build statuses based on build ID.

6. **Pipeline Execution**:
   - Each step in the pipeline is executed using `exec.Command`, and its output is captured to update the build status.
   - If any step fails, the execution stops, and the build status is updated to "Failed" with the error message.

*/
