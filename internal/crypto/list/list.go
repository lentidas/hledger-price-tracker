/*
 * hledger-price-tracker - a CLI tool to get market prices for commodities
 * Copyright (C) 2024 Gonçalo Carvalheiro Heleno
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <https://www.gnu.org/licenses/>.
 */

package list

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/lentidas/hledger-price-tracker/internal"
	"github.com/lentidas/hledger-price-tracker/internal/flags"
)

const url = "https://www.alphavantage.co/digital_currency_list/"

type Cryptos map[string]string

func (obj *Cryptos) GenerateOutput(body []byte, format flags.OutputFormat) (string, error) {
	switch format {
	case flags.OutputFormatHledger, flags.OutputFormatJSON, flags.OutputFormatTableLong:
		errorMessage := fmt.Sprintf("[(*Cryptos).GenerateOutput] %s output format not supported", format)
		return "", errors.New(errorMessage)
	case flags.OutputFormatCSV:
		return string(body), nil
	case flags.OutputFormatTable:
		// Read CSV data.
		csvReader := csv.NewReader(bytes.NewReader(body))
		data, err := csvReader.ReadAll()
		if err != nil {
			return "", fmt.Errorf("[(*Cryptos).GenerateOutput] failure to read CSV data: %w", err)
		}

		// Parse CSV data into a slice of Crypto objects.
		for i, line := range data {
			if i > 0 { // Skip the header line.
				var currencyCode string
				var currencyName string
				for j, field := range line {
					switch j {
					case 0:
						currencyCode = field
					case 1:
						currencyName = field
					}
				}
				(*obj)[currencyCode] = currencyName
			}
		}

		// Create a table and send it to the output.
		t := table.NewWriter()
		t.SetStyle(table.StyleLight)
		t.AppendHeader(table.Row{"Code", "Currency Name"})
		for currencyCode, currencyName := range *obj {
			t.AppendRow(table.Row{currencyCode, currencyName})
		}
		t.SortBy([]table.SortBy{{Name: "Code", Mode: table.Asc}})
		return t.Render() + "\n", nil
	default:
		return "", errors.New("[(*Cryptos).GenerateOutput] invalid output format")
	}
}

func Execute(format flags.OutputFormat) (string, error) {
	body, err := internal.HTTPRequest(url)
	if err != nil {
		return "", err
	}

	response := Cryptos{}

	return response.GenerateOutput(body, format)
}
