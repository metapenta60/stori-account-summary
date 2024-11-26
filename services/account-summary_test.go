package services

import (
	"stori-account-summary/model"
	"testing"
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

	// Crear una instancia del servicio
	ass := accountSummaryService{}

	// Llamar a la función a probar
	report := ass.analyseAccount(rows)

	// Verificar los resultados esperados
	expectedSum := 250.0
	if report.Sum != expectedSum {
		t.Errorf("expected Sum %v, got %v", expectedSum, report.Sum)
	}

	expectedTotalCredit := 300.0
	if report.TotalCredit != expectedTotalCredit {
		t.Errorf("expected TotalCredit %v, got %v", expectedTotalCredit, report.TotalCredit)
	}

	expectedTotalDebit := -50.0
	if report.TotalDebit != expectedTotalDebit {
		t.Errorf("expected TotalDebit %v, got %v", expectedTotalDebit, report.TotalDebit)
	}

	expectedNumCreditTransaction := 2.0
	if report.NumCreditTransaction != expectedNumCreditTransaction {
		t.Errorf("expected NumCreditTransaction %v, got %v", expectedNumCreditTransaction, report.NumCreditTransaction)
	}

	expectedNumDebitTransaction := 1.0
	if report.NumDebutTransaction != expectedNumDebitTransaction {
		t.Errorf("expected NumDebutTransaction %v, got %v", expectedNumDebitTransaction, report.NumDebutTransaction)
	}

	expectedAvgDebit := -50.0
	if report.AvgDebit != expectedAvgDebit {
		t.Errorf("expected AvgDebit %v, got %v", expectedAvgDebit, report.AvgDebit)
	}

	expectedAvgCredit := 150.0
	if report.AvgCredit != expectedAvgCredit {
		t.Errorf("expected AvgCredit %v, got %v", expectedAvgCredit, report.AvgCredit)
	}

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
