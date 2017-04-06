package gatherer

import (
    l "github.com/cloudurable/simplelog/logging"
	c "github.com/cloudurable/metricsd/common"
)

type Gatherers struct {
    logger l.Logger
    gatherers []c.MetricsGatherer
}

func NewGatherers(config *c.Config) *Gatherers {

    logger := c.GetLogger(config.Debug, "gatherers")
    gatherers := []c.MetricsGatherer{}

	for _, gathererName := range config.Gatherers {
		switch gathererName {
		case c.GATHERER_CPU:
			cpu := NewCPUMetricsGatherer(nil, config)
			if cpu != nil {
                gatherers = append(gatherers, cpu)
			}

		case c.GATHERER_DISK:
			disk := NewDiskMetricsGatherer(nil, config)
			if disk != nil {
				gatherers = append(gatherers, disk)
			}

		case c.GATHERER_FREE:
			free := NewFreeMetricGatherer(nil, config)
			if free != nil {
				gatherers = append(gatherers, free)
			}

		case c.GATHERER_CASSANDRA:
			tools := NewCassandraMetricGatherers(nil, config)
			if tools != nil {
				for _, tool := range tools {
					gatherers = append(gatherers, tool)
				}
			}
		}
	}

	return &Gatherers{logger, gatherers}
}

func (this *Gatherers) Gather() []c.Metric {

    metrics := []c.Metric{}

    for _, gatherer := range this.gatherers {
        more, err := gatherer.Gather()
        if err != nil {
            this.logger.PrintError("Problem getting metrics from gatherer: " + gatherer.Name(), err)
        } else if more != nil {
            metrics = append(metrics, more...)
        }
    }

    return metrics
}
