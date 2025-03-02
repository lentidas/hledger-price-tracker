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
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <https://www.gnu.org/licenses/>.
 */

package stock

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lentidas/hledger-price-tracker/internal"
	"github.com/lentidas/hledger-price-tracker/internal/flags"
)

const apiFunctionSearch = "SYMBOL_SEARCH"

type SearchBody struct {
	BestMatches []struct {
		Symbol      string `json:"1. symbol"`
		Name        string `json:"2. name"`
		Type        string `json:"3. type"`
		Region      string `json:"4. region"`
		MarketOpen  string `json:"5. marketOpen"`
		MarketClose string `json:"6. marketClose"`
		Timezone    string `json:"7. timezone"`
		Currency    string `json:"8. currency"`
		MatchScore  string `json:"9. matchScore"`
	} `json:"bestMatches"`
}

func parseBody(body []byte) (SearchBody, error) {
	var parsedBody SearchBody
	err := json.Unmarshal(body, &parsedBody)
	if err != nil {
		return SearchBody{}, fmt.Errorf("[internal.stock.parseBody] failure to unmarshal JSON body: %v", err)
	}
	return parsedBody, nil
}

func Search(query string, format flags.OutputFormat) (string, error) {
	// Verify function parameters and variables.
	if internal.ApiKey == "" {
		return "", errors.New("[stock.Search] api key is required")
	}
	if query == "" {
		return "", errors.New("[stock.Search] no search query provided")
	}

	// Create the request URL.
	url := strings.Builder{}
	url.WriteString(internal.ApiBaseUrl)
	url.WriteString("function=")
	url.WriteString(apiFunctionSearch)
	url.WriteString("&keywords=")
	url.WriteString(query)
	url.WriteString("&apikey=")
	url.WriteString(internal.ApiKey)
	if format == flags.OutputFormatCsv {
		url.WriteString("&datatype=csv")
	}

	// Perform the HTTP request.
	body, err := internal.HttpRequest(url.String())
	if err != nil {
		return "", err
	}

	// Return the output in the desired format.
	switch format {
	case flags.OutputFormatJson, flags.OutputFormatCsv:
		return string(body), nil
	case flags.OutputFormatTable, flags.OutputFormatTableLong:
		// Parse the JSON body.
		parsedBody, err := parseBody(body)
		if err != nil {
			return "", err
		}

		t := table.NewWriter()
		if format == flags.OutputFormatTableLong {
			t.AppendHeader(table.Row{"#", "Symbol", "Name", "Type", "Region", "Market Open", "Market Close", "Timezone", "Currency", "Match Score"})
			for i, result := range parsedBody.BestMatches {
				t.AppendRow([]interface{}{i + 1, result.Symbol, result.Name, result.Type, result.Region, result.MarketOpen, result.MarketClose, result.Timezone, result.Currency, result.MatchScore})
			}
		} else {
			t.AppendHeader(table.Row{"#", "Symbol", "Name", "Type", "Region", "Currency", "Match Score"})
			for i, result := range parsedBody.BestMatches {
				t.AppendRow([]interface{}{i + 1, result.Symbol, result.Name, result.Type, result.Region, result.Currency, result.MatchScore})
			}
		}

		return t.Render(), nil
	default:
		return "", errors.New("[internal.stock.Search] invalid output format")
	}
}
