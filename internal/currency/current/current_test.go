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

package current

import (
	"testing"

	"github.com/lentidas/hledger-price-tracker/internal"
	"github.com/lentidas/hledger-price-tracker/internal/flags"
)

func TestCurrent(t *testing.T) {
	internal.ApiKey = "demo"

	// t.Run("success", func(t *testing.T) {
	// 	// TODO Implement expected response.
	// })

	t.Run("no origin currency", func(t *testing.T) {
		if _, err := Execute("", "JPY", flags.OutputFormatHledger); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("no destination currency", func(t *testing.T) {
		if _, err := Execute("USD", "", flags.OutputFormatHledger); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("invalid currency", func(t *testing.T) {
		if _, err := Execute("INVALID", "JPY", flags.OutputFormatHledger); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("invalid output format", func(t *testing.T) {
		if _, err := Execute("USD", "JPY", "csv"); err == nil {
			t.Error("expected error, got nil")
		}
	})

	internal.ApiKey = ""

	t.Run("no API key", func(t *testing.T) {
		if _, err := Execute("USD", "JPY", flags.OutputFormatHledger); err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestCurrentURLBuilder(t *testing.T) {
	internal.ApiKey = "demo"

	t.Run("currency to currency", func(t *testing.T) {
		expected := "https://www.alphavantage.co/query?function=CURRENCY_EXCHANGE_RATE&from_currency=USD&to_currency=JPY&apikey=demo"

		url, err := buildURL("USD", "JPY")
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		} else if url != expected {
			t.Errorf("expected %s, got %s", expected, url)
		}
	})

	t.Run("crypto to currency", func(t *testing.T) {
		expected := "https://www.alphavantage.co/query?function=CURRENCY_EXCHANGE_RATE&from_currency=BTC&to_currency=EUR&apikey=demo"

		url, err := buildURL("BTC", "EUR")
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		} else if url != expected {
			t.Errorf("expected %s, got %s", expected, url)
		}
	})

	t.Run("invalid currency", func(t *testing.T) {
		_, err := buildURL("INVALID", "USD")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
