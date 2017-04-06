package test

import (
    c "github.com/cloudurable/metricsd/common"
    s "github.com/cloudurable/metricsd/service"
    "testing"
    "time"
    "fmt"
)

func TestEmail(test *testing.T) {

    config := c.Config{
        Debug: true,
        SmtpConfig: c.SmtpConfig{
            Host: "gator3128.hostgator.com",
            Port: 465,
            Username: "donotreply@arondight.com",
            Password: "d0n0treply",
            FromAddress: "donotreply@arondight.com",
            IgnoreCert: false,
        },
    }

    mailer := s.NewSmtp(nil, &config)
    t := time.Now()
    subject := fmt.Sprintf("MetricsD command line test subject %s", t.Format("Mon, Jan 02 2006 15:04:05 -0700"))
    body := fmt.Sprintf("MetricsD command line test body %s", t.Format("Mon, Jan 02 2006 15:04:05 -0700"))
    mailer.SendEmail([]string{"scottfauerbach@gmail.com"}, subject, body)
}
