package blog

import (
	"log"
	"time"
	"bytes"
	"strings"
	"text/template"
	"appengine"
	"appengine/mail"
	"code.google.com/p/d2b/data"
)

func PublishEmail (c appengine.Context, message string) {
	sender := data.GetConfigOption(c, "email.sender");
	// comma delimited list of emails
	var blog_addresses []string;
	blog_addresses = strings.Split(data.GetConfigOption(c, "email.addresses"), ",");

	// template
	subject := data.GetConfigOption(c, "email.subject_template");
	if subject == "" {
		subject = "Links for {{ .Day }} {{ .Month }}, {{ .Year }}";
	}
	templ_subject := template.Must(template.New("subject").Parse(subject));
	var subject_buffer bytes.Buffer;
	templ_subject.Execute(&subject_buffer, time.Now())

	msg := &mail.Message{
		Sender: sender,
		To: blog_addresses,
		Subject: subject_buffer.String(),
		HTMLBody: message,
	}

	err := mail.Send(c, msg);
	if err != nil {
		log.Fatal("Error publishing to email: ", err);
	}

	data.SetConfigOption(c, "lastPublished", time.Now().Format(time.RFC3339));
}
