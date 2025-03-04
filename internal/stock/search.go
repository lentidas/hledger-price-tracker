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

package stock

import (
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/text"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lentidas/hledger-price-tracker/internal"
	"github.com/lentidas/hledger-price-tracker/internal/flags"
)

const apiFunctionSearch = "SYMBOL_SEARCH"

type SearchResponse struct {
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

func (response *SearchResponse) DecodeBody(body []byte) error {
	return internal.DecodeBody(body, response)
}

type SearchResponseTyped struct {
	BestMatches []struct {
		Symbol      string
		Name        string
		Type        string
		Region      string
		MarketOpen  time.Time
		MarketClose time.Time
		Timezone    string
		Currency    string
		MatchScore  float32
	}
}

func (responseTyped *SearchResponseTyped) TypeBody(response SearchResponse) error {
	for _, result := range response.BestMatches {
		score, err := strconv.ParseFloat(result.MatchScore, 32)
		if err != nil {
			return fmt.Errorf("[stock.TypeBody] failure to parse match score: %v", err)
		}

		// Parse the timezone offset
		offset, err := strconv.Atoi(result.Timezone[3:])
		if err != nil {
			return fmt.Errorf("[stock.TypeBody] failure to parse timezone offset: %v", err)
		}

		// Create a fixed timezone
		tz := time.FixedZone(result.Timezone, offset*3600)

		openTime, err := time.ParseInLocation("15:04", result.MarketOpen, tz)
		if err != nil {
			return fmt.Errorf("[stock.TypeBody] failure to parse market open time: %v", err)
		}

		closeTime, err := time.ParseInLocation("15:04", result.MarketClose, tz)
		if err != nil {
			return fmt.Errorf("[stock.TypeBody] failure to parse market close time: %v", err)
		}

		responseTyped.BestMatches = append(responseTyped.BestMatches, struct {
			Symbol      string
			Name        string
			Type        string
			Region      string
			MarketOpen  time.Time
			MarketClose time.Time
			Timezone    string
			Currency    string
			MatchScore  float32
		}{
			Symbol:      result.Symbol,
			Name:        result.Name,
			Type:        result.Type,
			Region:      result.Region,
			MarketOpen:  openTime,
			MarketClose: closeTime,
			Timezone:    result.Timezone,
			Currency:    result.Currency,
			MatchScore:  float32(score),
		})
	}

	return nil
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
	if format == flags.OutputFormatCSV {
		url.WriteString("&datatype=csv")
	}

	// Perform the HTTP request.
	body, err := internal.HTTPRequest(url.String())
	if err != nil {
		return "", err
	}

	// Return the output in the desired format.
	switch format {
	case flags.OutputFormatHledger:
		return "", errors.New("[internal.stock.Search] hledger output format not supported")
	case flags.OutputFormatJSON, flags.OutputFormatCSV:
		return string(body), nil
	case flags.OutputFormatTable, flags.OutputFormatTableLong:
		// Parse the JSON body.
		var parsedBody SearchResponse
		err := parsedBody.DecodeBody(body)
		if err != nil {
			return "", err
		}

		// Cast the attributes into proper types.
		var typedBody SearchResponseTyped
		err = typedBody.TypeBody(parsedBody)
		if err != nil {
			return "", err
		}

		t := table.NewWriter()
		if format == flags.OutputFormatTableLong {
			t.AppendHeader(table.Row{"#", "Symbol", "Name", "Type", "Region", "Market Open", "Market Close", "Timezone", "Currency", "Match Score"})
			for i, result := range typedBody.BestMatches {
				t.AppendRow([]interface{}{i + 1, result.Symbol, result.Name, result.Type, result.Region, result.MarketOpen.Format("15:04"), result.MarketClose.Format("15:04"), result.Timezone, result.Currency, fmt.Sprintf("%.2f%%", result.MatchScore*100)})
			}
			t.SetColumnConfigs([]table.ColumnConfig{
				{Number: 6, Align: text.AlignRight},
				{Number: 7, Align: text.AlignRight},
				{Number: 10, Align: text.AlignRight},
			})
		} else {
			t.AppendHeader(table.Row{"#", "Symbol", "Name", "Type", "Region", "Currency", "Match Score"})
			for i, result := range typedBody.BestMatches {
				t.AppendRow([]interface{}{i + 1, result.Symbol, result.Name, result.Type, result.Region, result.Currency, fmt.Sprintf("%.2f%%", result.MatchScore*100)})
			}
			t.SetColumnConfigs([]table.ColumnConfig{
				{Number: 7, Align: text.AlignRight},
			})
		}

		return t.Render(), nil
	default:
		return "", errors.New("[internal.stock.Search] invalid output format")
	}
}
