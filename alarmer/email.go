package alarmer

import (
	lg "github.com/cloudurable/simplelog/logging"
    c "github.com/cloudurable/metricsd/common"
    s "github.com/cloudurable/metricsd/service"
    "time"
    "fmt"
    "os"
    "crypto/md5"
    "encoding/hex"
)

type EmailMetricAlarmer struct {
	logger               lg.Logger
    config               *c.Config
    awsConfig            *s.AwsConfig
    alarmMap             map[string]time.Time
    to                   []string
    alarmIntervalSeconds int64
}

func (al EmailMetricAlarmer) Verify() bool {
    t := time.Now()
    subject := fmt.Sprintf("MetricsD verify %s", t.Format(c.STD_TIME_FORMAT))
    body := fmt.Sprintf("This is a MetricsD verification email, sent from '%s' at %s", al.awsConfig.EC2InstanceId, t.Format(c.STD_TIME_FORMAT))
    return al.send(subject, body)
}

func (al EmailMetricAlarmer) ProcessMetrics(context c.MetricContext, metrics []c.Metric) error {
    if len(metrics) > 0 {
        newMap := map[string]time.Time{}
        for _, m := range metrics {
            if m.Alarm {
                hash := hash(m)
                lastTime, present := al.alarmMap[hash]
                shouldSend := !present
                if present {
                    elapsed := time.Now().Unix() - lastTime.Unix()
                    if elapsed > al.alarmIntervalSeconds {
                        shouldSend = true
                    } else {
                        newMap[hash] = lastTime
                    }
                }

                if shouldSend {
                    subject := fmt.Sprintf("MetricsD alarm.  Instance'%s'  Provider'%s'  Name'%s'", al.awsConfig.EC2InstanceId, m.Provider, m.Name)
                    m.When.Format(c.STD_TIME_FORMAT)
                    body := "MetricsD alarm\r\n" + m.MetricFormattedString()
                    al.logger.Debug("Sending alarm email:", subject)
                    if al.send(subject, body) {
                        newMap[hash] = time.Now()
                        al.alarmMap[hash] = newMap[hash]
                    } else {
                        // didn't send make sure it will if it alarms again
                        delete(al.alarmMap, hash)
                    }
                }
            }
        }

        // NOT SURE WHY I HAVE TO DO THIS, I THINK IT HAS TO DO WITH POINTERS. AND YES, I JUST TRIED ASSIGNING
        deleteThese := []string{}
        for key, _ := range al.alarmMap {
            if _, present := newMap[key]; !present {
                deleteThese = append(deleteThese, key)
            }
        }
        for _,key := range deleteThese {
            delete(al.alarmMap, key)
        }
    }
    return nil
}

func hash(m c.Metric) string {
    hasher := md5.New()
    hasher.Write([]byte(fmt.Sprintf("%d %d %s %s", m.Type, m.Source, m.Name, m.Provider)))
    return hex.EncodeToString(hasher.Sum(nil));
}

func NewEmailMetricAlarmer(config *c.Config) *EmailMetricAlarmer {
    logger := c.EnsureLogger(nil, config.Debug, "email")

    alarmTo := 	c.ReadConfigStringArray("email alarm to", config.EmailAlarmTo, []string{}, logger)
    if len(alarmTo) == 0 {
        logger.Error("No email alarm To addresses are configured")
        os.Exit(11)
    }

    alarmIntervalSeconds := config.EmailAlarmIntervalSeconds
    if alarmIntervalSeconds < 30 {
        alarmIntervalSeconds = 30
    } else if alarmIntervalSeconds > 3600 {
        alarmIntervalSeconds = 3600
    }

    awsConfig := s.GetAwsConfig(config)

    return &EmailMetricAlarmer{logger, config, awsConfig, map[string]time.Time{}, alarmTo, int64(alarmIntervalSeconds)}
}

func (al EmailMetricAlarmer) send(subject string, body string) bool {
    mailer := s.NewMailer(al.logger, al.config)
    return mailer.SendEmail(al.to, subject, body)
}
