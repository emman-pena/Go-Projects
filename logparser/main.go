package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// A LogEntry struct represents each parsed log line.
type LogEntry struct {
	Timestamp string
	Level     string
	Message   string
}

/**
nil is a special constant used to represent the zero value of various types,
such as pointers, interfaces, maps, slices, channels, and function types.
It's a way of indicating that a variable does not point to a valid memory
location or does not contain a value.
*/

func main() {
	// Open the log file
	file, err := os.Open("sample.log")
	/**
	err != nil checks if the err variable is not nil,
	meaning an error occurred during the operation.
	*/
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		/**
		return immediately exits the current function (often main or another function).
		This is used to stop the execution of the program when a critical error occurs.
		*/
		return
	}
	defer file.Close()

	// Parse the log file
	logEntries, err := parseLogFile(file)
	if err != nil {
		fmt.Printf("Error parsing log file: %v\n", err)
		return
	}

	// Analyze the logs
	analyzeLogs(logEntries)
}

// parseLogFile reads and parses the log file into structured log entries
func parseLogFile(file *os.File) ([]LogEntry, error) {

	var logEntries []LogEntry

	/**
	Reads the file line by line using bufio.Scanner.
	Uses a regular expression to extract the timestamp, log level, and message
	from each line.

	scanner.Scan() is a method of the bufio.Scanner type in Go,
	which is used to read input line by line (usually from a file or a string).
	*/
	scanner := bufio.NewScanner(file)

	// Regex to match log lines

	/**
	The caret (^) is an anchor that matches the beginning of a line.
	This ensures that the pattern matches from the start of the line,
	so the line won't start with any characters other than those defined in the regex.

	(\S+ \S+) is a capturing group that matches two sequences of non-whitespace
	characters (\S+) separated by a single space.

	(.+) is another capturing group that matches one or more characters of any kind.
	The dot (.) matches any character except newline, and the plus (+) means one or
	more occurrences of any character.

	The dollar sign ($) is an anchor that matches the end of a line.
	This ensures that the pattern will match until the end of the line
	*/
	logLineRegex := regexp.MustCompile(`^(\S+ \S+) (\S+) (.+)$`)

	for scanner.Scan() {
		line := scanner.Text() //scanner.Text() retrieves the current line that was just read by the scanner.
		matches := logLineRegex.FindStringSubmatch(line)
		/**
		FindStringSubmatch(line) is a method of the regexp package,
		which attempts to match the string line against the regular expression.
		If there is a match, FindStringSubmatch returns a slice of strings containing:
		The full match (the entire line).
		Submatches corresponding to each capture group in the regular expression
		(usually parts of the line you’re interested in). matches will be a slice where:
		matches[0] is the entire matched line.
		matches[1], matches[2], and matches[3] are the capture groups (specific parts
		of the log line you're interested in, e.g., timestamp, log level, and message).

		if len(matches) == 4
		This checks if the regular expression found exactly 4 parts in the matches slice.
		*/
		if len(matches) == 4 {
			// This creates a new LogEntry struct with the following fields
			logEntries = append(logEntries, LogEntry{
				Timestamp: matches[1],
				Level:     matches[2],
				Message:   matches[3],
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return logEntries, nil
}

// analyzeLogs performs basic analysis on the parsed logs
func analyzeLogs(logEntries []LogEntry) {
	// Count log levels
	/**
	levelCount is a map that stores the count of log entries for each log level
	(e.g., "INFO", "ERROR").

	logEntries is expected to be a slice of LogEntry structs,
	where each LogEntry contains fields like Timestamp, Level, and Message.
	range logEntries returns two values

	levelCount[entry.Level]++
	entry.Level refers to the Level field of the current LogEntry object.
	The Level field should contain a string value, such as "INFO", "ERROR", etc.
	*/
	levelCount := make(map[string]int)
	for _, entry := range logEntries {
		levelCount[entry.Level]++
	}

	// Print analysis
	fmt.Println("Log Level Summary:")
	for level, count := range levelCount {
		fmt.Printf("  %s: %d\n", strings.ToUpper(level), count)
	}

	// Find error messages
	fmt.Println("\nError Messages:")
	for _, entry := range logEntries {
		if entry.Level == "ERROR" {
			fmt.Printf("  [%s] %s\n", entry.Timestamp, entry.Message)
		}
	}
}
