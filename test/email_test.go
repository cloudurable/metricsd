package test

import (
    c "github.com/cloudurable/metricsd/common"
    s "github.com/cloudurable/metricsd/service"
    "testing"
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
    mailer.SendEmail([]string{"scottfauerbach@gmail.com"}, "This is the subject 42", "This is a test 42")
}
