package repeater

import (
	c "github.com/cloudurable/metricsd/common"
    l "github.com/cloudurable/simplelog/logging"
)

func LoadRepeaters(config *c.Config) ([]c.MetricsRepeater) {

	var repeaters = []c.MetricsRepeater{}

	for _, repeaterName := range config.Repeaters {
		switch repeaterName {
		case c.REPEATER_AWS:
			repeater := NewAwsCloudMetricRepeater(config)
			if repeater != nil {
				repeaters = append(repeaters, repeater)
			}

		case c.REPEATER_LOGGER:
			repeater := NewLogMetricsRepeater()
			if repeater != nil {
				repeaters = append(repeaters, repeater)
			}

		case c.REPEATER_CONSOLE:
			repeater := NewConsoleMetricsRepeater()
			if repeater != nil {
				repeaters = append(repeaters, repeater)
			}
		}
	}

	return repeaters
}

func VerifyRepeater(repeaterName string, logger l.Logger, config *c.Config) {

    switch repeaterName {
    case c.REPEATER_AWS:     verify(NewAwsCloudMetricRepeater(config), repeaterName, logger)
    case c.REPEATER_LOGGER:  verify(NewLogMetricsRepeater(), repeaterName, logger)
    case c.REPEATER_CONSOLE: verify(NewConsoleMetricsRepeater(), repeaterName, logger)
    }
}

func verify(repeater c.MetricsRepeater, repeaterName string, logger l.Logger) {
    if repeater == nil {
        logger.Printf("%s repeater not constructed!", repeaterName)
    } else if !repeater.Verify() {
        logger.Printf("%s repeater not verified!", repeaterName)
    } else {
        logger.Printf("%s repeater verified!", repeaterName)
    }
}
