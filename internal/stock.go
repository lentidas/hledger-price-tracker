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

package internal

import (
	"io"
	"net/http"
	"strings"
)

func StockSearch(keyword string) string {
	// Create the request URL.
	url := strings.Builder{}
	url.WriteString(apiBaseUrl)
	url.WriteString("function=SYMBOL_SEARCH&keywords=")
	url.WriteString(keyword)
	url.WriteString("&apikey=")
	url.WriteString(ApiKey)

	// Perform the HTTP request.
	resp, err := http.Get(url.String())
	if err != nil {
		// TODO: Handle the error. Quit the program? Log the error?
	}
	defer resp.Body.Close()

	// Read the response body.
	body, err := io.ReadAll(resp.Body)

	// Return the response body as a string.
	return string(body)
}
