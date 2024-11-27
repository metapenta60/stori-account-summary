package model

import "fmt"

type Rows []Row

type Row struct {
	Id          string
	Date        Date
	Transaction float64
}

type Date struct {
	Day   string
	Month string
}

type AccountReport struct {
	Sum                  float64
	TransactionsPerMonth map[string]int
	AvgDebit             float64
	AvgCredit            float64

	NumCreditTransaction float64
	NumDebitTransaction  float64

	TotalCredit float64
	TotalDebit  float64
}

func (ar *AccountReport) IncreaseTransactionCount(month string) {
	ar.TransactionsPerMonth[month]++
}

func (ar *AccountReport) AddTransaction(transaction float64) {
	ar.Sum += transaction
	if transaction > 0 {
		ar.TotalCredit += transaction
		ar.NumCreditTransaction++
	} else {
		ar.TotalDebit += transaction
		ar.NumDebitTransaction++
	}
}

func (ar *AccountReport) UpdateAverageDebitAndCredit() {
	ar.AvgDebit = calculateAverage(ar.TotalDebit, ar.NumDebitTransaction)
	ar.AvgCredit = calculateAverage(ar.TotalCredit, ar.NumCreditTransaction)
}

func calculateAverage(total float64, numTransactions float64) float64 {
	if numTransactions > 0 {
		return total / numTransactions
	}

	return 0
}

var mapToMonth = map[string]string{
	"1":  "January",
	"2":  "February",
	"3":  "March",
	"4":  "April",
	"5":  "May",
	"6":  "June",
	"7":  "July",
	"8":  "August",
	"9":  "September",
	"10": "October",
	"11": "November",
	"12": "December",
}

func (ar AccountReport) TransactionPerMonth() []string {
	transactions := make([]string, 0)
	for i := 1; i <= 12; i++ {
		k := fmt.Sprintf("%d", i)
		v := ar.TransactionsPerMonth[k]

		if v > 0 {
			month := mapToMonth[k]
			count := v

			transactionsMessage := fmt.Sprintf("Number of transactions in %s: %d", month, count)
			transactions = append(transactions, transactionsMessage)
		}
	}

	return transactions
}
