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
	"testing"
)

func TestDecodeBody(t *testing.T) {
	type TestStruct struct {
		Name    string `json:"name"`
		Value   int    `json:"value"`
		Boolean bool   `json:"boolean"`
	}

	t.Run("success", func(t *testing.T) {
		jsonData := []byte(`{"name":"test","value":123, "boolean":true}`)
		var result TestStruct

		// Call function
		err := DecodeBody(jsonData, &result)

		// Check results
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result.Name != "test" {
			t.Errorf("expected name 'test', got '%s'", result.Name)
		}
		if result.Value != 123 {
			t.Errorf("expected value 123, got %d", result.Value)
		}
		if result.Boolean != true {
			t.Errorf("expected boolean true, got %v", result.Boolean)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		// Create test data with invalid JSON.
		invalidJSON := []byte(`{"name":"test", invalid json}`)
		var result TestStruct

		err := DecodeBody(invalidJSON, &result)

		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("non-pointer receiver", func(t *testing.T) {
		jsonData := []byte(`{"name":"test"}`)
		var result TestStruct

		err := DecodeBody(jsonData, result) // Passing value, not pointer.

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
