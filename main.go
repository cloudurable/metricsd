package main

import (
	l "github.com/cloudurable/simplelog/logging"
	c "github.com/cloudurable/metricsd/common"
	r "github.com/cloudurable/metricsd/repeater"
	g "github.com/cloudurable/metricsd/gatherer"
    s "github.com/cloudurable/metricsd/service"
	"os"
	"os/signal"
	"syscall"
	"time"
	"flag"
    "fmt"
)

func main() {

	version := flag.Bool("version", false, "Request Version")
	verifyAws := flag.Bool("verifyAws", false, "Verify AWS Credentials")
	configFileName := flag.String("config", "/etc/metricsd.conf", "metrics config file name")
    emailTo := flag.String("e", c.EMPTY, "Test Email, 'To' Address")

	flag.Parse()

	if *version            { println("0.4.5") } else
    if *verifyAws          { doVerifyAws(configFileName) } else
    if *emailTo != c.EMPTY { doEmailTo(configFileName, emailTo) } else
                           { run(configFileName) }
}

func prepare(configFileName *string, logName string) (l.Logger, *c.Config) {
    logger := l.NewSimpleLogger(logName)
    config, err := c.LoadConfig(*configFileName, logger)
    if err != nil {
        logger.CriticalError("Error reading config", err)
        os.Exit(1)
    }

    return logger, config
}

func run(configFileName *string) {
	// load the config file
	logger, config := prepare(configFileName, "main-run")

    // services
    // mailer := s.NewMailer(logger, config)

	logger = c.GetLogger(config.Debug, "main")
	logger.Debug("Init:", c.ObjectToString(config))

	// begin the work
	interval, intervalConfigRefresh, debug := readRunConfig(config)

	timer := time.NewTimer(interval)
	configTimer := time.NewTimer(intervalConfigRefresh)

	terminator := makeTerminateChannel()

	var gatherers []c.MetricsGatherer
	var repeaters []c.MetricsRepeater
	var configChanged bool = true

	for {
		select {
		case <-terminator:
			logger.Info("Exiting")
			os.Exit(0)

		case <-timer.C:
			if configChanged {
				configChanged = false
				gatherers = g.LoadGatherers(config)
				repeaters = r.LoadRepeaters(config)
			}
			metrics := collectMetrics(gatherers, logger)
			processMetrics(metrics, repeaters, config, logger)
			timer.Reset(interval)

		case <-configTimer.C:
			if newConfig, err := c.LoadConfig(*configFileName, logger); err != nil {
				logger.Error("Error reading config, changes ignored!", err)
			} else {
                changed := !c.ConfigEquals(config, newConfig)
				if changed {
					config = newConfig
					interval, intervalConfigRefresh, debug = readRunConfig(config)
					if debug {
						logger.Debug("Changed:", c.ObjectToString(config))
					}
					configChanged = true
				} else if debug {
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

func readRunConfig(config *c.Config) (time.Duration, time.Duration, bool){
	return config.TimePeriodSeconds * time.Second,
		config.ReadConfigSeconds * time.Second,
		config.Debug
}

func processMetrics(metrics []c.Metric, repeaters []c.MetricsRepeater, context *c.Config, logger l.Logger) {
	for _, r := range repeaters {
		if r.RepeatForContext() {
			if err := r.ProcessMetrics(context, metrics); err != nil {
				logger.PrintError("Repeater failed", err)
			}
		}
	}

	noIdContext := context.GetNoIdContext()

	for _, r := range repeaters {
		if r.RepeatForNoIdContext() {
			if err := r.ProcessMetrics(noIdContext, metrics); err != nil {
				logger.PrintError("Repeater failed", err)
			}
		}
	}
}

func collectMetrics(gatherers []c.MetricsGatherer, logger l.Logger) []c.Metric {

	metrics := []c.Metric{}

	for _, g := range gatherers {
		m, err := g.GetMetrics()
		if err != nil {
			logger.PrintError("Problem getting metrics from gatherer", err)
		}
		metrics = append(metrics, m...)
	}

	return metrics
}

func doVerifyAws(configFileName *string) {
    logger, config := prepare(configFileName, "verify")
    r.VerifyRepeater(c.REPEATER_AWS, logger, config)
}

func doEmailTo(configFileName *string, emailTo *string) {
    logger, config := prepare(configFileName, "argsEmail")
    mailer := s.NewMailer(logger, config)
    t := time.Now()
    subject := fmt.Sprintf("MetricsD command line test - subject %s", t.Format("2006-01-02T15:04:05.999-07:00"))
    body := fmt.Sprintf("MetricsD command line test - body %s", t.Format("2006-01-02T15:04:05.999-07:00"))
    mailer.SendEmail([]string{*emailTo}, subject, body)
    logger.Info("Mail sent to " + *emailTo)
}
