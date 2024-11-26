package services

import (
	model "stori-account-summary/model"
)

type SendGridClient interface {
	Send(ec model.EmailPayload) error
}

type EmailSender struct {
	SendGridClient SendGridClient
	From           string
}

func New(client SendGridClient) EmailSender {
	return EmailSender{
		SendGridClient: client,
	}
}

func (es EmailSender) SendMail(config model.EmailConfig) error {
	body := es.body(config)
	payload := model.EmailPayload{
		To:      config.To,
		Payload: body,
	}

	err := es.SendGridClient.Send(payload)
	if err != nil {
		return err
	}

	return nil
}

func (es EmailSender) body(config model.EmailConfig) model.SendEmailRequest {
	const (
		senderAddress = "Calle 81 # 11 - 08"
		senderName    = "Stori"
		senderCity    = "Bogota"
		senderState   = "Colombia"
	)

	return model.SendEmailRequest{
		From: model.Email{
			Email: es.From,
		},
		Personalizations: model.Personalizations{
			To: model.Email{
				Email: config.To,
			},
		},
		DynamicTemplateData: model.DynamicTemplateData{
			TotalBalance:         config.Report.Sum,
			AverageDebit:         config.Report.AvgDebit,
			AverageCredit:        config.Report.AvgCredit,
			SenderName:           senderName,
			SenderAddress:        senderAddress,
			SenderCity:           senderCity,
			SenderState:          senderState,
			TransactionsPerMonth: config.Report.TransactionPerMonth(),
		},
		TemplateID: config.TemplateID,
	}
}
