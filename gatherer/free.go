package gatherer

import (
	"fmt"
	l "github.com/advantageous/go-logback/logging"
	c "github.com/cloudurable/metricsd/common"
	"strings"
)

type FreeMetricGatherer struct {
	logger l.Logger
	debug  bool
	command string
}

func NewFreeMetricGatherer(logger l.Logger, config *c.Config) *FreeMetricGatherer {

	logger = c.EnsureLogger(logger, config.Debug, c.PROVIDER_FREE, c.FLAG_FREE)
	command := readFreeConfig(config, logger)

	return &FreeMetricGatherer{
		logger: logger,
		debug:  config.Debug,
		command: command,
	}
}

func readFreeConfig(config *c.Config, logger l.Logger) (string) {
	command := "/usr/bin/free"
	label := c.DEFAULT_LABEL

	if config.FreeCommand != c.EMPTY {
		command = config.FreeCommand
		label = c.CONFIG_LABEL
	}

	if config.Debug {
		logger.Println("Free gatherer initialized by:", label, "as:", command)
	}
	return command
}

func (gatherer *FreeMetricGatherer) GetMetrics() ([]c.Metric, error) {
	output, err := c.ExecCommand(gatherer.command)
	if err != nil {
		return nil, err
	}

	//	            total        used        free      shared  buff/cache   available
	//Mem:        8081400     6177200      404096      566316     1500104      959436
	//Swap:       8296444       23424     8273020

	var metrics = []c.Metric{}

	lines := strings.Split(output, c.NEWLINE)
	line1 := lines[1]
	line2 := lines[2]

	var total uint64
	var free uint64
	var used uint64
	var shared uint64
	var buffer uint64
	var available uint64
	var name string

	fmt.Sscanf(line1, "%s %d %d %d %d %d %d", &name, &total, &used, &free, &shared, &buffer, &available)

	if gatherer.debug {
		gatherer.logger.Printf("name %s  total %d  used %d  free %d  shared %d  buffer %d  available %d",
			                    name,    total,    used,    free,    shared,    buffer,    available)
	}

	metrics = append(metrics, c.Metric{c.MT_SIZE_K, c.MetricValue(free), c.EMPTY, "mFreeLvl", c.PROVIDER_RAM})
	metrics = append(metrics, c.Metric{c.MT_SIZE_K, c.MetricValue(used), c.EMPTY, "mUsedLvl", c.PROVIDER_RAM})
	metrics = append(metrics, c.Metric{c.MT_SIZE_K, c.MetricValue(shared), c.EMPTY, "mSharedLvl", c.PROVIDER_RAM})
	metrics = append(metrics, c.Metric{c.MT_SIZE_K, c.MetricValue(buffer), c.EMPTY, "mBufLvl", c.PROVIDER_RAM})
	metrics = append(metrics, c.Metric{c.MT_SIZE_K, c.MetricValue(available), c.EMPTY, "mAvailableLvl", c.PROVIDER_RAM})

	totalF := float64(total)

	freePercent := c.Percent(float64(free), totalF)
	metrics = append(metrics, c.Metric{c.MT_PERCENT, c.MetricValue(int64(freePercent)), c.EMPTY, "mFreePer", c.PROVIDER_RAM})

	usedPercent := c.Percent(float64(used), totalF)
	metrics = append(metrics, c.Metric{c.MT_PERCENT, c.MetricValue(int64(usedPercent)), c.EMPTY, "mUsedPer", c.PROVIDER_RAM})

	fmt.Sscanf(line2, "%s %d %d %d", &name, &total, &used, &free)

	if free == 0 && used == 0 && total == 0 {
		// do nothing
	} else {
		metrics = append(metrics, c.Metric{c.MT_SIZE_K, c.MetricValue(free), c.EMPTY, "mSwpFreeLvl", c.PROVIDER_RAM})
		metrics = append(metrics, c.Metric{c.MT_SIZE_K, c.MetricValue(used), c.EMPTY, "mSwpUsedLvl", c.PROVIDER_RAM})
		metrics = append(metrics, c.Metric{c.MT_SIZE_K, c.MetricValue(shared), c.EMPTY, "mSwpSharedLvl", c.PROVIDER_RAM})

		totalF = float64(total)
		freePercent = (float64(free) / totalF) * 100.0
		metrics = append(metrics, c.Metric{c.MT_PERCENT, c.MetricValue(int64(freePercent)), c.EMPTY, "mSwpFreePer", c.PROVIDER_RAM})
		usedPercent = (float64(used) / totalF) * 100.0
		metrics = append(metrics, c.Metric{c.MT_PERCENT, c.MetricValue(int64(usedPercent)), c.EMPTY, "mSwpUsedPer", c.PROVIDER_RAM})
	}

	return metrics, nil

}
