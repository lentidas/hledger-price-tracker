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
	"github.com/lentidas/hledger-price-tracker/internal/flags"
	"time"
)

const apiFunctionTimeSeriesDailyAdjusted = "TIME_SERIES_DAILY_ADJUSTED" // Requires premium API key.

// TODO Comment functions and such.

type RawDailyAdjusted struct {
	MetaData   RawMetadataDaily             `json:"Meta Data"`
	TimeSeries map[string]RawPricesAdjusted `json:"Time Series (Daily)"`
}

type TypedDailyAdjusted struct {
	MetaData   TypedMetadataDaily
	TimeSeries map[time.Time]TypedPricesAdjusted
}

type DailyAdjusted struct {
	Raw   RawDailyAdjusted
	Typed TypedDailyAdjusted
}

func (obj *DailyAdjusted) TypeBody() error {
	// TODO

	return nil
}

func (obj *DailyAdjusted) GenerateOutput(body []byte, begin time.Time, end time.Time, format flags.OutputFormat) (string, error) {
	// TODO

	return "DailyAdjusted", nil
}
