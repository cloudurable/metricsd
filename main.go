package main

import (
	l "github.com/cloudurable/simplelog/logging"
	c "github.com/cloudurable/metricsd/common"
	r "github.com/cloudurable/metricsd/repeater"
	g "github.com/cloudurable/metricsd/gatherer"
    a "github.com/cloudurable/metricsd/alarmer"
	"os"
	"os/signal"
	"syscall"
	"time"
	"flag"
)

func main() {

    version := flag.Bool("version", false, "Request Version")
    vaws := flag.Bool("vaws", false, "Verify AWS Credentials")
    vemail := flag.Bool("vemail", false, "Verify Email Alarmer")
    configFileName := flag.String("config", "/etc/metricsd.conf", "metrics config file name")

    flag.Parse()

    if *version { println("0.4.5") } else
    if *vaws    { doVerifyAws(configFileName) } else
    if *vemail  { doVerifyEmail(configFileName) } else
    { run(configFileName) }
}

func loadConfigForFirstTime(configFileName *string) *c.Config {
    logger := l.NewSimpleLogger("metricsd")
    config, err := c.LoadConfig(*configFileName, logger)
    if err != nil {
        logger.CriticalError("Error reading config", err)
        os.Exit(1)
    }
    return config
}

func run(configFileName *string) {
	// load the config file
	config := loadConfigForFirstTime(configFileName)
    logger := c.GetLogger(config.Debug, "main")

    logger.Debug("Init:", c.ToJsonLabeledString(config))
    interval, intervalConfigRefresh := readRunConfig(config)

    // services
    // mailer := s.NewMailer(logger, config)

	// terminator listener
    terminator := makeTerminateChannel()

    // workers
    gatherers := g.NewGatherers(config)
    repeaters := r.NewRepeaters(config)
    alarmers  := a.NewAlarmers(config)

    // timers
	timer := time.NewTimer(interval)
    configTimer := time.NewTimer(intervalConfigRefresh)

	var configChanged bool = false

	for {
		select {
		case <-terminator:
			logger.Info("Exiting")
			os.Exit(0)

		case <-timer.C:
			if configChanged {
				configChanged = false
                logger = c.GetLogger(config.Debug, "main")
				gatherers = g.NewGatherers(config)
                repeaters = r.NewRepeaters(config)
                alarmers = a.NewAlarmers(config)
			}
			metrics := gatherers.Gather()
            repeaters.Repeat(metrics)
            alarmers.Alarm(metrics)
			timer.Reset(interval)

		case <-configTimer.C:
			if newConfig, err := c.LoadConfig(*configFileName, logger); err != nil {
				logger.Error("Error reading config, changes ignored!", err)
			} else {
                changed := !c.ConfigEquals(config, newConfig)
				if changed {
                    configChanged = true
					config = newConfig
					interval, intervalConfigRefresh = readRunConfig(config)
                    logger.Debug("Changed:", c.ToJsonLabeledString(config))
				} else {
                    logger.Debug("Same Config")
				}
			}
			configTimer.Reset(intervalConfigRefresh)
		}
	}
}

func makeTerminateChannel() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	return ch
}

func readRunConfig(config *c.Config) (time.Duration, time.Duration){
	return config.TimePeriodSeconds * time.Second, config.ReadConfigSeconds * time.Second
}

func doVerifyAws(configFileName *string) {
    config := loadConfigForFirstTime(configFileName)
    logger := c.GetLogger(config.Debug, "verifyAws")
    r.VerifyRepeater(c.REPEATER_AWS, logger, config)
}

func doVerifyEmail(configFileName *string) {
    config := loadConfigForFirstTime(configFileName)
    logger := c.GetLogger(config.Debug, "verifyEmail")
    a.VerifyAlarmer(c.ALARMER_EMAIL, logger, config)
}
