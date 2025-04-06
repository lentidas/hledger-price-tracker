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
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const ApiBaseUrl string = "https://www.alphavantage.co/query?"

var ApiKey string
var DefaultCurrency string
var DebugMode bool

// HTTPRequest is a helper function to make HTTP requests and return the body as an array of bytes.
func HTTPRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("[internal.HTTPRequest] HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("[internal.HTTPRequest] failure to read HTTP body: %w", err)
	}

	return body, nil
}

func ParseCurrenciesCSV(body []byte) (map[string]string, error) {
	// Read CSV data.
	csvReader := csv.NewReader(bytes.NewReader(body))
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("[internal.ParseCurrenciesCSV] failure to read CSV data: %w", err)
	}

	currencies := make(map[string]string)

	// Parse CSV data into a slice of Currency objects.
	for i, line := range data {
		if i > 0 { // Skip the header line.
			var code string
			var name string
			for j, field := range line {
				switch j {
				case 0:
					code = field
				case 1:
					name = field
				}
			}
			currencies[code] = name
		}
	}

	return currencies, nil
}

// ValidateDates checks if begin and end dates are valid then returns its parsed values
// in the time.Time format.
func ValidateDates(begin, end string) (time.Time, time.Time, error) {
	var err error
	var beginTime time.Time
	var endTime time.Time = time.Now()
	if begin != "" {
		beginTime, err = time.Parse("2006-01-02", begin)
		if err != nil {
			return time.Now(), time.Now(), fmt.Errorf("[internal.ValidateDates] failed to parse begin date: %w", err)
		}
		if beginTime.After(time.Now()) {
			return time.Now(), time.Now(), errors.New("[internal.ValidateDates] begin date is in the future")
		}
	}
	if end != "" {
		endTime, err = time.Parse("2006-01-02", end)
		if err != nil {
			return time.Now(), time.Now(), fmt.Errorf("[internal.ValidateDates] failed to parse end date: %w", err)
		}
	}
	if beginTime.After(endTime) {
		return time.Now(), time.Now(), errors.New("[internal.ValidateDates] begin date is after end date")
	}

	return beginTime, endTime, nil
}
