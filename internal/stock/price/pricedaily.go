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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lentidas/hledger-price-tracker/internal/flags"
	"time"
)

const apiFunctionTimeSeriesDaily = "TIME_SERIES_DAILY"

// TODO Comment functions and such.

type RawDaily struct {
	MetaData   RawMetadataDaily     `json:"Meta Data"`
	TimeSeries map[string]RawPrices `json:"Time Series (Daily)"`
}

type TypedDaily struct {
	MetaData   TypedMetadataDaily
	TimeSeries map[time.Time]TypedPrices
}

type Daily struct {
	Raw   RawDaily
	Typed TypedDaily
}

func (obj *Daily) TypeBody() error {
	err := obj.Typed.MetaData.TypeBody(obj.Raw.MetaData)
	if err != nil {
		return fmt.Errorf("[(*Daily).TypeBody] failure to cast metadata body: %w", err)
	}

	obj.Typed.TimeSeries = make(map[time.Time]TypedPrices)

	for date, prices := range obj.Raw.TimeSeries {
		dateTyped, err := time.Parse("2006-01-02", date)
		if err != nil {
			return fmt.Errorf("[(*Daily).TypeBody] error parsing date: %w", err)
		}

		var pricesTyped TypedPrices
		err = pricesTyped.TypeBody(prices)
		if err != nil {
			return fmt.Errorf("[(*Daily).TypeBody] error casting prices body: %w", err)
		}

		obj.Typed.TimeSeries[dateTyped] = pricesTyped
	}

	return nil
}

func (obj *Daily) GenerateOutput(body []byte, begin time.Time, end time.Time, format flags.OutputFormat) (string, error) {
	err := json.Unmarshal(body, &obj.Raw)
	if err != nil {
		return "", fmt.Errorf("[(*Daily).GenerateOutput] failure to unmarshal JSON body: %w", err)
	}

	switch format {
	case flags.OutputFormatJSON, flags.OutputFormatCSV:
		return string(body), nil
	case flags.OutputFormatHledger, flags.OutputFormatTable, flags.OutputFormatTableLong:
		err = obj.TypeBody()
		if err != nil {
			return "", fmt.Errorf("[(*Daily).GenerateOutput] error casting response attributes: %w", err)
		}

		dates := getDatesNormal(obj.Typed.TimeSeries, begin, end)

		fmt.Println(dates) // TODO Remove

		return string(body), nil // TODO Remove
	default:
		return "", errors.New("[(*Daily).GenerateOutput] invalid output format")
	}
}
