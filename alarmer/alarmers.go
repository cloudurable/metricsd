package alarmer

import (
	c "github.com/cloudurable/metricsd/common"
    l "github.com/cloudurable/simplelog/logging"
)

func LoadAlarmers(config *c.Config) ([]c.MetricsAlarmer) {

	var alarmers = []c.MetricsAlarmer{}

	for _, alarmerName := range config.Alarmers {
		switch alarmerName {
		case c.ALARMER_EMAIL:
			alarmer := NewEmailMetricAlarmer(config)
			if alarmer != nil {
				alarmers = append(alarmers, alarmer)
			}
		}
	}

	return alarmers
}

func VerifyAlarmer(alarmerName string, logger l.Logger, config *c.Config) {
    switch alarmerName {
    case c.ALARMER_EMAIL: verify(NewEmailMetricAlarmer(config), alarmerName, logger)
    }
}

func verify(alarmer c.MetricsAlarmer, alarmerName string, logger l.Logger) {
    if alarmer == nil {
        logger.Printf("%s alarmer not constructed!", alarmerName)
    } else if !alarmer.Verify() {
        logger.Printf("%s alarmer not verified!", alarmerName)
    } else {
        logger.Printf("%s alarmer verified!", alarmerName)
    }
}
