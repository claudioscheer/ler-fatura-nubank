package main

import (
	"fmt"
	"log"

	"github.com/claudioscheer/ler-faturna-nubank"
)

func main() {
	fatura, err := fatura_nubank.ReadFatura("example/demo.pdf")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Fatura: %v\n", fatura)
}
