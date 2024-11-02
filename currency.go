package fatura_nubank

import (
	"strconv"
	"strings"
)

func parseCurrency(value string) (float64, error) {
	if strings.HasPrefix(value, "BRL") {
		value = strings.ReplaceAll(value, "BRL", "")
		value = strings.Split(value, " =")[0]
	}

	value = strings.TrimSpace(strings.ReplaceAll(value, "R$", ""))
	value = strings.ReplaceAll(value, ".", "")
	value = strings.ReplaceAll(value, ",", ".")

	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}

	return result, nil
}
