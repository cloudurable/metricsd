package service

import (
    "fmt"
    "net/smtp"
    "crypto/tls"
    "strings"
    c "github.com/cloudurable/metricsd/common"
    l "github.com/cloudurable/simplelog/logging"
)

type Smtp struct {
    logger l.Logger
    smtpConfig *c.SmtpConfig
}

func NewSmtp(logger l.Logger, config *c.Config) *Smtp {
    logger = c.EnsureLogger(logger, config.Debug, "smtp-service")
    return &Smtp{
        logger: logger,
        smtpConfig: &config.SmtpConfig,
    }
}

func (this *Smtp) serverName() string {
    return fmt.Sprintf("%s:%d", this.smtpConfig.Host, this.smtpConfig.Port)
}

func (this *Smtp) BuildMessage(toIds []string, subject string, body string) string {
    message := ""
    message += fmt.Sprintf("From: %s\r\n", this.smtpConfig.FromAddress)
    if len(toIds) > 0 {
        message += fmt.Sprintf("To: %s\r\n", strings.Join(toIds, ";"))
    }

    message += fmt.Sprintf("Subject: %s\r\n", subject)
    message += "\r\n" + body

    return message
}

func (this *Smtp) SendEmail(toIds []string, subject string, body string) bool {
    messageBody := this.BuildMessage(toIds, subject, body)

    auth := smtp.PlainAuth(c.EMPTY, this.smtpConfig.Username, this.smtpConfig.Password, this.smtpConfig.Host)

    tlsconfig := &tls.Config{
        InsecureSkipVerify: this.smtpConfig.IgnoreCert,
        ServerName:         this.smtpConfig.Host,
    }

    conn, err := tls.Dial("tcp", this.serverName(), tlsconfig)

    client, err := smtp.NewClient(conn, this.smtpConfig.Host)
    if err == nil {
        // step 1: Use Auth
        err = client.Auth(auth)
        if err == nil {
            // step 2: add all from and to
            err = client.Mail(this.smtpConfig.FromAddress)
            if err == nil {
                for _, k := range toIds {
                    if err == nil {
                        err = client.Rcpt(k)
                    }
                }
                if err == nil {
                    w, err := client.Data()
                    if err == nil {
                        _, err = w.Write([]byte(messageBody))
                        if err == nil {
                            w.Close()
                        }
                    }
                }
            }
        }
    }

    if err != nil {
        this.logger.Warn(c.ToJsonLabeledString(err))
        return false
    }

    if client != nil {
        err = client.Quit()
        if err != nil {
            this.logger.Warn(c.ToJsonLabeledString(err))
        }
    }
    return true
}
