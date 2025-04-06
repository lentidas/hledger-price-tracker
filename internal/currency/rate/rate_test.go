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

package rate

import (
	"testing"

	"github.com/lentidas/hledger-price-tracker/internal"
	"github.com/lentidas/hledger-price-tracker/internal/flags"
)

func TestRate(t *testing.T) {
	internal.ApiKey = "demo"

	t.Run("success from EUR to USD daily", func(t *testing.T) {
		if _, err := Execute("EUR", "USD", flags.OutputFormatHledger, flags.IntervalDaily, "", "", false); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("success from EUR to USD daily full", func(t *testing.T) {
		if _, err := Execute("EUR", "USD", flags.OutputFormatHledger, flags.IntervalDaily, "", "", true); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("success from EUR to USD weekly", func(t *testing.T) {
		if _, err := Execute("EUR", "USD", flags.OutputFormatHledger, flags.IntervalWeekly, "", "", false); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("success from EUR to USD monthly", func(t *testing.T) {
		if _, err := Execute("EUR", "USD", flags.OutputFormatHledger, flags.IntervalMonthly, "", "", false); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("no origin currency", func(t *testing.T) {
		if _, err := Execute("", "USD", flags.OutputFormatHledger, flags.IntervalDaily, "", "", false); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("no destination currency", func(t *testing.T) {
		if _, err := Execute("EUR", "", flags.OutputFormatHledger, flags.IntervalDaily, "", "", false); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("invalid origin currency", func(t *testing.T) {
		if _, err := Execute("INVALID", "USD", flags.OutputFormatHledger, flags.IntervalDaily, "", "", false); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("invalid destination currency", func(t *testing.T) {
		if _, err := Execute("EUR", "INVALID", flags.OutputFormatHledger, flags.IntervalDaily, "", "", false); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("invalid output format", func(t *testing.T) {
		if _, err := Execute("EUR", "USD", "invalid", flags.IntervalDaily, "", "", false); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("invalid interval", func(t *testing.T) {
		if _, err := Execute("EUR", "USD", flags.OutputFormatHledger, "invalid", "", "", false); err == nil {
			t.Error("expected error, got nil")
		}
	})

	internal.ApiKey = ""

	t.Run("no API key", func(t *testing.T) {
		if _, err := Execute("USD", "JPY", flags.OutputFormatHledger, flags.IntervalDaily, "", "", false); err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestRateURLBuilder(t *testing.T) {
	internal.ApiKey = "demo"

	t.Run("daily", func(t *testing.T) {
		expected := "https://www.alphavantage.co/query?function=FX_DAILY&from_symbol=EUR&to_symbol=USD&apikey=demo"

		url, err := buildURL("EUR", "USD", flags.OutputFormatHledger, flags.IntervalDaily, false)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		} else if url != expected {
			t.Errorf("expected %s, got %s", expected, url)
		}
	})

	t.Run("daily full", func(t *testing.T) {
		expected := "https://www.alphavantage.co/query?function=FX_DAILY&from_symbol=EUR&to_symbol=USD&outputsize=full&apikey=demo"

		url, err := buildURL("EUR", "USD", flags.OutputFormatHledger, flags.IntervalDaily, true)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		} else if url != expected {
			t.Errorf("expected %s, got %s", expected, url)
		}
	})

	t.Run("weekly", func(t *testing.T) {
		expected := "https://www.alphavantage.co/query?function=FX_WEEKLY&from_symbol=EUR&to_symbol=USD&apikey=demo"

		url, err := buildURL("EUR", "USD", flags.OutputFormatHledger, flags.IntervalWeekly, false)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		} else if url != expected {
			t.Errorf("expected %s, got %s", expected, url)
		}
	})

	t.Run("weekly ignore full", func(t *testing.T) {
		expected := "https://www.alphavantage.co/query?function=FX_WEEKLY&from_symbol=EUR&to_symbol=USD&apikey=demo"

		url, err := buildURL("EUR", "USD", flags.OutputFormatHledger, flags.IntervalWeekly, true)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		} else if url != expected {
			t.Errorf("expected %s, got %s", expected, url)
		}
	})

	t.Run("monthly", func(t *testing.T) {
		expected := "https://www.alphavantage.co/query?function=FX_MONTHLY&from_symbol=EUR&to_symbol=USD&apikey=demo"

		url, err := buildURL("EUR", "USD", flags.OutputFormatHledger, flags.IntervalMonthly, false)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		} else if url != expected {
			t.Errorf("expected %s, got %s", expected, url)
		}
	})

	t.Run("monthly ignore full", func(t *testing.T) {
		expected := "https://www.alphavantage.co/query?function=FX_MONTHLY&from_symbol=EUR&to_symbol=USD&apikey=demo"

		url, err := buildURL("EUR", "USD", flags.OutputFormatHledger, flags.IntervalMonthly, true)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		} else if url != expected {
			t.Errorf("expected %s, got %s", expected, url)
		}
	})

	t.Run("CSV", func(t *testing.T) {
		expected := "https://www.alphavantage.co/query?function=FX_DAILY&from_symbol=EUR&to_symbol=USD&outputsize=full&apikey=demo&datatype=csv"

		url, err := buildURL("EUR", "USD", flags.OutputFormatCSV, flags.IntervalDaily, true)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		} else if url != expected {
			t.Errorf("expected %s, got %s", expected, url)
		}
	})
}
