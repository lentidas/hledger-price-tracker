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
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <https://www.gnu.org/licenses/>.
 */

package stock

import (
	"errors"
	"io"
	"net/http"
	"strings"

	// Import internal dependencies.
	"github.com/lentidas/hledger-price-tracker/internal"
	"github.com/lentidas/hledger-price-tracker/internal/flags"
)

const apiFunctionSearch = "SYMBOL_SEARCH"

func Search(query string, output flags.OutputFormat) (string, error) {
	// Verify function parameters and variables.
	if internal.ApiKey == "" {
		return "", errors.New("[stock.Search] api key is required")
	}
	if query == "" {
		return "", errors.New("[stock.Search] no search query provided")
	}

	// Create the request URL.
	url := strings.Builder{}
	url.WriteString(internal.ApiBaseUrl)
	url.WriteString("function=")
	url.WriteString(apiFunctionSearch)
	url.WriteString("&keywords=")
	url.WriteString(query)
	url.WriteString("&apikey=")
	url.WriteString(internal.ApiKey)

	// Perform the HTTP request.
	resp, err := http.Get(url.String())
	if err != nil {
		// TODO: Handle the error. Quit the program? Log the error?
	}
	defer resp.Body.Close()

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Return the output in the desired format.
	switch output {
	case flags.OutputFormatJson:
		return string(body), nil
	case flags.OutputFormatTable:
		return "table", nil
	case flags.OutputFormatTableShort:
		return "table-short", nil
	case flags.OutputFormatCsv:
		return "csv", nil
	default:
		return "", errors.New("[stock.Search] invalid output format")
	}
}
