package smtp

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"

	"github.com/circa10a/go-aws-news/providers"

	"github.com/circa10a/go-aws-news/news"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type config struct {
	Providers struct {
		Provider Provider `yaml:"smtp"`
	} `yaml:"providers"`
}

// Provider is an implementation of the `go-aws-news/providers` Provider interface.
type Provider struct {
	IsEnabled bool     `yaml:"enabled"`
	Server    string   `yaml:"server"`
	Port      string   `yaml:"port"`
	Username  string   `yaml:"username"`
	Password  string   `yaml:"password"`
	Subject   string   `yaml:"subject"`
	From      string   `yaml:"from"`
	To        []string `yaml:"to"`
}

type email struct {
	addr    string
	from    string
	to      []string
	subject string
	body    string
}

// init initializes the provider from the provided config.
func init() {
	var c config
	if err := yaml.Unmarshal(providers.Config, &c); err != nil {
		log.Fatal(err)
	}

	providers.RegisterProvider(&c.Providers.Provider)
}

// Enabled returns true if the provider is enabled in the config.
func (p *Provider) Enabled() bool {
	return p.IsEnabled
}

// GetName returns the Provider's name.
func (*Provider) GetName() string {
	return "smtp"
}

// Notify is the function executed to POST to a provider's webhook url.
func (p *Provider) Notify(news news.Announcements) {

	m := &email{
		addr:    fmt.Sprintf("%s:%s", p.Server, p.Port),
		from:    p.From,
		to:      p.To,
		subject: p.Subject,
		body:    news.HTML(),
	}

	var auth smtp.Auth

	if p.Username != "" && p.Password != "" {
		auth = smtp.PlainAuth("", p.Username, p.Password, p.Server)
	}

	m.parseTemplate(defaultTemplate)

	log.Info(fmt.Sprintf("[%v] Firing notification", p.GetName()))
	err := m.sendMail(auth)
	if err != nil {
		log.Error(fmt.Sprintf("[%v] %v", p.GetName(), err))
	}
}

func (e *email) sendMail(auth smtp.Auth) error {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("Subject: %s\n", e.subject))
	b.WriteString("MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n\n")
	b.WriteString(e.body)

	err := smtp.SendMail(e.addr, auth, e.from, e.to, b.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (e *email) parseTemplate(tmpl string) error {
	t, err := template.New("mail").Parse(tmpl)
	if err != nil {
		return err
	}

	var b bytes.Buffer
	err = t.Execute(&b, struct{ News string }{e.body})
	if err != nil {
		return err
	}

	// wrap body with parsed template
	e.body = b.String()
	return nil
}

var defaultTemplate = `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html>
</head>
<body>
{{ .News }}
</body>
</html>
`
