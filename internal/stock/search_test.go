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

func TestStockSearchExceptionNoApiKey(t *testing.T) {
	internal.ApiKey = ""

	if _, err := Search("tesco", flags.OutputFormatJSON); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestStockSearchExceptionNoQuery(t *testing.T) {
	internal.ApiKey = "demo"

	if _, err := Search("", flags.OutputFormatJSON); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestStockSearchExceptionInvalidOutputFlag(t *testing.T) {
	internal.ApiKey = "demo"

	if _, err := Search("tesco", "invalid"); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestStockSearchURLBuilder(t *testing.T) {
	internal.ApiKey = "demo"

	expected := "https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=tesco&apikey=demo"
	url, err := buildSearchURL("tesco", flags.OutputFormatJSON)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	if url != expected {
		t.Errorf("expected %s, got %s", expected, url)
	}
}

func TestStockSearchURLBuilderCSV(t *testing.T) {
	internal.ApiKey = "demo"

	expected := "https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=tesco&apikey=demo&datatype=csv"
	url, err := buildSearchURL("tesco", flags.OutputFormatCSV)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if url != expected {
		t.Errorf("expected %s, got %s", expected, url)
	}
}
