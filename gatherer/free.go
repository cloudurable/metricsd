package gatherer

import (
	"fmt"
	l "github.com/cloudurable/simplelog/logging"
	c "github.com/cloudurable/metricsd/common"
	"strings"
)

type FreeMetricGatherer struct {
	logger l.Logger
	command string
}

func NewFreeMetricGatherer(logger l.Logger, config *c.Config) *FreeMetricGatherer {

	logger = c.EnsureLogger(logger, config.Debug, c.GATHERER_FREE)
	command := c.ReadConfigString("free command", config.FreeConfig.Command, "/usr/bin/free", logger)

	return &FreeMetricGatherer{
		logger: logger,
		command: command,
	}
}

func (this *FreeMetricGatherer) Name() string {
    return c.GATHERER_FREE
}

func (this *FreeMetricGatherer) Gather() ([]c.Metric, error) {
	output, err := c.ExecCommand(this.command)
	if err != nil {
		return nil, err
	}

	//	            total        used        free      shared  buff/cache   available
	//Mem:        8081400     6177200      404096      566316     1500104      959436
	//Swap:       8296444       23424     8273020

	var metrics = []c.Metric{}

	lines := strings.Split(output, c.NEWLINE)

	var total int64
	var free int64
	var used int64
	var shared int64
	var buffer int64
	var available int64
	var name string

	// Mem:
	fmt.Sscanf(lines[1], "%s %d %d %d %d %d %d", &name, &total, &used, &free, &shared, &buffer, &available)
	this.logger.Debugf("%s, total %d, used %d, free %d, shared %d, buffer %d, available %d", name, total, used, free, shared, buffer, available)

	metrics = append(metrics, *c.NewMetricInt(c.MT_SIZE_KB, free, "frMemFree", c.EMPTY, c.PROVIDER_FREE))
	metrics = append(metrics, *c.NewMetricInt(c.MT_SIZE_KB, used, "frMemUsed", c.EMPTY, c.PROVIDER_FREE))
	metrics = append(metrics, *c.NewMetricInt(c.MT_SIZE_KB, shared, "frMemShared", c.EMPTY, c.PROVIDER_FREE))
	metrics = append(metrics, *c.NewMetricInt(c.MT_SIZE_KB, buffer, "frMemBuf", c.EMPTY, c.PROVIDER_FREE))
	metrics = append(metrics, *c.NewMetricInt(c.MT_SIZE_KB, available, "frMemAvailable", c.EMPTY, c.PROVIDER_FREE))

	totalF := float64(total)

	pct := c.PercentRoundDp(float64(free), totalF, 2)
	metrics = append(metrics, *c.NewMetricFloat(c.MT_PERCENT, pct, "frMemFreePct", c.EMPTY, c.PROVIDER_FREE))

    pct = c.PercentRoundDp(float64(used), totalF, 2)
	metrics = append(metrics, *c.NewMetricFloat(c.MT_PERCENT, pct, "frMemUsedPct", c.EMPTY, c.PROVIDER_FREE))

    // Swap:
	fmt.Sscanf(lines[2], "%s %d %d %d", &name, &total, &used, &free)
    this.logger.Debugf("%s, total %d, used %d, free %d", name, total, used, free)

	if total != 0 || free != 0 || used != 0 {
		metrics = append(metrics, *c.NewMetricInt(c.MT_SIZE_KB, free, "frSwapFreeLvl", c.EMPTY, c.PROVIDER_FREE))
		metrics = append(metrics, *c.NewMetricInt(c.MT_SIZE_KB, used, "frSwapUsedLvl", c.EMPTY, c.PROVIDER_FREE))
		metrics = append(metrics, *c.NewMetricInt(c.MT_SIZE_KB, shared, "frSwapSharedLvl", c.EMPTY, c.PROVIDER_FREE))

		totalF = float64(total)
        pct = c.PercentRoundDp(float64(free), totalF, 2)
		metrics = append(metrics, *c.NewMetricFloat(c.MT_PERCENT, pct, "frSwapFreePct", c.EMPTY, c.PROVIDER_FREE))
        pct = c.PercentRoundDp(float64(used), totalF, 2)
		metrics = append(metrics, *c.NewMetricFloat(c.MT_PERCENT, pct, "frSwapUsedPct", c.EMPTY, c.PROVIDER_FREE))
	}

	return metrics, nil
}
