package util

import (
	"crypto/tls"
	"fmt"

	"gopkg.in/gomail.v2"
)

const (
	from     = "From"
	to       = "To"
	subject  = "Subject"
	textType = "text/html"
)

var (
	Smtp = ""
	Port = 465
	User = ""
	Pwd  = ""
)

func SendEmail(toMail, subjectCon, body string) error {
	if !isValidSend() {
		return fmt.Errorf("invalide email setup to sendEmail")
	}

	msg := gomail.NewMessage()
	msg.SetHeader(from, User)
	msg.SetHeader(to, toMail)
	msg.SetHeader(subject, subjectCon)
	msg.SetBody(textType, body)

	//"Grj5DmU4q2VmZj4"
	d := gomail.NewDialer(Smtp, Port, User, Pwd)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	err := d.DialAndSend(msg)
	return err
}

func isValidSend() bool {

	if Smtp != "" && User != "" && Pwd != "" {
		return true
	} else {
		return false
	}
}
