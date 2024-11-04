# Fatura Nubank

## Overview

The `fatura_nubank` package provides functionality to read and parse PDF files containing transaction data from Nubank statements. It extracts relevant information such as transaction dates, descriptions, values, and overall summary details of the statement.

## Installation

To use the `fatura_nubank` package in your Go project, follow these steps:

1. **Create a new Go module** (if you haven't already):

   ```bash
   mkdir myproject
   cd myproject
   go mod init myproject
   ```

2. **Install the package**: Run the following command:

   ```bash
   go get github.com/claudioscheer/ler-fatura-nubank
   ```

## Usage

Here's how to use the ReadFatura function to read a Nubank PDF statement:

```go
package main

import (
    "fmt"
    "log"
    "github.com/claudioscheer/ler-fatura-nubank"
)

func main() {
    pdfFilePath := "path/to/your/nubank_statement.pdf"

    fatura, err := fatura_nubank.ReadFatura(pdfFilePath)
    if err != nil {
    	log.Fatalf("Error reading fatura: %v", err)
    }

    fmt.Printf("Due Date: %s\n", fatura.DueDate.Format("2006-01-02"))
    fmt.Printf("Total Amount: %.2f\n", fatura.Total)
    fmt.Printf("Previous Month Total: %.2f\n", fatura.PreviousMonthTotal)

    for _, transaction := range fatura.Transactions {
    	fmt.Printf("Date: %s, Description: %s, Value: %.2f\n",
    		transaction.Date.Format("2006-01-02"), transaction.Description, transaction.Value)
    }
}
```

## Function Signature

```go
func ReadFatura(pdfFilePath string) (Fatura, error)
```

### Parameters

- `pdfFilePath`: A string representing the path to the PDF file containing the Nubank statement.

### Returns

- `Fatura`: A struct containing all transactions and summary information.
- `error`: An error object if something goes wrong during the reading process.

## Structs

### `Transaction`

The `Transaction` struct represents an individual transaction with the following fields:

- `Date`: The date of the transaction.
- `DateString`: A string representation of the date.
- `Description`: A description of the transaction.
- `Value`: The monetary value of the transaction.

### `Fatura`

The Fatura struct contains information about the entire statement:

- `Transactions`: A slice of Transaction structs representing all transactions in the statement.
- `DueDate`: The due date for payment.
- `PreviousMonthTotal`: The total amount from the previous month.
- `Total`: The total amount for the current statement.
