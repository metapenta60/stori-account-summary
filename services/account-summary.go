package services

import (
	"stori-account-summary/model"
)

type accountSummaryService struct {
	rows model.Rows
}

func NewReportService(rows model.Rows) accountSummaryService {
	return accountSummaryService{
		rows: rows,
	}
}

func (as accountSummaryService) createEmptyReport() model.AccountReport {
	return model.AccountReport{
		Sum: 0.0,
		TransactionsPerMonth: map[string]int{
			"1":  0,
			"2":  0,
			"3":  0,
			"4":  0,
			"5":  0,
			"6":  0,
			"7":  0,
			"8":  0,
			"9":  0,
			"10": 0,
			"11": 0,
			"12": 0,
		},
		AvgDebit:  0.0,
		AvgCredit: 0.0,

		NumCreditTransaction: 0.0,
		NumDebitTransaction:  0.0,
	}
}

func (as accountSummaryService) AnalyseAccount() model.AccountReport {
	report := as.createEmptyReport()

	for _, row := range as.rows {
		report.AddTransaction(row.Transaction)
		report.IncreaseTransactionCount(row.Date.Month)
	}

	report.UpdateAverageDebitAndCredit()

	return report
}
