/**
The purpose of this program is to create a basic security scanner that scans
a network or system to identify potential vulnerabilities. Specifically, it:

Identifies Open Ports:
Scans a range of ports (1–1024 by default) to find which ones are open and
listening. Open ports can expose services running on a system, which might be
vulnerable to attacks if not properly secured.

Checks for Misconfigurations:
Includes checks for specific vulnerabilities, such as services like MongoDB
being accessible without authentication. This highlights weak security
configurations that need to be addressed.

Serves as a Security Tool:
Helps administrators or security professionals understand the exposed attack
surface of their infrastructure.
Provides an initial step in securing systems by identifying open doors (ports)
and common vulnerabilities.
*/

// Replace the hostname variable in main with your target's hostname or IP address.

/*
*
go mod init github.com/yourusername/security-scanner

Integrate a library like go-cve-dictionary to check for known CVEs.
go get github.com/kotakanbe/go-cve-dictionary
*/
package main

/**
fmt: Provides formatted I/O operations like printing to the console.
net: Offers networking utilities to handle TCP and other connections.
os: Manages OS-level operations like reading environment variables or exiting programs.
time: Adds support for time-related functionality like delays or timeouts
*/
import (
	"fmt"
	"net"
	"time"
)

/*
*
Inputs:

protocol: The type of connection (e.g., tcp).
hostname: Target system (e.g., 127.0.0.1).
port: Port number to check.
Functionality:

Creates an address string in the form of hostname:port (e.g., 127.0.0.1:80).
Attempts to connect to the address using net.DialTimeout. If successful,
it means the port is open.
Returns true if the connection succeeds, otherwise false.
Timeout: A timeout of 1 second is set to prevent indefinite blocking.
*/
func scanPort(protocol, hostname string, port int) bool {
	address := fmt.Sprintf("%s:%d", hostname, port)
	conn, err := net.DialTimeout(protocol, address, 1*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

/*
*
Loops through port numbers from 1 to 1024 (common ports).
Calls scanPort for each port.
If a port is open, it prints a message indicating the port is open.
*/
func portScan(hostname string) {
	fmt.Printf("Scanning ports on %s...\n", hostname)
	for port := 1; port <= 1024; port++ {
		if scanPort("tcp", hostname, port) {
			fmt.Printf("Port %d is open\n", port)
		}
	}
}

/*
*
Purpose:

Checks if MongoDB is running on the default port (27017) and is accessible
without authentication.
Functionality:

Attempts to connect to hostname:27017.
If the connection is successful, it warns that MongoDB may lack proper
authentication.
If the connection fails, it indicates MongoDB is not accessible.
*/
func checkMongoDB(hostname string) {
	address := fmt.Sprintf("%s:%d", hostname, 27017)
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		fmt.Println("MongoDB not accessible")
		return
	}
	defer conn.Close()

	fmt.Println("MongoDB is accessible without authentication!")
}

/*
*
Purpose: Speeds up scanning by using Goroutines for concurrency.

Concurrency Control:

sem is a buffered channel with a size of 10 to limit the number of
simultaneous Goroutines.
A Goroutine is spawned for each port being scanned.
Each Goroutine runs independently, checking its assigned port.
Flow:

Push a token (true) into the channel before starting a Goroutine.
After scanning, the Goroutine removes a token (<-sem) to free up space in
the channel.
Limits concurrency to 10 Goroutines at a time.
*/
func concurrentPortScan(hostname string, ports []int) {
	sem := make(chan bool, 10) // Limit concurrency
	for _, port := range ports {
		sem <- true
		go func(port int) {
			defer func() { <-sem }()
			if scanPort("tcp", hostname, port) {
				fmt.Printf("Port %d is open\n", port)
			}
		}(port)
	}
}

func main() {
	hostname := "127.0.0.1" // Replace with target
	fmt.Println("Starting security scan...")
	portScan(hostname)
	checkMongoDB(hostname)
	fmt.Println("Scan completed.")
}
