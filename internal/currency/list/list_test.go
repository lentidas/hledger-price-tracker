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

package list

import "testing"

func TestListCurrencyExists(t *testing.T) {
	// Define a slice of currencies to test.
	currencies := []string{"EUR", "CHF", "GBP", "USD"}

	for _, currency := range currencies {
		t.Run("currency exists "+currency, func(t *testing.T) {
			result, err := CurrencyExists(currency)
			if err != nil {
				t.Errorf("expected nil, got %v", err)
			} else if !result {
				t.Errorf("expected %s to exist, got false", currency)
			}
		})
	}

	// Define a slice of cryptocurrencies to test.
	cryptos := []string{"BTC", "ETH", "LTC", "XRP"}

	for _, crypto := range cryptos {
		t.Run("crypto does not exist "+crypto, func(t *testing.T) {
			result, err := CurrencyExists(crypto)
			if err != nil {
				t.Errorf("expected nil, got %v", err)
			} else if result {
				t.Errorf("expected %s to not exist, got true", crypto)
			}
		})
	}

	// Test for empty currency.
	t.Run("failure empty", func(t *testing.T) {
		result, err := CurrencyExists("")
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		} else if result {
			t.Errorf("expected false, got true")
		}
	})

	// Test for non-existent currency.
	t.Run("failure INVALID", func(t *testing.T) {
		result, err := CurrencyExists("INVALID")
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		} else if result {
			t.Errorf("expected false, got true")
		}
	})
}
