package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/oyvindsk/post2mail"
)

type server struct {
	emailTo       string
	smtpInfo      post2mail.SMTPInfo
	indexTemplate *template.Template
}

func main() {

	// A server type to share some more or less global data with our handlers
	s := server{}
	s.indexTemplate = template.Must(template.New("index").Parse(templateStr))

	//
	// Read command line parameters

	// SMTP info
	flag.StringVar(&s.smtpInfo.Server, "smtp-server", "", "SMTP server to send mail through (gmail is: 'smtp.gmail.com')")
	flag.IntVar(&s.smtpInfo.Port, "smtp-port", 587, "Port to connect to on the smtp server, usually 25 or 587")
	flag.StringVar(&s.smtpInfo.Username, "smtp-user", "", "Username for the SMTP server")
	flag.StringVar(&s.smtpInfo.Password, "smtp-pass", "", "Password for the SMTP server")

	// Email info
	flag.StringVar(&s.emailTo, "email-to", "", "Email address to send the emails to")

	// Other parameters
	var httpPort int
	flag.IntVar(&httpPort, "http-port", 80, "The port the http server should be listening on (no https support)")

	flag.Parse()

	// Ensure that we have required commandline parameters, there's no default for these ATM
	if s.smtpInfo.Server == "" || s.emailTo == "" {
		flag.Usage()
		log.Fatal("-smtp-server and -email-to are required")
	}

	//
	// HTTP handlers
	http.Handle("/", http.HandlerFunc(s.handleGet))      // Get, for testing
	http.Handle("/post", http.HandlerFunc(s.handlePost)) // Receive the POST from the form

	// HTTP listen
	log.Printf("About to listen on: :%d", httpPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil) // Remeber to use the same port in the <form> :)
	if err != nil {
		log.Fatalf("ListenAndServe: %s", err)
	}
}

// handleGet just returns a simple html form - for testing
func (s server) handleGet(w http.ResponseWriter, req *http.Request) {
	s.indexTemplate.Execute(w, req.FormValue("s")) // FIXME some error handling
}

// Status returned to the HTTP client
// modifying this will affect the json returned to the client
// if so, remember the hardcoded error json string as well
type statusMSG struct {
	Success bool
	Status  string
}

func (s server) handlePost(w http.ResponseWriter, req *http.Request) {
	log.Println("handlePost")
	// log.Println(req.Form)

	var email post2mail.EmailData

	email.To = s.emailTo // Receiving email address, sat at server start

	email.FromName = req.FormValue("name")
	email.FromEmail = req.FormValue("from")
	email.Subject = "Someone filled out your form: " + req.FormValue("subject")
	email.Text = req.FormValue("text")

	// Do stupid spam filtering :P
	spam, reason := post2mail.IsSpam(email)
	if spam {
		log.Printf("handleContactform: skipping spammy post: reason: %q, Refferer: %q, IP: %q, UA: %q", reason, req.Referer(), req.RemoteAddr, req.UserAgent())
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, `{ "Status" : "Not acceptable" , "Success" : "false" }`)
		return
	}

	// Send the info on email
	err := post2mail.FormatAndSendEmail(email, s.smtpInfo)

	// Return some infor to the client
	var status statusMSG
	if err == nil {
		status.Success = true
		status.Status = "OK"
	} else {
		status.Success = false
		status.Status = fmt.Sprintf("Error: %s", err)
	}

	// Return a json status
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(status)
	// json encoding errors
	if err != nil {
		log.Printf("handlePost: %d: failed to handle request: json encode failed: %s. Status: %v", http.StatusInternalServerError, err, status)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{ "Status" : "error: json encode failed: %s" , "Success" : "false" }`, err)
		return
	}

	// other problems
	if !status.Success {
		log.Printf("handlePost: %d: failed to handle request: Status: %v", http.StatusInternalServerError, status)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", j)
		return
	}

	_, err = w.Write(j)
	if err != nil {
		log.Printf("handlePost: failed write to client: error: %s. Status: %v", err, status)
	}

	// Success!
	log.Printf("Email sent:\n\t%+v", status)

}

const templateStr = `
<html>
<head>
<title>test</title>
</head>
<body>
{{if .}}
<h1>?</h1>
{{.}}
<br>
<br>
{{end}}
<form action="/post" name=f method="POST">
    <input maxLength=1024 size=70   name=name       value="" placeholder="Name"    title="Name">
    <input maxLength=1024 size=70   name=from       value="" placeholder="From"    title="From Email">
    <input maxLength=1024 size=70   name=subject    value="" placeholder="Subject" title="Subject">
    <input maxLength=1024 size=500  name=text       value="" placeholder="Message" title="Message">
    <input type=submit value="Send" name=send>
</form>

</body>
</html>
`
