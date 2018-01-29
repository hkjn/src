// Package probes provides some probe implementations.
//
// This package defines some helpers to send alert emails, while
// actual probes are defined in subpackages.
//
// It uses sendgrid.com to send emails, which seems quite reliable,
// cheap and easy to use. You need an account with sengrid.com with
// user/password to use this package.
//
// To send alert emails, at minimum the following configuration is
// required:
//   - Config.SendGrid: sendgrid credentials
//   - Config.Alert.Recipient: who receives the emails
package probes // import "hkjn.me/src/probes"

import (
	"bytes"
	"fmt"
	"html/template"

	"hkjn.me/src/prober"

	"github.com/golang/glog"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Config is the email configuration.
var Config = EmailConfig{TemplateName: "email"}

// EmailConfig describes the structure of the email configuration.
type EmailConfig struct {
	// Template for HTML email. See EmailData for what's passed to the
	// template when an alert email is generated.
	Template     string
	TemplateName string // name of the template
	Alert        struct {
		Sender    string   // From: address
		Recipient string   // To: address
		CCs       []string // CC: addresses
	}
	SendgridToken string // sendgrid token
}

// EmailData describes the data available in EmailConfig.Template.
type EmailData struct {
	Name, Desc string
	Badness    int
	Records    prober.Records
}

// getClient returns the SendGrid client.
func getClient() (*sendgrid.Client, error) {
	token := Config.SendgridToken
	if token == "" {
		return nil, fmt.Errorf("no sendgrid token specified - set Config.SendgridToken")
	}
	return sendgrid.NewSendClient(token), nil
}

// SendAlertEmail sends an alert email using SendGrid.
//
// This is provided to simplify prober.Probe implementations for Alert().
func SendAlertEmail(name, desc string, badness int, records prober.Records) error {
	glog.V(1).Infof("sending alert email..\n")

	var t *template.Template
	var err error
	if Config.Alert.Recipient == "" {
		err = fmt.Errorf("missing email recipient - set Config.Alert.Recipient")
	}
	// TODO: would be nice to only parse the template once, but little
	// complicated without complicating config for user, since this will
	// be called from separate goroutines..
	t, err = template.New(Config.TemplateName).Parse(Config.Template)
	if err != nil {
		err = fmt.Errorf("failed to parse email: %v", err)
	}
	if err != nil {
		return err
	}

	var html bytes.Buffer
	data := EmailData{name, desc, badness, records}
	err = t.ExecuteTemplate(
		&html, Config.TemplateName, data)
	if err != nil {
		return fmt.Errorf("failed to construct email from template: %v", err)
	}

	sgClient, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create mail client: %v", err)
	}

	subject := fmt.Sprintf("%s failed (badness %d)", name, badness)
	from := mail.NewEmail("Alert email sender", Config.Alert.Sender)
	to := mail.NewEmail("Alert email recipient", Config.Alert.Recipient)
	email := mail.NewSingleEmail(from, subject, to, html.String(), html.String())

	resp, err := sgClient.Send(email)
	if err != nil {
		return fmt.Errorf("failed to send mail: %v", err)
	}
	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		return fmt.Errorf("unexpected response from sendgrid, want status 200-299: %v", resp)
	}
	glog.Infof("sent alert email to %s\n", Config.Alert.Recipient)
	return nil
}
