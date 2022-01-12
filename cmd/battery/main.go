package main

import (
	"fmt"
	"github.com/markusressel/polybar-addons/internal/util"
	"math"
	"strings"
)

// Outputs the remaining time to fully charge/discharge a battery
//
// Examples:
// > battery
// 00:45
func main() {
	// TODO: optional -d --device parameter
	battery := "BAT0"

	// get value
	charging, err := isBatteryCharging(battery)
	if err != nil {
		fmt.Printf("ERR")
		return
	}
	energyTarget, err := getEnergyTarget(battery)
	if err != nil {
		fmt.Printf("ERR")
		return
	}
	energyNow, err := getEnergyNow(battery)
	if err != nil {
		fmt.Printf("ERR")
		return
	}
	powerNow, err := getPowerNow(battery)
	if err != nil {
		fmt.Printf("ERR")
		return
	}
	if err != nil {
		fmt.Printf("ERR")
		return
	}

	if powerNow == 0 {
		fmt.Printf("∞")
		return
	}

	var remainingTimeInSeconds int
	if charging == true {
		remainingTimeInSeconds = calculateRemainingTime(energyTarget-energyNow, powerNow)
	} else {
		remainingTimeInSeconds = calculateRemainingTime(energyNow, powerNow)
	}

	remainingHours := int(math.Min(99, float64(remainingTimeInSeconds/60/60)))
	remainingMinutes := (remainingTimeInSeconds / 60) % 60

	fmt.Printf("%02d:%02d", remainingHours, remainingMinutes)
}

func getEnergyTarget(battery string) (int, error) {
	chargeControlEndThreshold := getChargeControlEndThreshold(battery)
	energyFull, err := getEnergyFull(battery)
	return int((float64(energyFull) / 100) * float64(chargeControlEndThreshold)), err
}

func getChargeControlEndThreshold(battery string) int {
	path := "/sys/class/power_supply/" + battery + "/charge_control_end_threshold"
	value, err := util.ReadIntFromFile(path)
	if err != nil {
		value = 100
	}
	return value
}

func calculateRemainingTime(wh int, w int) int {
	return int((float64(wh) / float64(w)) * 60 * 60)
}

func isBatteryCharging(battery string) (bool, error) {
	path := "/sys/class/power_supply/" + battery + "/status"
	status, err := util.ReadTextFromFile(path)
	status = strings.TrimSpace(status)
	charging := status == "Charging"
	return charging, err
}

func getEnergyFull(battery string) (int, error) {
	path := "/sys/class/power_supply/" + battery + "/energy_full"
	return util.ReadIntFromFile(path)
}

func getEnergyNow(battery string) (int, error) {
	path := "/sys/class/power_supply/" + battery + "/energy_now"
	return util.ReadIntFromFile(path)
}

func getPowerNow(battery string) (int, error) {
	path := "/sys/class/power_supply/" + battery + "/power_now"
	return util.ReadIntFromFile(path)
}
