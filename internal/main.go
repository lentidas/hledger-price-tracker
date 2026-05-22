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
	"encoding/json"
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

// httpGet performs a single HTTP GET and returns the raw response body.
func httpGet(url string) ([]byte, error) {
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

// apiErrorMessage checks whether a response body is an Alpha Vantage error envelope
// (rate-limit, daily cap, invalid key, etc.) and returns the message, or an empty
// string when the body is a normal response.
func apiErrorMessage(body []byte) string {
	if len(body) == 0 || body[0] != '{' {
		return ""
	}
	var envelope struct {
		Note        string `json:"Note"`
		Information string `json:"Information"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return ""
	}
	if envelope.Note != "" {
		return envelope.Note
	}
	return envelope.Information
}

// HTTPRequest makes an HTTP GET request and returns the body as a byte slice.
// This is the main entry point for all API calls, so it also checks for Alpha Vantage error envelopes
// and sends them back to the caller if an error happens.
// I've hardcoded a single retry after one second to circumvent the per-second rate-limit that Alpha Vantage added
// recently. If an user uses a premium API key they won't be affected by this, but free users will get a transparent
// retry when they hit the limit.
// After a single retry, if the response is still an error envelope, the message is returned directly to the caller.
func HTTPRequest(url string) ([]byte, error) {
	body, err := httpGet(url)
	if err != nil {
		return []byte{}, err
	}

	if msg := apiErrorMessage(body); msg != "" {
		time.Sleep(time.Second)
		body, err = httpGet(url)
		if err != nil {
			return []byte{}, err
		}
		if msg = apiErrorMessage(body); msg != "" {
			return []byte{}, fmt.Errorf("[internal.HTTPRequest] Alpha Vantage API error: %s", msg)
		}
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
