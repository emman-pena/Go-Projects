package aggregator

import (
	"fmt"
	"sync"
)

type LogEntry struct {
	Source  string
	Message string
}

type Aggregator struct {
	mu    sync.Mutex
	logs  []LogEntry
	input chan LogEntry
}

// NewAggregator creates a new instance of Log Aggregator
func NewAggregator() *Aggregator {
	return &Aggregator{
		logs:  make([]LogEntry, 0),
		input: make(chan LogEntry, 100),
	}
}

// Start begins listening for incoming logs
func (a *Aggregator) Start() {
	go func() {
		for log := range a.input {
			a.mu.Lock()
			a.logs = append(a.logs, log)
			fmt.Printf("Received log from %s: %s\n", log.Source, log.Message)
			a.mu.Unlock()
		}
	}()
}

// Submit adds a new log entry
func (a *Aggregator) Submit(log LogEntry) {
	a.input <- log
}

// GetLogs retrieves all logs
func (a *Aggregator) GetLogs() []LogEntry {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]LogEntry(nil), a.logs...)
}
