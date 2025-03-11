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

const apiFunctionTimeSeriesWeeklyAdjusted = "TIME_SERIES_WEEKLY_ADJUSTED"

// TODO Comment functions and such.

type RawWeeklyAdjusted struct {
	MetaData   RawMetadata                  `json:"Meta Data"`
	TimeSeries map[string]RawPricesAdjusted `json:"Weekly Adjusted Time Series"`
}

type TypedWeeklyAdjusted struct {
	MetaData   TypedMetadata
	TimeSeries map[time.Time]TypedPricesAdjusted
}

type WeeklyAdjusted struct {
	Raw   RawWeeklyAdjusted
	Typed TypedWeeklyAdjusted
}

func (obj *WeeklyAdjusted) TypeBody() error {
	// TODO

	return nil
}

func (obj *WeeklyAdjusted) GenerateOutput(body []byte, begin time.Time, end time.Time, format flags.OutputFormat) (string, error) {
	// TODO

	return "WeeklyAdjusted", nil
}
