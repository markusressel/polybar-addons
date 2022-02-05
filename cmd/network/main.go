package main

import (
	"fmt"
	"github.com/markusressel/polybar-addons/internal/util"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	NetStatPath      = "/proc/net/dev"
	TmpStatsFilePath = "/dev/shm/network_traffic_last_state"

	// SectorSize see: https://stackoverflow.com/questions/37248948/how-to-get-disk-read-write-bytes-per-second-from-proc-in-programming-on-linux
	SectorSize = 512
)

// Outputs the current disk IO
//
//
// Examples:
// > disk
//    0 B/s   16.3 B/s
//
func main() {
	lastStats, lastTime, err := loadLastStats()
	_ = updateLastStats()
	if err != nil {
		fmt.Printf("Calculating...")
		return
	}

	currentStats, err := getCurrentStats()
	if err != nil {
		fmt.Printf("ERR")
		return
	}

	lastTotalReads, lastTotalWrites := aggregate(lastStats)
	currTotalReads, currTotalWrites := aggregate(currentStats)

	diff := time.Now().Sub(lastTime)

	readsSinceLast := currTotalReads - lastTotalReads
	writesSinceLast := currTotalWrites - lastTotalWrites

	formattedReads := util.FormatDataRate(readsSinceLast, diff)
	formattedWrites := util.FormatDataRate(writesSinceLast, diff)

	fmt.Printf("\uE2C6%s \uE2C4%s\n", formattedReads, formattedWrites)
}

func loadLastStats() ([]StatItem, time.Time, error) {
	finfo, _ := os.Stat(TmpStatsFilePath)
	stats, err := getStats(TmpStatsFilePath)
	return stats, finfo.ModTime(), err
}

func getCurrentStats() ([]StatItem, error) {
	return getStats(NetStatPath)
}

func getStats(path string) ([]StatItem, error) {
	text, err := readStatsFile(path)
	if err != nil {
		return nil, err
	}

	stats, err := parseStats(text)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func updateLastStats() error {
	text, err := readStatsFile(NetStatPath)
	if err != nil {
		return err
	}
	return util.WriteTextToFile(text, TmpStatsFilePath)
}

func aggregate(stats []StatItem) (int64, int64) {
	var totalReceived int64 = 0
	var totalTransmitted int64 = 0
	for _, stat := range stats {
		totalReceived += stat.receiveBytes
		totalTransmitted += stat.transmitBytes
	}
	return totalReceived, totalTransmitted
}

// StatItem https://www.kernel.org/doc/Documentation/iostats.txt
type StatItem struct {
	device string

	receiveBytes      int64
	receivePackets    int64
	receiveErrors     int64
	receiveDrops      int64
	receiveFifo       int64
	receiveFrame      int64
	receiveCompressed int64
	receiveMulticast  int64

	transmitBytes      int64
	transmitPackets    int64
	transmitErrors     int64
	transmitDrops      int64
	transmitFifo       int64
	transmitFrame      int64
	transmitCompressed int64
	transmitMulticast  int64
}

func parseStats(text string) ([]StatItem, error) {
	var items []StatItem

	lines := strings.Split(text, "\n")

	for i, line := range lines {
		if i < 2 {
			// skip irrelevant infos
			continue
		}
		parts := strings.Split(line, ":")
		name := parts[0]
		stats := parts[1]
		fields := strings.Fields(stats)
		statItem := createStatItem(name, fields)
		items = append(items, statItem)
	}

	return items, nil
}

func readStatsFile(path string) (string, error) {
	text, err := util.ReadTextFromFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

func createStatItem(name string, fields []string) StatItem {
	return StatItem{
		device:             name,
		receiveBytes:       toIntOr0(fields[0]),
		receivePackets:     toIntOr0(fields[1]),
		receiveErrors:      toIntOr0(fields[2]),
		receiveDrops:       toIntOr0(fields[3]),
		receiveFifo:        toIntOr0(fields[4]),
		receiveFrame:       toIntOr0(fields[5]),
		receiveCompressed:  toIntOr0(fields[6]),
		receiveMulticast:   toIntOr0(fields[7]),
		transmitBytes:      toIntOr0(fields[8]),
		transmitPackets:    toIntOr0(fields[9]),
		transmitErrors:     toIntOr0(fields[10]),
		transmitDrops:      toIntOr0(fields[11]),
		transmitFifo:       toIntOr0(fields[12]),
		transmitFrame:      toIntOr0(fields[13]),
		transmitCompressed: toIntOr0(fields[14]),
		transmitMulticast:  toIntOr0(fields[15]),
	}
}

func toIntOr0(s string) int64 {
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0
	} else {
		return int64(value)
	}
}
