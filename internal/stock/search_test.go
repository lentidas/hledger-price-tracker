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
		if _, err := Search("", flags.OutputFormatJSON); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("invalid output format", func(t *testing.T) {
		if _, err := Search("tesco", "invalid"); err == nil {
			t.Error("expected error, got nil")
		}
	})

	internal.ApiKey = ""

	t.Run("no API key", func(t *testing.T) {
		if _, err := Search("tesco", flags.OutputFormatJSON); err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestSearchURLBuilder(t *testing.T) {
	internal.ApiKey = "demo"

	t.Run("normal", func(t *testing.T) {
		expected := "https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=tesco&apikey=demo"
		url, err := buildSearchURL("tesco", flags.OutputFormatJSON)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if url != expected {
			t.Errorf("expected %s, got %s", expected, url)
		}
	})

	t.Run("CSV", func(t *testing.T) {
		expected := "https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=tesco&apikey=demo&datatype=csv"
		url, err := buildSearchURL("tesco", flags.OutputFormatCSV)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if url != expected {
			t.Errorf("expected %s, got %s", expected, url)
		}
	})
}
