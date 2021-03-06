/*
inactive is a command line tool that can be used to determine if a computer has not been
active recently. For details see the README or https://github.com/nthnca/inactive.
*/
package main

import (
	"fmt"
	"log"
	"log/syslog"
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

var (
	logWriter *syslog.Writer
)

func infof(str string, args ...interface{}) {
	err := logWriter.Info(str)
	if err != nil {
		log.Fatalf(str, args...)
	}
}

func fatalf(str string, args ...interface{}) {
	infof(fmt.Sprintf(str, args...))
	os.Exit(1)
}

func main() {
	var err error

	logWriter, err = syslog.New(syslog.LOG_INFO, "AmIStillActive")
	if err != nil {
		log.Fatalf("failed to open logger %s", err)
	}

	uptime := getUptime()
	if uptime < minUptime {
		fatalf("System active. Uptime: %s", uptime)
	}

	lastActivity := getLastActivityTime()
	if lastActivity < minActivityThreshold {
		fatalf("System active. No activity seen for: %s",
			lastActivity)
	}

	infof("System INactive. Uptime: %s, last activity seen %s ago.\n",
		uptime, lastActivity)
}

func getUptime() time.Duration {
	uptime, err := os.Open("/proc/uptime")
	if err != nil {
		fatalf("Can't open /proc/uptime: %s", err)
	}

	b := make([]byte, 50)
	n, err := uptime.Read(b)
	if err != nil {
		fatalf("Can't read /proc/uptime: %s", err)
	}

	tmp := string(b[:n])
	split := strings.Index(tmp, ".")
	chars := tmp[:split]
	seconds, err := strconv.Atoi(chars)
	if err != nil {
		fatalf("/proc/uptime can't be converted to an integer: %s",
			tmp)
	}

	return time.Duration(seconds) * time.Second
}

func getLastActivityTime() time.Duration {
	f, err := os.Open(tmpFilePath)
	if err != nil {
		fatalf("Can't open: %s %s", tmpFilePath, err)
	}

	fi, err := f.Readdir(2000)
	if err != nil {
		fatalf("Can't readdir: %s %s", tmpFilePath, err)
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
