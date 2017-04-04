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
        SmtpHost: "gator3128.hostgator.com",
        SmtpPort: 465,
        SmtpUsername: "donotreply@arondight.com",
        SmtpPassword: "d0n0treply",
        SmtpFromAddress: "donotreply@arondight.com",
        SmtpIgnoreCert: false,
    }

    mailer := s.NewMailer(nil, &config)
    t := time.Now()
    subject := fmt.Sprintf("MetricsD command line test subject %s", t.Format("Mon, Jan 02 2006 15:04:05 -0700"))
    body := fmt.Sprintf("MetricsD command line test body %s", t.Format("Mon, Jan 02 2006 15:04:05 -0700"))
    mailer.SendEmail([]string{"scottfauerbach@gmail.com"}, subject, body)
}
