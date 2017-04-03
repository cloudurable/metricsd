package common

import (
    "fmt"
    "net/smtp"
    "crypto/tls"
    "strings"
)

type Mail struct {
    senderId string
    toIds    []string
    subject  string
    body     string
}

type SmtpServer struct {
    host string
    port string
}

func (s *SmtpServer) ServerName() string {
    return s.host + ":" + s.port
}

func (mail *Mail) BuildMessage() string {
    message := ""
    message += fmt.Sprintf("From: %s\r\n", mail.senderId)
    if len(mail.toIds) > 0 {
        message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";"))
    }

    message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
    message += "\r\n" + mail.body

    return message
}

func SendEmail(to []string, subject string, body string) {
    mail := Mail{}
    mail.senderId = "donotreply@arondight.com"
    mail.toIds = to
    mail.subject = subject
    mail.body = body

    messageBody := mail.BuildMessage()

    smtpServer := SmtpServer{host: "gator3128.hostgator.com", port: "465"}

    auth := smtp.PlainAuth(EMPTY, "donotreply@arondight.com", "d0n0treply", "gator3128.hostgator.com")

    tlsconfig := &tls.Config{
        InsecureSkipVerify: true,
        ServerName:         smtpServer.host,
    }

    conn, err := tls.Dial("tcp", smtpServer.ServerName(), tlsconfig)
    if err != nil {
        // log something
    }

    client, err := smtp.NewClient(conn, smtpServer.host)
    if err != nil {
        fmt.Println(ObjectToString(err))
    }

    // step 1: Use Auth
    if err = client.Auth(auth); err != nil {
        fmt.Println(ObjectToString(err))
    }

    // step 2: add all from and to
    if err = client.Mail(mail.senderId); err != nil {
        fmt.Println(ObjectToString(err))
    }

    for _, k := range mail.toIds {
        if err = client.Rcpt(k); err != nil {
            fmt.Println(ObjectToString(err))
        }
    }

    // Data
    w, err := client.Data()
    if err != nil {
        fmt.Println(ObjectToString(err))
    }

    _, err = w.Write([]byte(messageBody))
    if err != nil {
        fmt.Println(ObjectToString(err))
    }

    err = w.Close()
    if err != nil {
        fmt.Println(ObjectToString(err))
    }

    client.Quit()
}
