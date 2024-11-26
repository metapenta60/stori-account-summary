package model

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
	NumDebutTransaction  float64

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
		ar.NumDebutTransaction++
	}
}

func (ar *AccountReport) UpdateAverageDebitAndCredit() {
	ar.AvgDebit = calculateAverage(ar.TotalDebit, ar.NumDebutTransaction)
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

func (ar AccountReport) TransactionPerMonth() []Transaction {
	transactions := make([]Transaction, 0)
	for k, v := range ar.TransactionsPerMonth {
		transactions = append(transactions, Transaction{
			Month: mapToMonth[k],
			Count: v,
		})
	}

	return transactions
}
