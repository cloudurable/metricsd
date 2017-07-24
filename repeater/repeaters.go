package repeater

import (
	c "github.com/cloudurable/metricsd/common"
    l "github.com/cloudurable/simplelog/logging"
)

type Repeaters struct {
    logger l.Logger
    repeaters []c.MetricsRepeater
}

func NewRepeaters(config *c.Config) *Repeaters {

    logger := c.GetLogger(config.Debug, "repeaters")
	repeaters := []c.MetricsRepeater{}

	for _, repeaterName := range config.Repeaters {
		switch repeaterName {
		case c.REPEATER_AWS:
			repeater := NewAwsCloudMetricRepeater(config)
			if repeater != nil {
				repeaters = append(repeaters, repeater)
			}

        case c.REPEATER_LOG:
            repeater := NewLogMetricsRepeater()
            if repeater != nil {
                repeaters = append(repeaters, repeater)
            }

        case c.REPEATER_KAFKA:
            repeater := NewKafkaMetricsRepeater(config)
            if repeater != nil {
                repeaters = append(repeaters, repeater)
            }
		}
	}

	return &Repeaters{logger, repeaters}
}

func (this *Repeaters) Repeat(metrics []c.Metric) {
    for _, r := range this.repeaters {
        if err := r.Repeat(metrics); err != nil {
            this.logger.PrintError("Repeater failed", err)
        }
    }
}

func VerifyRepeater(repeaterName string, logger l.Logger, config *c.Config) {
    switch repeaterName {
    case c.REPEATER_AWS:   verify(NewAwsCloudMetricRepeater(config), logger)
    case c.REPEATER_LOG:   verify(NewLogMetricsRepeater(), logger)
    case c.REPEATER_KAFKA: verify(NewKafkaMetricsRepeater(config), logger)
    }
}

func verify(repeater c.MetricsRepeater, logger l.Logger) {
    if repeater == nil {
        logger.Printf("%s repeater not constructed!", repeater.Name())
    } else if !repeater.Verify() {
        logger.Printf("%s repeater not verified!", repeater.Name())
    } else {
        logger.Printf("%s repeater verified!", repeater.Name())
    }
}
