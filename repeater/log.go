package repeater

import (
	lg "github.com/cloudurable/simplelog/logging"
	c "github.com/cloudurable/metricsd/common"
)

type LogMetricsRepeater struct {
	logger lg.Logger
}

func (this LogMetricsRepeater) Repeat(metrics []c.Metric) error {
	for _, m := range metrics {
		this.logger.Println(c.ToJson(m))
	}
	return nil
}

func (this LogMetricsRepeater) Verify() bool { return true }

func (this LogMetricsRepeater) Name() string { return this.logger.Name() } // already there so just use it

func NewLogMetricsRepeater() *LogMetricsRepeater {
	return &LogMetricsRepeater{ lg.NewSimpleLogger(c.REPEATER + c.REPEATER_LOG)	}
}
