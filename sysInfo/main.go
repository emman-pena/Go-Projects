/**
mkdir sysinfo
cd sysinfo
go mod init sysinfo

Install dependencies
go get -u github.com/shirou/gopsutil/...

go run main.go

optional
create executable
go build -o sysinfo

run executable
./sysinfo

*/

package main

/**fmt: Provides formatting for input/output (e.g., printing to the console).
os: Allows interaction with the operating system, such as fetching environment
variables.

runtime: Contains utilities for OS and architecture detection.

time: Used for converting uptime into a human-readable format.

gopsutil: A third-party library to retrieve detailed system information like
memory, CPU, and host details.
*/

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func main() {

	/**
	runtime.GOOS: Returns the operating system (windows, linux, darwin, etc.).

	runtime.GOARCH: Returns the system architecture (amd64, arm, etc.).

	fmt.Printf: Formats and prints strings with placeholders (%s for strings).
	*/

	// OS and architecture
	fmt.Println("System Information")
	fmt.Println("===================")
	fmt.Printf("Operating System: %s\n", runtime.GOOS)
	fmt.Printf("Architecture: %s\n", runtime.GOARCH)

	/**
	host.Info(): Fetches details about the host machine, including hostname,
	uptime, and boot time.

	hostInfo.Hostname: Machine name.

	hostInfo.Uptime: Total uptime (in seconds); converted into a human-readable format
	using time.Duration.

	hostInfo.BootTime: Time the system last booted, formatted using time.Unix and
	time.RFC1123.

	err: Checks for errors. If there's an issue fetching data, it prints an error message.
	*/
	// Host information
	hostInfo, err := host.Info()
	if err == nil {
		fmt.Printf("Hostname: %s\n", hostInfo.Hostname)
		fmt.Printf("Uptime: %s\n", time.Duration(hostInfo.Uptime)*time.Second)
		fmt.Printf("Boot Time: %s\n", time.Unix(int64(hostInfo.BootTime), 0).Format(time.RFC1123))
	} else {
		fmt.Println("Error fetching host info:", err)
	}

	/**
	cpu.Info(): Fetches details about the CPU.

	cpuInfo[0].ModelName: Name of the CPU (e.g., Intel(R) Core(TM) i7-8700).

	cpuInfo[0].Cores: Number of cores in the CPU.
	*/

	// CPU information
	cpuInfo, err := cpu.Info()
	if err == nil && len(cpuInfo) > 0 {
		fmt.Printf("CPU: %s\n", cpuInfo[0].ModelName)
		fmt.Printf("Cores: %d\n", cpuInfo[0].Cores)
	} else {
		fmt.Println("Error fetching CPU info:", err)
	}

	/**
	mem.VirtualMemory(): Fetches virtual memory details.

	memInfo.Total: Total memory (in bytes), converted to gigabytes (/1e9).

	memInfo.Available: Available memory (in bytes), converted similarly.
	*/

	// Memory information
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		fmt.Printf("Total Memory: %.2f GB\n", float64(memInfo.Total)/1e9)
		fmt.Printf("Available Memory: %.2f GB\n", float64(memInfo.Available)/1e9)
	} else {
		fmt.Println("Error fetching memory info:", err)
	}

	/**
	os.Environ(): Returns a slice of all environment variables in KEY=VALUE format.

	The for loop iterates over the slice and prints each variable
	*/

	// Environment variables
	fmt.Println("\nEnvironment Variables:")
	for _, env := range os.Environ() {
		fmt.Println(env)
	}
}
