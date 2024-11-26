package model

type AccountReport struct {
	Sum                  float64
	TransactionsPerMonth map[string]int
	AvgDebit             float64
	AvgCredit            float64
}

type SendEmailRequest struct {
	From                Email               `json:"from"`
	Personalizations    Personalizations    `json:"Personalizations"`
	DynamicTemplateData DynamicTemplateData `json:"dynamic_template_data"`
	TemplateID          string              `json:"template_id"`
}

type Email struct {
	Email string `json:"email"`
}

type Personalizations struct {
	To Email `json:"to"`
}

type DynamicTemplateData struct {
	TotalBalance         float64       `json:"total_balance"`
	TransactionsPerMonth []Transaction `json:"transactions_per_month"`
	AverageDebit         float64       `json:"average_debit"`
	AverageCredit        float64       `json:"average_credit"`
	SenderName           string        `json:"Sender_Name"`
	SenderAddress        string        `json:"Sender_Address"`
	SenderCity           string        `json:"Sender_City"`
	SenderState          string        `json:"Sender_State"`
}

type Transaction struct {
	Month string `json:"month"`
	Count int    `json:"count"`
}

type EmailConfig struct {
	To         string
	Report     AccountReport
	TemplateID string
}

type EmailPayload struct {
	To      string
	Payload SendEmailRequest
}
