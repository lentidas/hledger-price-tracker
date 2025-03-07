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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const ApiBaseUrl string = "https://www.alphavantage.co/query?"

var ApiKey string
var DefaultCurrency string

type JSONResponseInterface interface {
	DecodeBody(body []byte) error
}

type JSONResponse struct {
	Content any
}

func (response *JSONResponse) DecodeBody(body []byte) error {
	err := json.Unmarshal(body, response.Content)
	if err != nil {
		return fmt.Errorf("[internal.DecodeBody] failure to unmarshal JSON body: %v", err)
	} else {
		return nil
	}
}

type JSONResponseTypedInterface interface {
	TypeBody(response JSONResponse) error
}

type JSONResponseTyped struct {
	Content any
}

func (responseTyped *JSONResponseTyped) TypeBody(response JSONResponse) error {
	// TODO
	return nil
}

// HTTPRequest is a helper function to make HTTP requests and return the body as an array of bytes.
func HTTPRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("[internal.HTTPRequest] HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("[internal.HTTPRequest] failure to read HTTP body: %v", err)
	}

	return body, nil
}
