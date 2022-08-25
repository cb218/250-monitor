package probers

import (
	pb "250-monitor/pinger"
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"time"
)

type Monitor struct {
	lastReachable time.Time
}

type EmailTemplate struct {
	HostIP    string
	Timestamp string
}

func NewMonitor() *Monitor {
	monitor := new(Monitor)
	monitor.lastReachable = time.Time{}
	return monitor
}

func (monitor *Monitor) Start(IP string, period int) {
	fmt.Printf("Starting periodic pinger...\n")
	replyCh := make(chan *pb.PingReply)
	var p *pb.PingReply
	var emailData EmailTemplate
	emailData.HostIP = IP
	ticker := time.NewTicker(time.Duration(period) * time.Minute)
	go func() {
		for range ticker.C {
			fmt.Printf("\nPinging host @ %s...\n", IP)
			icmp := NewICMPProbe(IP, replyCh)
			icmp.Start()
			p = <-replyCh
			if p.Reachable == true {
				monitor.lastReachable = time.Now()
				fmt.Printf("Host %s was reached at %s\n", IP, monitor.lastReachable.String())
			} else {
				fmt.Printf("Host %s could not be reached.\n", IP)
				subject := emailData.HostIP + " could not be reached at " + emailData.Timestamp
				Mail(subject, emailData)
			}
		}
	}()
}

const MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}
func (r *Request) SendEmail() (bool, error) {
	mime := MIME
	subject := "Subject: " + r.subject + "\n"
	msg := []byte(subject + mime + "\n" + r.body)

	if err := smtp.SendMail(os.Getenv("HOST")+":"+os.Getenv("MAIL_PORT"), smtp.PlainAuth("", os.Getenv("EMAIL"), os.Getenv("PASSWORD"), os.Getenv("HOST")), os.Getenv("EMAIL"), r.to, msg); err != nil {
		fmt.Printf("Error sending email: %s", err)
		return false, err
	}
	return true, nil
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		fmt.Printf("Error at ParseFiles %s", err)
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		fmt.Printf("Error at Execute %s", err)
		return err
	}
	r.body = buf.String()
	return nil
}

func Mail(subject string, emailData EmailTemplate) {
	r := NewRequest([]string{os.Getenv("DESTINATION")}, subject, "body")
	if err := r.ParseTemplate("template.html", emailData); err == nil {
		r.SendEmail()
		fmt.Printf("Email sent %s\n", subject)
	}
}
