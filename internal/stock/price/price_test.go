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

package price

import (
	"testing"

	"github.com/lentidas/hledger-price-tracker/internal"
	"github.com/lentidas/hledger-price-tracker/internal/flags"
)

// TODO Add unitary tests for the parsing of the price response, per interval, and adjusted or not.

func TestPrice(t *testing.T) {
	internal.ApiKey = "demo"

	// t.Run("success", func(t *testing.T) {
	// 	// TODO Implement expected response.
	// })

	t.Run("no symbol", func(t *testing.T) {
		if _, err := Execute("", flags.OutputFormatHledger, flags.IntervalWeekly, "", "", false, false); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("invalid output format", func(t *testing.T) {
		if _, err := Execute("tesco", "invalid", flags.IntervalWeekly, "", "", false, false); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("invalid interval", func(t *testing.T) {
		if _, err := Execute("tesco", flags.OutputFormatHledger, "invalid", "", "", false, false); err == nil {
			t.Error("expected error, got nil")
		}
	})

	internal.ApiKey = ""

	t.Run("no API key", func(t *testing.T) {
		if _, err := Execute("tesco", flags.OutputFormatHledger, flags.IntervalWeekly, "", "", false, false); err == nil {
			t.Error("expected error, got nil")
		}
	})
}

// FIXME Fix the tests to use the new functions
// FIXME Use current_test.go as an example because we should use a if else instead of consecutive if statements.

func TestPriceURLBuilder(t *testing.T) {
	internal.ApiKey = "demo"

	t.Run("daily", func(t *testing.T) {
		expected := "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=IBM&outputsize=full&apikey=demo"

		url, err := buildURL("IBM", flags.OutputFormatHledger, flags.IntervalDaily, false, true)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		} else if url != expected {
			t.Errorf("expected %s, got %s", expected, url)
		}
	})

	t.Run("daily adjusted", func(t *testing.T) {
		expected := "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY_ADJUSTED&symbol=IBM&outputsize=full&apikey=demo"

		url, err := buildURL("IBM", flags.OutputFormatHledger, flags.IntervalDaily, true, true)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		} else if url != expected {
			t.Errorf("expected %s, got %s", expected, url)
		}
	})

	t.Run("weekly", func(t *testing.T) {
		expected := "https://www.alphavantage.co/query?function=TIME_SERIES_WEEKLY&symbol=IBM&apikey=demo"

		url, err := buildURL("IBM", flags.OutputFormatHledger, flags.IntervalWeekly, false, false)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		} else if url != expected {
			t.Errorf("expected %s, got %s", expected, url)
		}
	})

	t.Run("weekly adjusted", func(t *testing.T) {
		expected := "https://www.alphavantage.co/query?function=TIME_SERIES_WEEKLY_ADJUSTED&symbol=IBM&apikey=demo"

		url, err := buildURL("IBM", flags.OutputFormatHledger, flags.IntervalWeekly, true, false)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		} else if url != expected {
			t.Errorf("expected %s, got %s", expected, url)
		}
	})

	t.Run("monthly", func(t *testing.T) {
		expected := "https://www.alphavantage.co/query?function=TIME_SERIES_MONTHLY&symbol=IBM&apikey=demo"

		url, err := buildURL("IBM", flags.OutputFormatHledger, flags.IntervalMonthly, false, false)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		} else if url != expected {
			t.Errorf("expected %s, got %s", expected, url)
		}
	})

	t.Run("monthly adjusted", func(t *testing.T) {
		expected := "https://www.alphavantage.co/query?function=TIME_SERIES_MONTHLY_ADJUSTED&symbol=IBM&apikey=demo"

		url, err := buildURL("IBM", flags.OutputFormatHledger, flags.IntervalMonthly, true, false)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		} else if url != expected {
			t.Errorf("expected %s, got %s", expected, url)
		}
	})

	t.Run("CSV", func(t *testing.T) {
		expected := "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=IBM&outputsize=full&apikey=demo&datatype=csv"

		url, err := buildURL("IBM", flags.OutputFormatCSV, flags.IntervalDaily, false, true)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		} else if url != expected {
			t.Errorf("expected %s, got %s", expected, url)
		}
	})
}
