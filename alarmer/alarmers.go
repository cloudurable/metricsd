package alarmer

import (
	c "github.com/cloudurable/metricsd/common"
    l "github.com/cloudurable/simplelog/logging"
)


type Alarmers struct {
    logger l.Logger
    alarmers []c.MetricsAlarmer
}

func NewAlarmers(config *c.Config) *Alarmers {
    logger := c.GetLogger(config.Debug, "alarmers")
	alarmers := []c.MetricsAlarmer{}

	for _, alarmerName := range config.Alarmers {
		switch alarmerName {
        case c.ALARMER_AWS:
            alarmer := NewAwsCloudMetricAlarmer(config)
            if alarmer != nil {
                alarmers = append(alarmers, alarmer)
            }

        case c.ALARMER_EMAIL:
            alarmer := NewEmailMetricAlarmer(config)
            if alarmer != nil {
                alarmers = append(alarmers, alarmer)
            }
		}
	}

	return &Alarmers{logger, alarmers}
}

func (this *Alarmers) Alarm(metrics []c.Metric) {
    for _, a := range this.alarmers {
        if err := a.Alarm(metrics); err != nil {
            this.logger.PrintError("Alarmer failed", err)
        }
    }
}

func VerifyAlarmer(alarmerName string, logger l.Logger, config *c.Config) {
    switch alarmerName {
    case c.ALARMER_AWS: verify(NewAwsCloudMetricAlarmer(config), alarmerName, logger)
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
