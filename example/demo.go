package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/claudioscheer/ler-fatura-nubank"
)

func main() {
	fatura, err := fatura_nubank.ReadFatura("example/demo-1.pdf")
	if err != nil {
		log.Fatal(err)
	}

	value, _ := json.MarshalIndent(fatura, "", "  ")
	fmt.Printf("Fatura: %v\n", string(value))
}
