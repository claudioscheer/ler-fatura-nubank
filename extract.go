package fatura_nubank

import (
	"regexp"
	"strings"
	"time"
)

const (
	PrefixPreviousMonthTotal = "Fatura anterior"
	PrefixTotal              = "Total a pagar"
	PrefixDueDate            = "Data do vencimento"
	PrefixTransactions       = "TRANSAÇÕES"
)

func extractPreviousMonthTotal(rows []string) (float64, error) {
	previousMonthTotal, err := parseCurrency(rows[1])
	if err != nil {
		return 0, err
	}
	return previousMonthTotal, nil
}

func extractTotal(rows []string) (float64, error) {
	total, err := parseCurrency(rows[1])
	if err != nil {
		return 0, err
	}
	return total, nil
}

func extractDueDate(row string) (time.Time, error) {
	values := strings.Split(row, ":")
	dueDate, err := time.Parse("02 Jan 2006", strings.TrimSpace(values[1]))
	if err != nil {
		return time.Time{}, err
	}
	return dueDate, nil
}

func extractTransactions(rows []string) ([]Transaction, error) {
	transactions := []Transaction{}

	startingTransactionPattern := `^\d{1,2} [A-Z]{3}$`
	reStartingTransaction := regexp.MustCompile(startingTransactionPattern)

	for i := 0; i < len(rows); i++ {
		row := rows[i]
		if !reStartingTransaction.MatchString(row) {
			continue
		}

		transaction, err := extractTransaction(rows[i : i+3])
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func extractTransaction(rows []string) (Transaction, error) {
	transaction := Transaction{}

  // TODO: parse date
  // consider starting month can be dez and ending month can be jan
	transaction.DateString = strings.TrimSpace(rows[0])
	transaction.Description = strings.TrimSpace(rows[1])
	value, err := parseCurrency(rows[2])
	if err != nil {
		return transaction, err
	}
	transaction.Value = value

	return transaction, nil
}
