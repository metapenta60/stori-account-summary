package pkg

import (
	"encoding/json"
	"log"
	model "stori-account-summary/model"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
)

const (
	sendEmailURL = "/v3/mail/send"
)

type SendGridClient struct {
	ProviderHost string
	ApiKey       string
	EmailFrom    string
}

func NewSendGridClient(providerHost, apiKey string) SendGridClient {
	return SendGridClient{
		ProviderHost: providerHost,
		ApiKey:       apiKey,
	}
}

func (sgc SendGridClient) Send(payload model.EmailPayload) error {
	request, err := sgc.request(payload)
	if err != nil {
		return err
	}

	res, err := sendgrid.API(*request)

	log.Print(res)
	if err != nil {
		return err
	}
	return nil
}

func (sgc SendGridClient) request(mp model.EmailPayload) (*rest.Request, error) {
	request := sendgrid.GetRequest(sgc.ApiKey, sendEmailURL, sgc.ProviderHost)
	request.Method = "POST"
	jsonData, err := json.Marshal(mp.Payload)
	if err != nil {
		return nil, err
	}

	request.Body = jsonData
	return &request, nil
}
