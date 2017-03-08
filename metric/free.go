package metric

import (
	"fmt"
	l "github.com/advantageous/go-logback/logging"
	"os/exec"
	"runtime"
	"strings"
)

type FreeMetricGatherer struct {
	logger l.Logger
	debug  bool
	command string
}

func NewFreeMetricGatherer(logger l.Logger, config *Config) *FreeMetricGatherer {

	logger = ensureLogger(logger, config.Debug, "free", "MT_FREE_DEBUG")

	command := "/usr/bin/free"
	label := LINUX_LABEL

	if config.FreeCommand != EMPTY {
		command = config.FreeCommand
		label = CONFIG_LABEL
	} else if runtime.GOOS == GOOS_DARWIN {
		command = "/usr/local/bin/free"
		label = DARWIN_LABEL
	}

	if config.Debug {
		logger.Println("Free gatherer initialized by:", label, "as:", command)
	}

	return &FreeMetricGatherer{
		logger: logger,
		debug:  config.Debug,
		command: command,
	}
}

func (gatherer *FreeMetricGatherer) GetMetrics() ([]Metric, error) {
	var metrics = []Metric{}
	var output string

	if out, err := exec.Command(gatherer.command).Output(); err != nil {
		return nil, err
	} else {
		output = string(out)
	}

	lines := strings.Split(output, "\n")
	line1 := lines[1]
	line2 := lines[2]

	var total uint64
	var free uint64
	var used uint64
	var shared uint64
	var buffer uint64
	var available uint64
	var mem string

	fmt.Sscanf(line1, "%s %d %d %d %d %d %d", &mem, &total, &used, &free, &shared, &buffer, &available)

	if gatherer.debug {
		gatherer.logger.Printf("name %s total %d, used %d, free %d,"+
			" shared %d , buffer %d, available %d\n", mem, total, used, free, shared, buffer, available)
	}

	metrics = append(metrics, metric{LEVEL, MetricValue(free), "mFreeLvl", "ram", EMPTY})
	metrics = append(metrics, metric{LEVEL, MetricValue(used), "mUsedLvl", "ram", EMPTY})
	metrics = append(metrics, metric{LEVEL, MetricValue(shared), "mSharedLvl", "ram", EMPTY})
	metrics = append(metrics, metric{LEVEL, MetricValue(buffer), "mBufLvl", "ram", EMPTY})
	metrics = append(metrics, metric{LEVEL, MetricValue(available), "mAvailableLvl", "ram", EMPTY})

	totalF := float64(total)

	freePercent := (float64(free) / totalF) * 100.0
	metrics = append(metrics, metric{LEVEL_PERCENT, MetricValue(int64(freePercent)), "mFreePer", "ram", EMPTY})

	usedPercent := (float64(used) / totalF) * 100.0
	metrics = append(metrics, metric{LEVEL_PERCENT, MetricValue(int64(usedPercent)), "mUsedPer", "ram", EMPTY})

	fmt.Sscanf(line2, "%s %d %d %d", &mem, &total, &used, &free)

	if free == 0 && used == 0 && total == 0 {
		// do nothing
	} else {
		metrics = append(metrics, metric{LEVEL, MetricValue(free), "mSwpFreeLvl", "ram", EMPTY})
		metrics = append(metrics, metric{LEVEL, MetricValue(used), "mSwpUsedLvl", "ram", EMPTY})
		metrics = append(metrics, metric{LEVEL, MetricValue(shared), "mSwpSharedLvl", "ram", EMPTY})

		totalF = float64(total)
		freePercent = (float64(free) / totalF) * 100.0
		metrics = append(metrics, metric{LEVEL_PERCENT, MetricValue(int64(freePercent)), "mSwpFreePer", "ram", EMPTY})
		usedPercent = (float64(used) / totalF) * 100.0
		metrics = append(metrics, metric{LEVEL_PERCENT, MetricValue(int64(usedPercent)), "mSwpUsedPer", "ram", EMPTY})
	}

	return metrics, nil

}
