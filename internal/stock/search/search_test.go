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
	"testing"

	"github.com/lentidas/hledger-price-tracker/internal"
	"github.com/lentidas/hledger-price-tracker/internal/flags"
)

func TestSearch(t *testing.T) {
	internal.ApiKey = "demo"

	t.Run("success", func(t *testing.T) {
		// TODO Implement expected response.
	})

	t.Run("no search query", func(t *testing.T) {
		if _, err := Execute("", flags.OutputFormatJSON); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("invalid output format", func(t *testing.T) {
		if _, err := Execute("tesco", "invalid"); err == nil {
			t.Error("expected error, got nil")
		}
	})

	internal.ApiKey = ""

	t.Run("no API key", func(t *testing.T) {
		if _, err := Execute("tesco", flags.OutputFormatJSON); err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestSearchURLBuilder(t *testing.T) {
	internal.ApiKey = "demo"

	t.Run("normal", func(t *testing.T) {
		expected := "https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=tesco&apikey=demo"
		url, err := buildURL("tesco", flags.OutputFormatJSON)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if url != expected {
			t.Errorf("expected %s, got %s", expected, url)
		}
	})

	t.Run("CSV", func(t *testing.T) {
		expected := "https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=tesco&apikey=demo&datatype=csv"
		url, err := buildURL("tesco", flags.OutputFormatCSV)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
		if url != expected {
			t.Errorf("expected %s, got %s", expected, url)
		}
	})
}

func TestSearchResponseTypeBody(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var response Search
		response.Raw = Raw{
			BestMatches: []struct {
				Symbol      string `json:"1. symbol"`
				Name        string `json:"2. name"`
				Type        string `json:"3. type"`
				Region      string `json:"4. region"`
				MarketOpen  string `json:"5. marketOpen"`
				MarketClose string `json:"6. marketClose"`
				Timezone    string `json:"7. timezone"`
				Currency    string `json:"8. currency"`
				MatchScore  string `json:"9. matchScore"`
			}{
				{
					Symbol:      "AAPL",
					Name:        "Apple Inc.",
					Type:        "Equity",
					Region:      "United States",
					MarketOpen:  "09:30",
					MarketClose: "16:00",
					Timezone:    "UTC-05",
					Currency:    "USD",
					MatchScore:  "0.95",
				},
			},
		}

		err := response.TypeBody()
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if len(response.Typed.BestMatches) != 1 {
			t.Fatalf("expected 1 match, got %d", len(response.Typed.BestMatches))
		}

		match := response.Typed.BestMatches[0]
		if match.Symbol != "AAPL" {
			t.Errorf("expected AAPL, got %s", match.Symbol)
		}
		if match.Name != "Apple Inc." {
			t.Errorf("expected Apple Inc., got %s", match.Name)
		}
		if match.Type != "Equity" {
			t.Errorf("expected Equity, got %s", match.Type)
		}
		if match.Region != "United States" {
			t.Errorf("expected United States, got %s", match.Region)
		}
		if match.Currency != "USD" {
			t.Errorf("expected USD, got %s", match.Currency)
		}
		if match.MatchScore != 0.95 {
			t.Errorf("expected 0.95, got %f", match.MatchScore)
		}
	})

	t.Run("invalid match score", func(t *testing.T) {
		var response Search
		response.Raw = Raw{
			BestMatches: []struct {
				Symbol      string `json:"1. symbol"`
				Name        string `json:"2. name"`
				Type        string `json:"3. type"`
				Region      string `json:"4. region"`
				MarketOpen  string `json:"5. marketOpen"`
				MarketClose string `json:"6. marketClose"`
				Timezone    string `json:"7. timezone"`
				Currency    string `json:"8. currency"`
				MatchScore  string `json:"9. matchScore"`
			}{
				{
					Symbol:      "AAPL",
					Name:        "Apple Inc.",
					Type:        "Equity",
					Region:      "United States",
					MarketOpen:  "09:30",
					MarketClose: "16:00",
					Timezone:    "UTC-05",
					Currency:    "USD",
					MatchScore:  "invalid",
				},
			},
		}

		err := response.TypeBody()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("invalid timezone", func(t *testing.T) {
		var response Search
		response.Raw = Raw{
			BestMatches: []struct {
				Symbol      string `json:"1. symbol"`
				Name        string `json:"2. name"`
				Type        string `json:"3. type"`
				Region      string `json:"4. region"`
				MarketOpen  string `json:"5. marketOpen"`
				MarketClose string `json:"6. marketClose"`
				Timezone    string `json:"7. timezone"`
				Currency    string `json:"8. currency"`
				MatchScore  string `json:"9. matchScore"`
			}{
				{
					Symbol:      "AAPL",
					Name:        "Apple Inc.",
					Type:        "Equity",
					Region:      "United States",
					MarketOpen:  "09:30",
					MarketClose: "16:00",
					Timezone:    "invalid",
					Currency:    "USD",
					MatchScore:  "0.95",
				},
			},
		}

		err := response.TypeBody()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("invalid market open time", func(t *testing.T) {
		var response Search
		response.Raw = Raw{
			BestMatches: []struct {
				Symbol      string `json:"1. symbol"`
				Name        string `json:"2. name"`
				Type        string `json:"3. type"`
				Region      string `json:"4. region"`
				MarketOpen  string `json:"5. marketOpen"`
				MarketClose string `json:"6. marketClose"`
				Timezone    string `json:"7. timezone"`
				Currency    string `json:"8. currency"`
				MatchScore  string `json:"9. matchScore"`
			}{
				{
					Symbol:      "AAPL",
					Name:        "Apple Inc.",
					Type:        "Equity",
					Region:      "United States",
					MarketOpen:  "invalid",
					MarketClose: "16:00",
					Timezone:    "UTC-05",
					Currency:    "USD",
					MatchScore:  "0.95",
				},
			},
		}

		err := response.TypeBody()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("invalid market close time", func(t *testing.T) {
		var response Search
		response.Raw = Raw{
			BestMatches: []struct {
				Symbol      string `json:"1. symbol"`
				Name        string `json:"2. name"`
				Type        string `json:"3. type"`
				Region      string `json:"4. region"`
				MarketOpen  string `json:"5. marketOpen"`
				MarketClose string `json:"6. marketClose"`
				Timezone    string `json:"7. timezone"`
				Currency    string `json:"8. currency"`
				MatchScore  string `json:"9. matchScore"`
			}{
				{
					Symbol:      "AAPL",
					Name:        "Apple Inc.",
					Type:        "Equity",
					Region:      "United States",
					MarketOpen:  "09:30",
					MarketClose: "invalid",
					Timezone:    "UTC-05",
					Currency:    "USD",
					MatchScore:  "0.95",
				},
			},
		}

		err := response.TypeBody()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
