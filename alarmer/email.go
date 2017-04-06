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
	logger                lg.Logger
    config                *c.Config
    awsContext            *c.AwsContext
    alarmMap              map[string]time.Time
    to                    []string
    resendIntervalSeconds int64
}

func (this *EmailMetricAlarmer) Verify() bool {
    t := time.Now()
    subject := fmt.Sprintf("MetricsD verify %s", t.Format(c.STD_TIME_FORMAT))
    body := fmt.Sprintf("This is a MetricsD verification email, sent from '%s' at %s", this.awsContext.EC2InstanceId, t.Format(c.STD_TIME_FORMAT))
    return this.send(subject, body)
}

func (this *EmailMetricAlarmer) Alarm(metrics []c.Metric) error {
    if len(metrics) > 0 {
        newSendTimes := map[string]time.Time{}
        for _, m := range metrics {
            if !m.Alarm.Empty() {
                hash := hash(m)
                lastSendTime, sentThisBefore := this.alarmMap[hash]
                newSendTime := lastSendTime
                shouldSend := !sentThisBefore
                if sentThisBefore {
                    if this.resendIntervalSeconds > 0 { // if it's < 0, never re-send
                        elapsed := time.Now().Unix() - lastSendTime.Unix()
                        if elapsed > this.resendIntervalSeconds {
                            shouldSend = true
                            newSendTime = time.Now()
                        }
                    }
                } else {
                    newSendTime = time.Now()
                }
                newSendTimes[hash] = newSendTime

                if shouldSend {
                    subject := fmt.Sprintf("MetricsD alarm.  Instance'%s'  Provider'%s'  Name'%s'", this.awsContext.EC2InstanceId, m.Provider, m.Name)
                    m.When.Format(c.STD_TIME_FORMAT)
                    body := "MetricsD alarm\r\n" + m.MetricStringFormatted()
                    this.logger.Debug("Sending alarm email:", subject)
                    if this.send(subject, body) {
                        this.alarmMap[hash] = newSendTimes[hash]
                    } else {
                        // didn't send make sure it will if it alarms again
                        delete(this.alarmMap, hash)
                    }
                }
            }
        }

        // NOT SURE WHY I HAVE TO DO THIS, I THINK IT HAS TO DO WITH POINTERS. AND YES, I JUST TRIED ASSIGNING
        deleteThese := []string{}
        for key, _ := range this.alarmMap {
            if _, present := newSendTimes[key]; !present {
                deleteThese = append(deleteThese, key)
            }
        }
        for _,key := range deleteThese {
            delete(this.alarmMap, key)
        }
    }
    return nil
}

func hash(m c.Metric) string {
    hasher := md5.New()
    hasher.Write([]byte(fmt.Sprintf("%d %d %s %s", m.Type, m.Source, m.Name, m.Provider)))
    return hex.EncodeToString(hasher.Sum(nil));
}
func (this *EmailMetricAlarmer) Name() string { return this.logger.Name() }

func NewEmailMetricAlarmer(config *c.Config) *EmailMetricAlarmer {
    logger := c.EnsureLogger(nil, config.Debug, c.ALARMER + c.ALARMER_AWS)

    alarmTo := 	c.ReadConfigStringArray("email alarm to", []string{"FIX ME !!! config.EmailAlarmTo"}, []string{}, logger, true)
    if len(alarmTo) == 0 {
        logger.Error("No email alarm To addresses are configured")
        os.Exit(501)
    }

    resendIntervalSeconds := config.EmailAlarmerConfig.ResendIntervalSeconds
    if resendIntervalSeconds == 0 {
        resendIntervalSeconds = 300
    } else if resendIntervalSeconds > 3600 {
        resendIntervalSeconds = 3600
    }

    awsContext := c.GetAwsContext(logger, config)

    return &EmailMetricAlarmer{logger, config, awsContext, map[string]time.Time{}, alarmTo, int64(resendIntervalSeconds)}
}

func (this *EmailMetricAlarmer) send(subject string, body string) bool {
    mailer := s.NewSmtp(this.logger, this.config)
    return mailer.SendEmail(this.to, subject, body)
}
