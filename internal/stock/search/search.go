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

package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/lentidas/hledger-price-tracker/internal"
	"github.com/lentidas/hledger-price-tracker/internal/flags"
)

const apiFunctionSearch = "SYMBOL_SEARCH"

type Response interface {
	TypeBody() error
	GenerateOutput(body []byte, format flags.OutputFormat) (string, error)
}

type Raw struct {
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

type Typed struct {
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

type Search struct {
	Raw   Raw
	Typed Typed
}

func (obj *Search) TypeBody() error {
	for _, result := range obj.Raw.BestMatches {
		score, err := strconv.ParseFloat(result.MatchScore, 32)
		if err != nil {
			return fmt.Errorf("[stock.search.TypeBody] failure to parse match score: %w", err)
		}

		// Parse the timezone offset.
		offset, err := strconv.Atoi(result.Timezone[3:])
		if err != nil {
			return fmt.Errorf("[stock.search.TypeBody] failure to parse timezone offset: %w", err)
		}

		// Create a fixed timezone.
		tz := time.FixedZone(result.Timezone, offset*3600)

		openTime, err := time.ParseInLocation("15:04", result.MarketOpen, tz)
		if err != nil {
			return fmt.Errorf("[stock.search.TypeBody] failure to parse market open time: %w", err)
		}

		closeTime, err := time.ParseInLocation("15:04", result.MarketClose, tz)
		if err != nil {
			return fmt.Errorf("[stock.search.TypeBody] failure to parse market close time: %w", err)
		}

		obj.Typed.BestMatches = append(obj.Typed.BestMatches, struct {
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

func (obj *Search) GenerateOutput(body []byte, format flags.OutputFormat) (string, error) {
	switch format {
	case flags.OutputFormatHledger:
		return "", errors.New("[internal.search.generateSearchOutput] hledger output format not supported")
	case flags.OutputFormatJSON, flags.OutputFormatCSV:
		return string(body), nil
	case flags.OutputFormatTable, flags.OutputFormatTableLong:
		// Create a struct to parse the JSON body.
		var obj Search

		// Parse the JSON body.
		err := json.Unmarshal(body, &obj.Raw)
		if err != nil {
			return "", fmt.Errorf("[(*Search).GenerateOutput] failure to unmarshal JSON body: %w", err)
		}

		// Cast the attributes into proper types.
		err = obj.TypeBody()
		if err != nil {
			return "", fmt.Errorf("[(*Search).GenerateOutput] error casting response attributes: %w", err)
		}

		t := table.NewWriter()
		t.SetStyle(table.StyleLight)
		if format == flags.OutputFormatTableLong {
			t.AppendHeader(table.Row{"#", "Symbol", "Name", "Type", "Region", "Market Open", "Market Close", "Timezone", "Currency", "Match Score"})
			for i, result := range obj.Typed.BestMatches {
				t.AppendRow([]interface{}{
					i + 1,
					result.Symbol,
					result.Name,
					result.Type,
					result.Region,
					result.MarketOpen.Format("15:04"),
					result.MarketClose.Format("15:04"),
					result.Timezone,
					result.Currency,
					fmt.Sprintf("%.2f%%", result.MatchScore*100),
				})
			}
			t.SetColumnConfigs([]table.ColumnConfig{
				{Number: 6, Align: text.AlignRight},
				{Number: 7, Align: text.AlignRight},
				{Number: 10, Align: text.AlignRight},
			})
		} else {
			t.AppendHeader(table.Row{"#", "Symbol", "Name", "Type", "Region", "Currency", "Match Score"})
			for i, result := range obj.Typed.BestMatches {
				t.AppendRow([]interface{}{
					i + 1,
					result.Symbol,
					result.Name,
					result.Type,
					result.Region,
					result.Currency,
					fmt.Sprintf("%.2f%%", result.MatchScore*100),
				})
			}
			t.SetColumnConfigs([]table.ColumnConfig{
				{Number: 7, Align: text.AlignRight},
			})
		}

		return t.Render(), nil
	default:
		return "", errors.New("[internal.search.generateSearchOutput] invalid output format")
	}
}

// buildURL creates the URL to make the HTTP request to the Alpha Vantage API.
func buildURL(query string, format flags.OutputFormat) (string, error) {
	if internal.ApiKey == "" {
		return "", errors.New("[search.buildURL] api key is required")
	}
	if query == "" {
		return "", errors.New("[search.buildURL] no search query provided")
	}
	switch format {
	case flags.OutputFormatHledger, flags.OutputFormatTable, flags.OutputFormatTableLong, flags.OutputFormatJSON, flags.OutputFormatCSV:
		// Do nothing.
	default:
		return "", errors.New("[search.buildURL] invalid output format")
	}

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

	return url.String(), nil
}

// GetCurrency performs a search for a certain stock symbol, then returns the currency of the first result.
func GetCurrency(symbol string) (string, error) {
	// Alpha Vantage does not use the same stocks for the search and price demos. As such, since we are essentially
	// also performing a search when getting the price of a stock, this would fail with the `demo` API key.
	if internal.DebugMode || internal.ApiKey == "demo" {
		return "NIL", nil
	}

	body, err := Execute(symbol, flags.OutputFormatJSON)
	if err != nil {
		return "", err
	}
	var response Raw
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return "", fmt.Errorf("[stock.search.GetCurrency] error unmarshalling JSON to get currency: %w", err)
	} else if len(response.BestMatches) == 0 {
		return "", fmt.Errorf("[stock.search.GetCurrency] error getting currency of symbol %s", symbol)
	}

	if len(response.BestMatches) < 1 {
		return "", errors.New("[stock.search.GetCurrency] no results found")
	}

	return response.BestMatches[0].Currency, nil
}

// Execute is the core function of the search package. It performs a search for a certain stock symbol and returns the
// results in the specified format.
func Execute(query string, format flags.OutputFormat) (string, error) {
	if internal.ApiKey == "" {
		return "", errors.New("[stock.search.Execute] API key is required")
	}
	if query == "" {
		return "", errors.New("[stock.search.Execute] no search query provided")
	}

	url, err := buildURL(query, format)
	if err != nil {
		return "", err
	}

	body, err := internal.HTTPRequest(url)
	if err != nil {
		return "", err
	}

	response := Search{}

	return response.GenerateOutput(body, format)
}
