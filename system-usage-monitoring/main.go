/**
go get github.com/shirou/gopsutil/v3

go get github.com/shirou/gopsutil/v3/cpu@v3.24.5
go get github.com/shirou/gopsutil/v3/disk@v3.24.5
go get golang.org/x/sys/windows@latest

*/

package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

func main() {
	for {
		// CPU Usage
		cpuPercent, err := cpu.Percent(0, false)
		if err != nil {
			fmt.Printf("Error fetching CPU usage: %v\n", err)
		} else {
			fmt.Printf("CPU Usage: %.2f%%\n", cpuPercent[0])
		}

		// Memory Usage
		vmStat, err := mem.VirtualMemory()
		if err != nil {
			fmt.Printf("Error fetching memory usage: %v\n", err)
		} else {
			fmt.Printf("Memory Usage: %.2f%% (%v/%v)\n", vmStat.UsedPercent, formatBytes(vmStat.Used), formatBytes(vmStat.Total))
		}

		// Disk Usage
		diskStat, err := disk.Usage("/")
		if err != nil {
			fmt.Printf("Error fetching disk usage: %v\n", err)
		} else {
			fmt.Printf("Disk Usage: %.2f%% (%v/%v)\n", diskStat.UsedPercent, formatBytes(diskStat.Used), formatBytes(diskStat.Total))
		}

		// Wait before the next iteration
		time.Sleep(1 * time.Second)
	}
}

// Helper function to format bytes into human-readable format
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

/**
Explanation
CPU Usage:

cpu.Percent(interval, percpu) calculates CPU usage as a percentage.
Pass 0 as the interval for instantaneous readings.
Memory Usage:

mem.VirtualMemory() retrieves the total, used, and free memory.
Use UsedPercent for the percentage of memory being used.
Disk Usage:

disk.Usage(path) retrieves disk usage statistics for the specified path (e.g., / for root).
Formatting:

The formatBytes function converts raw byte values into human-readable formats like KB, MB, GB, etc.
Looping:

The for loop continuously updates the metrics every second. You can adjust the sleep duration as needed.
*/
