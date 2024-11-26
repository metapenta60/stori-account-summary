package services

import (
	"testing"

	mock_sendemail "stori-account-summary/mocks"
	"stori-account-summary/model"

	"go.uber.org/mock/gomock"
)

func TestEmailSender_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSendGridClient := mock_sendemail.NewMockSendGridClient(ctrl)

	emailSender := EmailSender{
		SendGridClient: mockSendGridClient,
		From:           "sender@example.com",
	}

	sendEmailRequest := model.SendEmailRequest{
		From: model.Email{
			Email: emailSender.From,
		},
		Personalizations: model.Personalizations{
			To: model.Email{
				Email: "recipient@example.com",
			},
		},
		DynamicTemplateData: model.DynamicTemplateData{
			TotalBalance: 12345.67,
			TransactionsPerMonth: []model.Transaction{
				{Month: "January", Count: 10},
				{Month: "February", Count: 15},
				{Month: "March", Count: 8},
			},
			AverageDebit:  500.75,
			AverageCredit: -650.90,
			SenderName:    "Stori",
			SenderAddress: "Calle 81 # 11 - 08",
			SenderCity:    "Bogota",
			SenderState:   "Colombia",
		},
		TemplateID: "d-12345abcde67890fghij12345klmn67890",
	}

	emailPayload := model.EmailPayload{
		To:      "recipient@example.com",
		Payload: sendEmailRequest,
	}

	emailConfig := model.EmailConfig{
		To: "recipient@example.com",
		Report: model.AccountReport{
			Sum: 12345.67,
			TransactionsPerMonth: map[string]int{
				"1": 10,
				"2": 15,
				"3": 8,
			},
			AvgDebit:  500.75,
			AvgCredit: -650.90,
		},
		TemplateID: "d-12345abcde67890fghij12345klmn67890",
	}

	mockSendGridClient.EXPECT().Send(emailPayload).Return(nil)

	err := emailSender.SendMail(emailConfig)
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
}
