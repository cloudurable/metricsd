package service

import (
    "fmt"
    "net/smtp"
    "crypto/tls"
    "strings"
    c "github.com/cloudurable/metricsd/common"
    l "github.com/cloudurable/simplelog/logging"
)

type Mailer struct {
    logger l.Logger
    senderId string
    host string
    port int
    username string
    password string
    ignoreCert bool
}

func NewMailer(logger l.Logger, config *c.Config) *Mailer {
    logger = c.EnsureLogger(logger, config.Debug, "Mailer")
    return &Mailer{
        logger: logger,
        senderId: config.SmtpFromAddress,
        host: config.SmtpHost,
        port: config.SmtpPort,
        username: config.SmtpUsername,
        password: config.SmtpPassword,
        ignoreCert: config.SmtpIgnoreCert,
    }
}

func (mailer *Mailer) serverName() string {
    return fmt.Sprintf("%s:%d", mailer.host, mailer.port)
}

func (mailer *Mailer) BuildMessage(toIds []string, subject string, body string) string {
    message := ""
    message += fmt.Sprintf("From: %s\r\n", mailer.senderId)
    if len(toIds) > 0 {
        message += fmt.Sprintf("To: %s\r\n", strings.Join(toIds, ";"))
    }

    message += fmt.Sprintf("Subject: %s\r\n", subject)
    message += "\r\n" + body

    return message
}

func (mailer *Mailer) SendEmail(toIds []string, subject string, body string) {
    messageBody := mailer.BuildMessage(toIds, subject, body)

    auth := smtp.PlainAuth(c.EMPTY, mailer.username, mailer.password, mailer.host)

    tlsconfig := &tls.Config{
        InsecureSkipVerify: mailer.ignoreCert,
        ServerName: mailer.host,
    }

    conn, err := tls.Dial("tcp", mailer.serverName(), tlsconfig)

    client, err := smtp.NewClient(conn, mailer.host)
    if err == nil {
        // step 1: Use Auth
        err = client.Auth(auth)
        if err == nil {
            // step 2: add all from and to
            err = client.Mail(mailer.senderId)
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
        mailer.logger.Warn(c.ObjectToString(err))
    } else if client != nil {
        err = client.Quit()
        if err != nil {
            mailer.logger.Warn(c.ObjectToString(err))
        }
    }
}
