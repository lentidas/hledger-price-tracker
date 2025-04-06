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

package internal

import (
	"reflect"
	"testing"
	"time"
)

func TestParseCurrenciesCSV(t *testing.T) {
	t.Run("valid CSV data", func(t *testing.T) {
		data := []byte("code,name\nUSD,United States Dollar\nEUR,Euro")
		expected := map[string]string{
			"USD": "United States Dollar",
			"EUR": "Euro",
		}

		result, err := ParseCurrenciesCSV(data)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("empty CSV data", func(t *testing.T) {
		data := []byte("code,name\n")
		expected := map[string]string{}

		result, err := ParseCurrenciesCSV(data)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("malformed CSV data", func(t *testing.T) {
		data := []byte("code,name\nUSD,United States Dollar\nEUR")

		_, err := ParseCurrenciesCSV(data)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestValidateDates(t *testing.T) {
	t.Run("begin date in the future", func(t *testing.T) {
		date := time.Now()
		date = date.Add(3 * 24 * time.Hour) // Add 3 days to the current date.
		dateString := date.Format("2006-01-02")

		if _, _, err := ValidateDates(dateString, ""); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("begin date after end date", func(t *testing.T) {
		if _, _, err := ValidateDates("2024-01-01", "2023-12-31"); err == nil {
			t.Error("expected error, got nil")
		}
	})
}
