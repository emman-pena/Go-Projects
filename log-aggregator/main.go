package main

/**
fmt: For printing logs to the console.
net/http: To create and handle an HTTP server.
encoding/json: To handle JSON serialization and deserialization.
"log-aggregator/aggregator": Your custom package for aggregating logs.
*/
import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"log-aggregator/aggregator"
)

/*
*
sync.Once: Guarantees that a block of code runs only once,
even in a concurrent environment.
logInstance: Holds the single instance of the Aggregator.
*/
var (
	once        sync.Once
	logInstance *aggregator.Aggregator
)

/*
*
once.Do(func): Ensures the initialization block inside is executed only once.
aggregator.NewAggregator(): Creates a new instance of the log aggregator.
logInstance.Start(): Starts any internal operations for the aggregator
(e.g., background processing).
Returns the single logInstance.
*/
func getAggregatorInstance() *aggregator.Aggregator {
	once.Do(func() {
		logInstance = aggregator.NewAggregator()
		logInstance.Start()
	})
	return logInstance
}

/*
*
Validates HTTP Method: Ensures only POST requests are processed.
Otherwise, it returns a 405 Method Not Allowed response.
*/
func logHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	/**
	Decodes JSON Request Body: Parses the incoming JSON into a LogEntry struct.
	Error Handling: Returns a 400 Bad Request if the body is invalid.
	*/
	var logEntry aggregator.LogEntry
	err := json.NewDecoder(r.Body).Decode(&logEntry)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	/**
	Retrieves the singleton instance of the aggregator and submits
	the new log entry.
	*/
	getAggregatorInstance().Submit(logEntry)

	// Respond with a success message in JSON format.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"id":     "b1a1-b000dff4711e", // Replace with real logic for generating IDs if needed
	})
}

// Handles GET requests to /logs to retrieve all stored logs.
func getLogsHandler(w http.ResponseWriter, r *http.Request) {
	// Validates HTTP Method: Ensures only GET requests are processed.
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	// Retrieves all logs stored in the singleton aggregator.
	logs := getAggregatorInstance().GetLogs()

	/**
	Encodes Logs as JSON: Converts the logs into a JSON array and sends it
	as the response.
	Error Handling: Returns a 500 Internal Server Error if encoding fails.
	*/
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(logs)
	if err != nil {
		http.Error(w, "Failed to encode logs", http.StatusInternalServerError)
	}
}

/*
*
/log: Handled by logHandler, for adding logs.
/logs: Handled by getLogsHandler, for retrieving logs.
Starts the Server:

Listens on port 8080 and serves the registered routes.
Logs an error and terminates if the server fails to start.
*/
func main() {
	http.HandleFunc("/log", logHandler)
	http.HandleFunc("/logs", getLogsHandler)

	log.Println("Log Aggregator running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}

/**
Test the Application
Submit Logs
Using Invoke-RestMethod (Recommended)
Invoke-RestMethod -Uri "http://localhost:8080/log" `
                  -Method POST `
                  -ContentType "application/json" `
                  -Body '{"Source": "App1", "Message": "This is a test log"}'
or
Using Invoke-WebRequest
Invoke-WebRequest -Uri "http://localhost:8080/log" `
                  -Method POST `
                  -ContentType "application/json" `
                  -Body '{"Source": "App1", "Message": "This is a test log"}'

Retrieve Logs
Invoke-RestMethod -Uri "http://localhost:8080/logs" -Method GET
or
Invoke-WebRequest -Uri "http://localhost:8080/logs" -Method GET

*/
