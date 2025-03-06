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

package stock

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lentidas/hledger-price-tracker/internal"
	"github.com/lentidas/hledger-price-tracker/internal/flags"
)

const apiFunctionTimeSeriesDaily = "TIME_SERIES_DAILY"
const apiFunctionTimeSeriesDailyAdjusted = "TIME_SERIES_DAILY_ADJUSTED" // Requires premium API key.
const apiFunctionTimeSeriesWeekly = "TIME_SERIES_WEEKLY"
const apiFunctionTimeSeriesWeeklyAdjusted = "TIME_SERIES_WEEKLY_ADJUSTED"
const apiFunctionTimeSeriesMonthly = "TIME_SERIES_MONTHLY"
const apiFunctionTimeSeriesMonthlyAdjusted = "TIME_SERIES_MONTHLY_ADJUSTED"

type priceResponseBase struct {
	internal.JSONResponse
	MetaData struct {
		Information string `json:"1. Information"`
		Symbol      string `json:"2. Symbol"`
		LastRefresh string `json:"3. Last Refreshed"`
		OutputSize  string `json:"-"`
		TimeZone    string `json:"5. Time Zone"`
	} `json:"Meta Data"`
}

type prices struct {
	Prices struct {
		Open   string `json:"1. open"`
		High   string `json:"2. high"`
		Low    string `json:"3. low"`
		Close  string `json:"4. close"`
		Volume string `json:"5. volume"`
	}
}

type pricesAdjusted struct {
	Prices struct {
		Open             string `json:"1. open"`
		High             string `json:"2. high"`
		Low              string `json:"3. low"`
		Close            string `json:"4. close"`
		AdjustedClose    string `json:"5. adjusted close"`
		Volume           string `json:"6. volume"`
		DividendAmount   string `json:"7. dividend amount"`
		SplitCoefficient string `json:"-"`
	}
}

type priceResponseDaily struct {
	priceResponseBase
	TimeSeriesDaily prices `json:"Time Series (Daily)"`
}

type priceResponseDailyAdjusted struct {
	priceResponseBase
	TimeSeriesDailyAdjusted pricesAdjusted `json:"Time Series (Daily)"`
}

type priceResponseWeekly struct {
	priceResponseBase
	TimeSeriesWeekly prices `json:"Weekly Time Series"`
}

type priceResponseWeeklyAdjusted struct {
	priceResponseBase
	TimeSeriesDailyAdjusted pricesAdjusted `json:"Weekly Adjusted Time Series"`
}

type priceResponseMonthly struct {
	priceResponseBase
	TimeSeriesWeekly prices `json:"Monthly Time Series"`
}

type priceResponseMonthlyAdjusted struct {
	priceResponseBase
	TimeSeriesDailyAdjusted pricesAdjusted `json:"Monthly Adjusted Time Series"`
}

// buildPriceURL creates the URL to make the HTTP request to the Alpha Vantage API.
func buildPriceURL(symbol string, format flags.OutputFormat, interval flags.Interval, adjusted bool) (string, error) {
	if internal.ApiKey == "" {
		return "", errors.New("[stock.buildPriceURL] api key is required")
	}
	if symbol == "" {
		return "", errors.New("[stock.buildPriceURL] no search query provided")
	}
	switch format {
	case flags.OutputFormatHledger, flags.OutputFormatTable, flags.OutputFormatTableLong, flags.OutputFormatJSON, flags.OutputFormatCSV:
		// Do nothing.
	default:
		return "", errors.New("[stock.buildPriceURL] invalid output format")
	}

	url := strings.Builder{}
	url.WriteString(internal.ApiBaseUrl)
	url.WriteString("function=")

	if adjusted {
		switch interval {
		case flags.IntervalDaily:
			url.WriteString(apiFunctionTimeSeriesDailyAdjusted)
		case flags.IntervalWeekly:
			url.WriteString(apiFunctionTimeSeriesWeeklyAdjusted)
		case flags.IntervalMonthly:
			url.WriteString(apiFunctionTimeSeriesMonthlyAdjusted)
		default:
			return "", errors.New("[internal.stock.buildPriceURL] invalid interval for adjusted prices")
		}
	} else {
		switch interval {
		case flags.IntervalDaily:
			url.WriteString(apiFunctionTimeSeriesDaily)
		case flags.IntervalWeekly:
			url.WriteString(apiFunctionTimeSeriesWeekly)
		case flags.IntervalMonthly:
			url.WriteString(apiFunctionTimeSeriesMonthly)
		default:
			return "", errors.New("[internal.stock.buildPriceURL] invalid interval")
		}
	}

	url.WriteString("&symbol=")
	url.WriteString(symbol)

	// Print entire time series if daily, because user can then limit the interval with `begin` and `end`.
	if interval == flags.IntervalDaily {
		url.WriteString("&outputsize=full")
	}
	url.WriteString("&apikey=")
	url.WriteString(internal.ApiKey)

	if format == flags.OutputFormatCSV {
		url.WriteString("&datatype=csv")
	}

	return url.String(), nil
}

// TODO Continue implementing unitary tests for this
func Price(symbol string, format flags.OutputFormat, interval flags.Interval, begin string, end string, adjusted bool) (string, error) {
	// TODO Remove these debug lines
	// fmt.Println("[DEBUG]")
	// fmt.Println("symbol:", symbol)
	// fmt.Println("format:", format)
	// fmt.Println("interval:", interval)
	// fmt.Println("begin:", begin)
	// fmt.Println("end:", end)
	// fmt.Println("adjusted:", adjusted)
	// fmt.Println("[END_DEBUG]")

	// Verify function parameters and variables.
	if internal.ApiKey == "" {
		return "", errors.New("[stock.Price] api key is required")
	}
	if symbol == "" {
		return "", errors.New("[stock.Price] no stock symbol provided")
	}

	url, err := buildPriceURL(symbol, format, interval, adjusted)
	if err != nil {
		return "", err
	}

	fmt.Println(url)

	// TODO
	return "", nil
}
