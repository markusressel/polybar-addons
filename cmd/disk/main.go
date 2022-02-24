package main

import (
	"fmt"
	"github.com/markusressel/polybar-addons/internal/util"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	DiskStatsPath    = "/proc/diskstats"
	TmpStatsFilePath = "/dev/shm/disk_traffic_last_stats"

	// SectorSize see: https://stackoverflow.com/questions/37248948/how-to-get-disk-read-write-bytes-per-second-from-proc-in-programming-on-linux
	SectorSize = 512
)

// Outputs the current disk IO
//
//
// Examples:
// > disk "%reads% %writes%"
//    0.0 B/s    7.6MB/s
//
func main() {
	template := os.Args[1]

	lastStats, lastTime, err := loadLastStats()
	_ = updateLastStats()

	currentStats, err := getCurrentStats()
	if err != nil {
		fmt.Printf("ERR")
		return
	}

	var lastTotalReads, lastTotalWrites int64
	if lastStats != nil {
		lastTotalReads, lastTotalWrites = aggregate(lastStats)
	} else {
		lastTotalReads, lastTotalWrites = 0, 0
	}

	currTotalReads, currTotalWrites := aggregate(currentStats)

	diff := time.Now().Sub(lastTime)

	readsSinceLast := currTotalReads - lastTotalReads
	writesSinceLast := currTotalWrites - lastTotalWrites

	formattedReads := util.FormatDataRate(readsSinceLast, diff)
	formattedWrites := util.FormatDataRate(writesSinceLast, diff)

	placeholders := map[string]string{}
	placeholders["reads"] = formattedReads
	placeholders["writes"] = formattedWrites

	result := util.ReplacePlaceholders(template, placeholders)

	fmt.Print(result)
}

func loadLastStats() ([]StatItem, time.Time, error) {
	finfo, err := os.Stat(TmpStatsFilePath)
	if err != nil {
		return nil, time.Time{}, err
	}
	stats, err := getStats(TmpStatsFilePath)
	if err != nil {
		return nil, time.Time{}, err
	}
	return stats, finfo.ModTime(), err
}

func getCurrentStats() ([]StatItem, error) {
	return getStats(DiskStatsPath)
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

	return getMainStats(stats), nil
}

func updateLastStats() error {
	text, err := readStatsFile(DiskStatsPath)
	if err != nil {
		return err
	}
	return util.WriteTextToFile(text, TmpStatsFilePath)
}

func aggregate(stats []StatItem) (int64, int64) {
	var totalReads int64 = 0
	var totalWrites int64 = 0
	for _, stat := range stats {
		totalReads += stat.sectorsRead * SectorSize
		totalWrites += stat.sectorsWritten * SectorSize
	}
	return totalReads, totalWrites
}

func getMainStats(stats []StatItem) []StatItem {
	pattern := regexp.MustCompile("^nvme\\dn\\d$|^sd[a-z]$")

	var result []StatItem
	for _, stat := range stats {
		if stat.device.minorNumber == 0 || pattern.MatchString(stat.device.name) {
			result = append(result, stat)
		}
	}
	return result
}

type Device struct {
	majorNumber int64
	minorNumber int64
	name        string
}

// StatItem https://www.kernel.org/doc/Documentation/iostats.txt
type StatItem struct {
	device Device

	// Field  1 -- # of reads completed
	// This is the total number of reads completed successfully.
	reads int64

	// Field  2 -- # of reads merged, field 6 -- # of writes merged
	// Reads and writes which are adjacent to each other may be merged for
	// efficiency.  Thus two 4K reads may become one 8K read before it is
	// ultimately handed to the disk, and so it will be counted (and queued)
	// as only one I/O.This field lets you know how often this was done.
	readsMerged int64

	// Field  3 -- # of sectors read
	// This is the total number of sectors read successfully.
	sectorsRead int64

	// Field  4 -- # of milliseconds spent reading
	// This is the total number of milliseconds spent by all reads (as
	// measured from __make_request() to end_that_request_last()).
	timeMillisReads int64

	// Field  5 -- # of writes completed
	// This is the total number of writes completed successfully.
	writes int64

	// Field  6 -- # of writes merged
	// See the description of field 2.
	writesMerged int64

	// Field  7 -- # of sectors written
	// This is the total number of sectors written successfully.
	sectorsWritten int64

	// Field  8 -- # of milliseconds spent writing
	// This is the total number of milliseconds spent by all writes (as
	// measured from __make_request() to end_that_request_last()).
	timeMillisWriting int64

	// Field  9 -- # of I/Os currently in progress
	// The only field that should go to zero.Incremented as requests are
	// given to appropriate struct request_queue and decremented as they finish.
	currentIops int64

	// Field 10 -- # of milliseconds spent doing I/Os
	// This field increases so long as field 9 is nonzero.
	timeMillisIops int64

	// Field 11 -- weighted # of milliseconds spent doing I/Os
	// This field is incremented at each I/O start, I/O completion, I/O
	// merge, or read of these stats by the number of I/Os in progress (field 9) times the number of milliseconds spent doing I/O since the
	// last update of this field.This can provide an easy measure of both
	// I/O completion time and the backlog that may be accumulating.
	weightedTimeMillisIops int64

	// Field 12 -- # of discards completed
	// This is the total number of discards completed successfully.
	discards int64

	// Field 13 -- # of discards merged
	// See the description of field 2
	discardsMerged int64

	// Field 14 -- # of sectors discarded
	// This is the total number of sectors discarded successfully.
	sectorsDiscarded int64

	// Field 15 -- # of milliseconds spent discarding
	// This is the total number of milliseconds spent by all discards (as
	// measured from __make_request() to end_that_request_last()).
	timeMillisDiscarding int64
}

func parseStats(stats string) ([]StatItem, error) {
	var items []StatItem

	lines := strings.Split(stats, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		statItem := createStatItem(fields)
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

func createStatItem(fields []string) StatItem {
	return StatItem{
		device: Device{
			majorNumber: toIntOr0(fields[0]),
			minorNumber: toIntOr0(fields[1]),
			name:        fields[2],
		},
		reads:                  toIntOr0(fields[3]),
		readsMerged:            toIntOr0(fields[4]),
		sectorsRead:            toIntOr0(fields[5]),
		timeMillisReads:        toIntOr0(fields[6]),
		writes:                 toIntOr0(fields[7]),
		writesMerged:           toIntOr0(fields[8]),
		sectorsWritten:         toIntOr0(fields[9]),
		timeMillisWriting:      toIntOr0(fields[10]),
		currentIops:            toIntOr0(fields[11]),
		timeMillisIops:         toIntOr0(fields[12]),
		weightedTimeMillisIops: toIntOr0(fields[13]),
		discards:               toIntOr0(fields[14]),
		discardsMerged:         toIntOr0(fields[15]),
		sectorsDiscarded:       toIntOr0(fields[16]),
		timeMillisDiscarding:   toIntOr0(fields[17]),
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
