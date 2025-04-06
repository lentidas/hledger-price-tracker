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
	"strings"
	"time"

	"github.com/lentidas/hledger-price-tracker/internal/flags"
)

const apiFunctionTimeSeriesMonthly = "TIME_SERIES_MONTHLY"

// TODO Comment functions and such.

type RawMonthly struct {
	MetaData   RawMetadata          `json:"Meta Data"`
	TimeSeries map[string]RawPrices `json:"Monthly Time Series"`
}

type TypedMonthly struct {
	MetaData   TypedMetadata
	TimeSeries map[time.Time]TypedPrices
}

type Monthly struct {
	Raw   RawMonthly
	Typed TypedMonthly
}

func (obj *Monthly) TypeBody() error {
	err := obj.Typed.MetaData.TypeBody(obj.Raw.MetaData)
	if err != nil {
		return fmt.Errorf("[(*Monthly).TypeBody] failure to cast metadata body: %w", err)
	}

	obj.Typed.TimeSeries = make(map[time.Time]TypedPrices, len(obj.Raw.TimeSeries))

	for date, prices := range obj.Raw.TimeSeries {
		dateTyped, err := time.Parse("2006-01-02", date)
		if err != nil {
			return fmt.Errorf("[(*Monthly).TypeBody] error parsing date: %w", err)
		}

		var pricesTyped TypedPrices
		err = pricesTyped.TypeBody(prices)
		if err != nil {
			return fmt.Errorf("[(*Monthly).TypeBody] error casting prices body: %w", err)
		}

		obj.Typed.TimeSeries[dateTyped] = pricesTyped
	}

	return nil
}

func (obj *Monthly) GenerateOutput(body []byte, begin time.Time, end time.Time, format flags.OutputFormat) (string, error) {
	switch format {
	case flags.OutputFormatJSON, flags.OutputFormatCSV:
		return string(body), nil
	case flags.OutputFormatHledger, flags.OutputFormatTable, flags.OutputFormatTableLong:
		// Parse the JSON body into the Raw struct.
		err := json.Unmarshal(body, &obj.Raw)
		if err != nil {
			return "", fmt.Errorf("[(*Monthly).GenerateOutput] failure to unmarshal JSON body: %w", err)
		}

		// Cast the attributes into proper types.
		err = obj.TypeBody()
		if err != nil {
			return "", fmt.Errorf("[(*Monthly).GenerateOutput] error casting response attributes: %w", err)
		}

		dates := getDatesNormal(obj.Typed.TimeSeries, begin, end)

		if format == flags.OutputFormatHledger {
			return generateOutputHledgerNormal(
					obj.Typed.TimeSeries,
					dates,
					obj.Typed.MetaData.Symbol,
					obj.Typed.MetaData.Currency),
				nil
		} else {
			out := strings.Builder{}
			out.WriteString(generateMetadataTable(
				obj.Typed.MetaData.Symbol,
				obj.Typed.MetaData.Currency,
				obj.Typed.MetaData.LastRefreshed,
				obj.Typed.MetaData.TimeZone))
			out.WriteString(generateTimeSeriesTableShort(
				obj.Typed.TimeSeries,
				dates))
			return out.String(), nil
		}
	default:
		return "", errors.New("[(*Monthly).GenerateOutput] invalid output format")
	}
}
