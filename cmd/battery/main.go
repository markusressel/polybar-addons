package main

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/markusressel/polybar-addons/util"
	"strings"
)

// Outputs the remaining time to fully charge/discharge a battery
//
// Examples:
//
//	>> python3 battery_time.py
//	1,8 hours
//
//	>> python3 battery_time.py -d "/org/freedesktop/UPower/devices/battery_BAT1"
//	1,8 hours
func main() {
	// TODO: optional -d --device parameter

	var devicePath string
	//if devicePath == nil {
	devices := getBatteryDevices()
	if len(devices) <= 0 {
		log.Fatalf("No battery")
	}
	devicePath = strings.TrimSpace(devices[0])
	//}

	// get value
	result, err := util.ExecCommand("upower", "-i", fmt.Sprintf("%s", devicePath))
	if err != nil {
		log.Error(err)
	}
	lines := strings.SplitAfter(result, "\n")

	filterFunc := func(s string) bool {
		return strings.Contains(strings.ToLower(s), "state") ||
			strings.Contains(strings.ToLower(s), "to\\ full") ||
			strings.Contains(strings.ToLower(s), "to\\ empty")
	}
	lines2 := filter(lines, filterFunc)
	state := findState(lines2)

	if state == "fully-charged" {
		fmt.Printf("")
	} else {
		timeRemaining := findTimeRemaining(lines)
		fmt.Printf("%s", timeRemaining)
	}
}

func getBatteryDevices() []string {
	result, err := util.ExecCommand("upower", "-e")
	if err != nil {
		return []string{}
	}

	lines := strings.SplitAfter(result, "\n")
	filterFunc := func(s string) bool { return strings.Contains(strings.ToLower(s), "bat") }
	return filter(lines, filterFunc)
}

func findState(result []string) string {
	filterFunc := func(s string) bool { return strings.Contains(strings.ToLower(s), "state") }
	lines := filter(result, filterFunc)
	if len(lines) <= 0 {
		return "Calculating..."
	}
	line := lines[0]
	state := strings.Split(line, ":")[1]
	state = strings.TrimSpace(state)
	return state
}

func findTimeRemaining(result []string) string {
	filterFunc := func(s string) bool { return strings.Contains(strings.ToLower(s), "time to") }
	lines := filter(result, filterFunc)
	if len(lines) <= 0 {
		return "Calculating..."
	}
	line := lines[0]
	state := strings.Split(line, ":")[1]
	state = strings.TrimSpace(state)
	return state
}

func filter(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}
