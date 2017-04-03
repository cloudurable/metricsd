package common

import (
    //"net/smtp"
    //"log"
)
import "fmt"

func SendEmail(to []string, subject string, body string) {

    //auth := smtp.PlainAuth("", "donotreply@arondight.com", "d0n0treply", "gator3128.hostgator.com")

    // Connect to the server, authenticate, set the sender and recipient,
    // and send the email all in one step.
    //msg := []byte(
    //    "To: " + to[0] + "\r\n" +
    //    "Subject: " + subject + "\r\n" +
    //    "\r\n" + body + "\r\n")
    //
    //err := smtp.SendMail("gator3128.hostgator.com:25", auth, "donotreply@arondight.com", to, msg)
    //if err != nil {
    //    log.Fatal(err)
    //}

    fmt.Print("EMAIL NOT WORKING YET!!!!")
}