package fatura_nubank

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	PrefixPreviousMonthTotal = "Fatura anterior"
	PrefixTotal              = "Total a pagar"
	PrefixDueDate            = "Data do vencimento"
	PrefixTransactions       = "TRANSAÇÕES"
)

var monthMap = map[string]int{
	"jan": 1,
	"fev": 2,
	"mar": 3,
	"abr": 4,
	"mai": 5,
	"jun": 6,
	"jul": 7,
	"ago": 8,
	"set": 9,
	"out": 10,
	"nov": 11,
	"dez": 12,
}

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

	dateParts := strings.Fields(strings.TrimSpace(values[1]))
	if len(dateParts) != 3 {
		return time.Time{}, errors.New("invalid date format")
	}

	dayStr := dateParts[0]
	monthAbbrev := strings.ToLower(dateParts[1])
	yearStr := dateParts[2]

	day, err := strconv.Atoi(dayStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid day in date string: %s", values[1])
	}

	month, ok := monthMap[monthAbbrev]
	if !ok {
		return time.Time{}, fmt.Errorf("invalid month abbreviation: %s", monthAbbrev)
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid year in date string: %s", values[1])
	}

	dueDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return dueDate, nil
}

func extractTransactions(rows []string, dueDate time.Time) ([]Transaction, error) {
	transactions := []Transaction{}

	startingTransactionPattern := `^(0[1-9]|[12][0-9]|3[01]) [A-Z]{3}$`
	reStartingTransaction := regexp.MustCompile(startingTransactionPattern)

	for i := 0; i < len(rows); i++ {
		row := rows[i]
		if !reStartingTransaction.MatchString(row) {
			continue
		}

		rowsToParse := rows[i : i+2]
		fmt.Println("rowsToParse", rowsToParse)
		rowTransactionValue := findRowWithTransactionValue(rows[i+2 : i+7])
		fmt.Println("rowTransactionValue", rowTransactionValue)
		rowsToParse = append(rowsToParse, rows[rowTransactionValue+i+2])
		fmt.Println("rowsToParse", rowsToParse)

		transaction, err := extractTransaction(rowsToParse)
		if err != nil {
			return nil, err
		}
		if transaction == (Transaction{}) {
			continue
		}

		transactions = append(transactions, transaction)

		i += 2 + rowTransactionValue
	}

	parsedTransactions, err := parseTransactionsDate(transactions, dueDate)
	if err != nil {
		return nil, err
	}

	return parsedTransactions, nil
}

func findRowWithTransactionValue(rows []string) int {
	for i, row := range rows {
		_, err := parseCurrency(row)
		if err == nil {
			return i
		}
	}
	return -1
}

func extractTransaction(rows []string) (Transaction, error) {
	transaction := Transaction{}

	description := strings.TrimSpace(rows[1])
	if len(description) == 0 || strings.HasPrefix(description, "Pagamento em") {
		return transaction, nil
	}

	transaction.DateString = strings.TrimSpace(rows[0])
	transaction.Description = description
	value, err := parseCurrency(rows[2])
	if err != nil {
		return transaction, err
	}
	transaction.Value = value

	return transaction, nil
}

func parseTransactionsDate(transactions []Transaction, dueDate time.Time) ([]Transaction, error) {
	for i, t := range transactions {
		dateParts := strings.Fields(t.DateString)
		if len(dateParts) != 2 {
			return nil, errors.New("invalid date format")
		}

		dayStr := dateParts[0]
		monthAbbrev := strings.ToLower(dateParts[1])

		day, err := strconv.Atoi(dayStr)
		if err != nil {
			return nil, fmt.Errorf("invalid day in date string: %s", t.DateString)
		}

		month, ok := monthMap[monthAbbrev]
		if !ok {
			return nil, fmt.Errorf("invalid month abbreviation: %s", monthAbbrev)
		}

		var year int
		if dueDate.Month() == time.Month(month) {
			year = dueDate.Year()
		} else if dueDate.Month() > time.Month(month) {
			year = dueDate.Year()
		} else {
			year = dueDate.Year() - 1
		}

		t.Date = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		transactions[i] = t
	}

	return transactions, nil
}
