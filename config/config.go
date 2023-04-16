package config

import (
	"fmt"
	"os"
	"regexp"
	"time"
)

type configData struct {
	watchFolder          string
	includeFiles         *regexp.Regexp
	startServerScript    string
	killServerScript     string
	pidFile              string
	waitForServerKill    time.Duration
	timeoutBetweenChecks time.Duration
}

func readEnv(name, deafult string) string {
	val := os.Getenv(name)
	if val == "" {
		return deafult
	}
	return val
}

func readFromEnv() configData {
	defaults := map[string]string{
		"WATCH_FOLDER":           "/app/dev",
		"INCLUDE_FILES":          ".*",
		"START_SERVER_SCRIPT":    "/app/scripts/start.sh",
		"KILL_SERVER_SCRIPT":     "/app/scripts/kill.sh",
		"PID_FILE":               "/tmp/server.pid",
		"WAIT_FOR_SERVER_KILL":   "100ms",
		"TIMEOUT_BETWEEN_CHECKS": "500ms",
	}
	values := map[string]string{}
	for key, value := range defaults {
		values[key] = readEnv(key, value)
	}
	waitForServerKill, err := time.ParseDuration(values["WAIT_FOR_SERVER_KILL"])
	if err != nil {
		fmt.Println("Invalid value for WAIT_FOR_SERVER_KILL: " + values["WAIT_FOR_SERVER_KILL"])
		fmt.Println("Using default value: " + defaults["WAIT_FOR_SERVER_KILL"])
		waitForServerKill, _ = time.ParseDuration(defaults["WAIT_FOR_SERVER_KILL"])
	}
	timeoutBetweenChecks, err := time.ParseDuration(values["TIMEOUT_BETWEEN_CHECKS"])
	if err != nil {
		fmt.Println("Invalid value for TIMEOUT_BETWEEN_CHECKS: " + values["TIMEOUT_BETWEEN_CHECKS"])
		fmt.Println("Using default value: " + defaults["TIMEOUT_BETWEEN_CHECKS"])
		timeoutBetweenChecks, _ = time.ParseDuration(defaults["TIMEOUT_BETWEEN_CHECKS"])
	}
	return configData{
		watchFolder:          values["WATCH_FOLDER"],
		includeFiles:         regexp.MustCompile(values["INCLUDE_FILES"]),
		startServerScript:    values["START_SERVER_SCRIPT"],
		killServerScript:     values["KILL_SERVER_SCRIPT"],
		pidFile:              values["PID_FILE"],
		waitForServerKill:    waitForServerKill,
		timeoutBetweenChecks: timeoutBetweenChecks,
	}
}

var config = readFromEnv()

func WatchFolder() string {
	return config.watchFolder
}

func IncludeFiles() *regexp.Regexp {
	return config.includeFiles
}

func StartServerScript() string {
	return config.startServerScript
}

func KillServerScript() string {
	return config.killServerScript
}

func PidFile() string {
	return config.pidFile
}

func WaitForServerKill() time.Duration {
	return config.waitForServerKill
}

func TimeoutBetweenChecks() time.Duration {
	return config.timeoutBetweenChecks
}
