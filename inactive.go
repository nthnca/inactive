/*
inactive is a command line tool that can be used to determine if a computer has not been
active recently. For details see the README or https://github.com/nthnca/inactive.
*/
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// We should allow user configuration of these values.
// Possibly through command line or settings file.
const (
	minUptime            = time.Minute * 5
	minActivityThreshold = time.Minute * 25
	tmpFilePath          = "/tmp"
	tmpFilePrefix        = "stayawake."
)

func main() {
	uptime := getUptime()
	if uptime < minUptime { // time.Minute*5 {
		fmt.Println("System active. Uptime:", uptime)
		os.Exit(1)
	}

	lastActivity := getLastActivityTime()
	if lastActivity < minActivityThreshold {
		fmt.Println("System active. No activity seen for:", lastActivity)
		os.Exit(1)
	}

	fmt.Printf("System INactive. Uptime: %s, last activity seen %s ago.\n",
		uptime, lastActivity)
}

func getUptime() time.Duration {
	uptime, err := os.Open("/proc/uptime")
	if err != nil {
		fmt.Println("Can't open /proc/uptime: ", err)
		os.Exit(1)
	}

	b := make([]byte, 50)
	n, err := uptime.Read(b)
	if err != nil {
		fmt.Println("Can't read /proc/uptime: ", err)
		os.Exit(1)
	}

	tmp := string(b[:n])
	split := strings.Index(tmp, ".")
	chars := tmp[:split]
	seconds, err := strconv.Atoi(chars)
	if err != nil {
		fmt.Println("/proc/uptime can't be converted to an integer: ", tmp)
		os.Exit(1)
	}

	return time.Duration(seconds) * time.Second
}

func getLastActivityTime() time.Duration {
	f, err := os.Open(tmpFilePath)
	if err != nil {
		fmt.Println("Can't open: ", tmpFilePath, err)
		os.Exit(1)
	}

	fi, err := f.Readdir(2000)
	if err != nil {
		fmt.Println("Can't readdir: ", tmpFilePath, err)
		os.Exit(1)
	}

	var lastActive time.Time

	for _, element := range fi {
		if len(element.Name()) < len(tmpFilePrefix) {
			continue
		}

		if element.Name()[:len(tmpFilePrefix)] != tmpFilePrefix {
			continue
		}

		if element.ModTime().After(lastActive) {
			lastActive = element.ModTime()
		}
	}

	return time.Now().Sub(lastActive)
}
