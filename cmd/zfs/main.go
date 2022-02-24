package main

import (
	"fmt"
	"github.com/markusressel/polybar-addons/internal/util"
	"os"
	"strings"
)

// Outputs ZFS pool stats
//
// Examples:
// > zfs "%bpool.cap% (%bpool.free%) | %rpool.cap% (%rpool.free%)"
// 5% (3.54G) | 21% (725G)
//
func main() {
	template := os.Args[1]

	stats, err := readStats()
	if err != nil {
		fmt.Printf("ERR")
		return
	}

	placeholders := map[string]string{}
	for _, stat := range stats {
		placeholders[stat.name+"."+"free"] = stat.free
		placeholders[stat.name+"."+"used"] = stat.alloc
		placeholders[stat.name+"."+"cap"] = stat.cap
		placeholders[stat.name+"."+"total"] = stat.size
	}

	result := util.ReplacePlaceholders(template, placeholders)

	fmt.Print(result)
}

func readStats() ([]statItem, error) {
	result, err := util.ExecCommand("zpool", "list", "-o", "name,size,alloc,cap,free", "-H")
	if err != nil {
		return nil, err
	}

	return parseStats(result)
}

type statItem struct {
	name  string
	size  string
	alloc string
	cap   string
	free  string
}

func parseStats(stats string) ([]statItem, error) {
	var items []statItem

	lines := strings.Split(stats, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 5 {
			continue
		}
		statItem := createStatItem(fields)
		items = append(items, statItem)
	}

	return items, nil
}

func createStatItem(fields []string) statItem {
	return statItem{
		fields[0],
		fields[1],
		fields[2],
		fields[3],
		fields[4],
	}
}
