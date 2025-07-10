package email

import (
	Coms "github.com/benjaminRoberts01375/Go-Communicate"
)

type config struct {
	EmailAPIKey        string `json:"email_api_key"`
	AllowSendingEmails bool   `json:"allow_sending_emails"`
	EmailFrom          string `json:"email_from"`
	BaseURL            string `json:"base_url"`
	ServiceName        string `json:"service_name"`
}

var Config config

func (config) URL() string {
	return `http://` + Config.BaseURL + `/api/`
}

func Setup() {
	Coms.ReadExternalConfig("email.json", &Config)
	if Config.EmailAPIKey == "" {
		panic("No email API key specified")
	}
}
