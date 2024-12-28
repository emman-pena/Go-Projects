package main

/**
fmt: This is used for formatted I/O, like printing output to the console (fmt.Printf).

log: Provides logging functions. Here it is used to log errors.

os: Provides a platform-independent interface for interacting with the operating system.
It's included but not used directly in the script. You could use it to handle environment
variables or input/output redirection.

golang.org/x/crypto/ssh: This is the Go SSH package used to establish SSH connections and
execute commands remotely.

strings: Useful for string manipulation, but it's not used directly in this script.
time: Used to set a timeout for SSH connections to avoid indefinite hanging on unreachable
servers.
*/
import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/ssh" //go get -u golang.org/x/crypto/ssh
)

// Server represents a server to connect to
/**
Defines a struct Server to hold details of a server that you need to SSH into.
The fields are:
Host: The IP address or hostname of the server.
Port: The SSH port (typically "22").
Username: The username for SSH login.
Password: The password for SSH authentication.
*/
type Server struct {
	Host     string
	Port     string
	Username string
	Password string
}

// SSH connection function
/**
Purpose: Establish an SSH connection to a given server.
Steps:
ssh.ClientConfig: Creates an SSH client configuration with:
User: The username for authentication.
Auth: Uses the password for authentication (ssh.Password(server.Password)).
HostKeyCallback: ssh.InsecureIgnoreHostKey() ignores SSH host key verification
(for simplicity, but not secure for production).
Timeout: The timeout duration for the SSH connection (10 seconds in this case).
ssh.Dial: Tries to establish an SSH connection to the server by specifying the
host and port.
If an error occurs while dialing the SSH server, it is returned. Otherwise,
the client object is returned.
*/
func sshConnect(server Server) (*ssh.Client, error) {
	// Create the SSH configuration
	config := &ssh.ClientConfig{
		User: server.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(server.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	// Connect to the server
	client, err := ssh.Dial("tcp", server.Host+":"+server.Port, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %v", server.Host, err)
	}
	return client, nil
}

// Execute a command on a server
/**
Purpose: Executes a given command on the remote server using the
established SSH connection.
Steps:
client.NewSession(): Creates a new SSH session for running commands.
session.CombinedOutput(cmd): Runs the provided command cmd on the server.
It returns both stdout and stderr as a combined string.
If the command fails, the error is returned. If successful, the output
(as a string) is returned.
*/
func executeCommand(client *ssh.Client, cmd string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	// Execute the command
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %v", err)
	}
	return string(output), nil
}

// Automate tasks across multiple servers
/**
Purpose: Automates the task of connecting to multiple servers and executing a
command on each.
Steps:
For Loop: Iterates through each server in the servers slice.
Connect to server: It calls sshConnect(server) to establish an SSH connection.
If the connection fails, it logs the error and continues with the next server.
Execute command: If the connection is successful, it calls
executeCommand(client, cmd) to execute the specified command on the server.
If the command execution fails, it logs the error and continues with the
next server.
If successful, it prints the command output from the server.
*/
func automateTasks(servers []Server, cmd string) {
	for _, server := range servers {
		fmt.Printf("Connecting to server: %s\n", server.Host)

		// Connect to server
		client, err := sshConnect(server)
		if err != nil {
			log.Printf("Error connecting to server %s: %v\n", server.Host, err)
			continue
		}
		defer client.Close()

		// Execute command on server
		output, err := executeCommand(client, cmd)
		if err != nil {
			log.Printf("Error executing command on server %s: %v\n", server.Host, err)
			continue
		}

		// Print the output
		fmt.Printf("Output from server %s:\n%s\n", server.Host, output)
	}
}

/*
*
Purpose: The entry point of the program, where you define the servers
and the command to run.
Steps:
Define Servers: The servers slice contains two servers, each with their IP
address, SSH port, username, and password. You can add more servers to the list.
Command: The command to be executed on each server is "uptime", which shows
how long the server has been running.
Call automateTasks: The automateTasks function is called to execute the task
across all the servers in the list.
*/
func main() {
	// Define servers
	servers := []Server{
		{"192.168.1.1", "22", "user", "password"},
		{"192.168.1.2", "22", "user", "password"},
	}

	// Command to be executed
	command := "uptime"

	// Automate tasks
	automateTasks(servers, command)
}
