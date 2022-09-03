package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	email := createEmail()
	errs := email.SendMessage(msg)
	if errs != nil {
		app.errorJSON(w, errs)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func createEmail() Mail {
	ports, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	email := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        ports,
		Username:    os.Getenv("MAIL_USERNAME"),
		password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
	}

	return email
}

type Mail struct {
	Domain      string `json:"name"`
	Host        string `json:"data"`
	Port        int
	Username    string
	password    string
	Encryption  string
	FromAddress string
	FromName    string
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	Datamap     map[string]any
}

func (m *Mail) SendMessage(msg Message) error {
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	data := map[string]any{
		"message": msg.Data,
	}

	msg.Datamap = data

	// formattedMessage, err := m.buildHTMLMessage(msg)
	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	SMTPClient, err := server.Connect()
	if err != nil {
		log.Println(err)
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)
	email.SetBody(mail.TextPlain, fmt.Sprintf("email %s", msg.From))
	email.AddAlternative(mail.TextHTML, fmt.Sprintf("email %s", msg.From))

	errs := email.Send(SMTPClient)
	if errs != nil {
		log.Println(errs)
		return errs
	}
	return nil
}

func (m *Mail) getEncryption(s string) mail.Encryption {
	if s == "tls" {
		return mail.EncryptionSTARTTLS
	} else if s == "ssl" {
		return mail.EncryptionSSL
	} else if s == "none" {
		return mail.EncryptionNone
	} else {
		return mail.EncryptionSSL
	}
}
