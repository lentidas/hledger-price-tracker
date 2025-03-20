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

const url = "https://www.alphavantage.co/physical_currency_list/"

type Currency struct {
	Code string
	Name string
}

type Currencies []Currency

func (obj *Currencies) GenerateOutput(body []byte, format flags.OutputFormat) (string, error) {
	switch format {
	case flags.OutputFormatHledger, flags.OutputFormatJSON, flags.OutputFormatTableLong:
		errorMessage := fmt.Sprintf("[(*Currencies).GenerateOutput] %s output format not supported", format)
		return "", errors.New(errorMessage)
	case flags.OutputFormatCSV:
		return string(body), nil
	case flags.OutputFormatTable:
		// Read CSV data.
		csvReader := csv.NewReader(bytes.NewReader(body))
		data, err := csvReader.ReadAll()
		if err != nil {
			return "", fmt.Errorf("[(*Currencies).GenerateOutput] failure to read CSV data: %w", err)
		}

		// Parse CSV data into a slice of Currency objects.
		for i, line := range data {
			if i > 0 { // Skip the header line.
				var currency Currency
				for j, field := range line {
					switch j {
					case 0:
						currency.Code = field
					case 1:
						currency.Name = field
					}
				}
				*obj = append(*obj, currency)
			}
		}

		// Create a table and send it to the output.
		t := table.NewWriter()
		t.SetStyle(table.StyleLight)
		t.AppendHeader(table.Row{"Code", "Currency Name"})
		for _, currency := range *obj {
			t.AppendRow(table.Row{currency.Code, currency.Name})
		}
		return t.Render() + "\n", nil
	default:
		return "", errors.New("[internal.search.generateSearchOutput] invalid output format")
	}
}

func Execute(format flags.OutputFormat) (string, error) {
	body, err := internal.HTTPRequest(url)
	if err != nil {
		return "", err
	}

	response := Currencies{}

	return response.GenerateOutput(body, format)
}
