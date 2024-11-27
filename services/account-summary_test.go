package services

import (
	"stori-account-summary/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test for account-summary

func TestAnalyseAccount(t *testing.T) {
	// Configurar los datos de entrada
	rows := model.Rows{
		{Transaction: 100.0,
			Date: model.Date{
				Day:   "01",
				Month: "01",
			}},
		{Transaction: -50.0,
			Date: model.Date{
				Day:   "02",
				Month: "01",
			}},
		{Transaction: 200.0,
			Date: model.Date{
				Day:   "01",
				Month: "02",
			},
		},
	}

	as := accountSummaryService{
		rows: rows,
	}
	report := as.AnalyseAccount()

	expectedSum := 250.0
	assert.Equal(t, expectedSum, report.Sum)

	expectedTotalCredit := 300.0
	assert.Equal(t, expectedTotalCredit, report.TotalCredit)

	expectedTotalDebit := -50.0
	assert.Equal(t, expectedTotalDebit, report.TotalDebit)

	expectedNumCreditTransaction := 2.0
	assert.Equal(t, expectedNumCreditTransaction, report.NumCreditTransaction)

	expectedNumDebitTransaction := 1.0
	assert.Equal(t, expectedNumDebitTransaction, report.NumDebitTransaction)

	expectedAvgDebit := -50.0
	assert.Equal(t, expectedAvgDebit, report.AvgDebit)

	expectedAvgCredit := 150.0
	assert.Equal(t, expectedAvgCredit, report.AvgCredit)

	expectedTransactionsPerMonth := map[string]int{
		"01": 2,
		"02": 1,
	}

	for month, count := range expectedTransactionsPerMonth {
		if report.TransactionsPerMonth[month] != count {
			t.Errorf("expected TransactionsPerMonth[%v] %v, got %v", month, count, report.TransactionsPerMonth[month])
		}
	}
}
