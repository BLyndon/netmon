package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	monitorInterval      = 2 * time.Second  // Check traffic every 2 seconds
	thresholdMBps        = 1                // Alert if traffic exceeds this value in MBps
	changeThreshold      = 0.1              // Only log if traffic change is >= 0.1 MBps
	notificationCooldown = 10 * time.Second // Cooldown period after showing a notification
)

var lastNotificationTime time.Time

func GetNetworkStats() (int64, int64, error) {
	cmd := exec.Command("netstat", "-ib")
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	lines := strings.Split(string(output), "\n")
	var totalRecv, totalSent int64

	for _, line := range lines {
		fields := regexp.MustCompile(`\s+`).Split(strings.TrimSpace(line), -1)
		if len(fields) > 9 && fields[0] != "Name" { // Exclude headers
			recv, _ := strconv.ParseInt(fields[6], 10, 64) // Bytes received
			sent, _ := strconv.ParseInt(fields[9], 10, 64) // Bytes sent
			totalRecv += recv
			totalSent += sent
		}
	}

	return totalRecv, totalSent, nil
}

func ShowMacOSNotification(message string) {
	if time.Since(lastNotificationTime) > notificationCooldown {
		cmd := exec.Command("osascript", "-e", fmt.Sprintf(`display notification "%s" with title "🚨 Network Alert"`, message))
		_ = cmd.Run()
		lastNotificationTime = time.Now() // Update the last notification time
	}
}

func MonitorTraffic() {
	var prevRecv, prevSent, prevRxRate, prevTxRate int64

	for {
		recv, sent, err := GetNetworkStats()
		if err != nil {
			fmt.Println("Error fetching network stats:", err)
			time.Sleep(monitorInterval)
			continue
		}

		if prevRecv > 0 && prevSent > 0 {
			rxRate := (recv - prevRecv) / 1024 / 1024 / int64(monitorInterval.Seconds()) // MBps
			txRate := (sent - prevSent) / 1024 / 1024 / int64(monitorInterval.Seconds()) // MBps

			if absDiff(float64(rxRate), float64(prevRxRate)) >= changeThreshold ||
				absDiff(float64(txRate), float64(prevTxRate)) >= changeThreshold {

				fmt.Printf("📡 Updated Traffic - RX: %d MBps, TX: %d MBps\n", rxRate, txRate)

				if rxRate > thresholdMBps || txRate > thresholdMBps {
					alertMsg := fmt.Sprintf("High Traffic! RX: %d MBps, TX: %d MBps", rxRate, txRate)
					fmt.Println("🚨 " + alertMsg)
					ShowMacOSNotification(alertMsg)
				}
			}

			prevRxRate, prevTxRate = rxRate, txRate
		}

		prevRecv, prevSent = recv, sent
		time.Sleep(monitorInterval)
	}
}

func absDiff(a, b float64) float64 {
	if a > b {
		return a - b
	}
	return b - a
}

func main() {
	fmt.Println("📡 Network Traffic Monitor started...")
	MonitorTraffic()
}
