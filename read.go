package fatura_nubank

import (
	"strings"
	"time"

	"github.com/gen2brain/go-fitz"
)

type Transaction struct {
	Date        time.Time
	DateString  string
	Description string
	Value       float64
}

type Fatura struct {
	Transactions       []Transaction
	DueDate            time.Time
	PreviousMonthTotal float64
	Total              float64
}

func ReadFatura(pdfFilePath string) (Fatura, error) {
	fatura := Fatura{}

	pdfFile, err := fitz.New(pdfFilePath)
	if err != nil {
		return fatura, err
	}
	defer pdfFile.Close()

	numPages := pdfFile.NumPage()
	for n := 0; n < numPages; n++ {
		text, err := pdfFile.Text(n)
		if err != nil {
			panic(err)
		}

		rows := strings.Split(text, "\n")
		for i := 0; i < len(rows); i++ {
			row := rows[i]
			if strings.Contains(strings.TrimSpace(row), PrefixPreviousMonthTotal) {
				previousMonthTotal, err := extractPreviousMonthTotal(rows[i : i+2])
				if err != nil {
					return fatura, err
				}
				fatura.PreviousMonthTotal = previousMonthTotal

				i++
			}

			if strings.Contains(strings.TrimSpace(row), PrefixTotal) {
				total, err := extractTotal(rows[i : i+2])
				if err != nil {
					return fatura, err
				}
				fatura.Total = total

				i++
			}

			if strings.Contains(strings.TrimSpace(row), PrefixDueDate) {
				dueDate, err := extractDueDate(row)
				if err != nil {
					return fatura, err
				}
				fatura.DueDate = dueDate
			}

			if strings.Contains(strings.TrimSpace(row), PrefixTransactions) {
				transactions, err := extractTransactions(rows[i+1:])
				if err != nil {
					return fatura, err
				}
				fatura.Transactions = transactions

				break
			}
		}
	}

	return fatura, nil
}
