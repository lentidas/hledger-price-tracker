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
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lentidas/hledger-price-tracker/internal"
	"github.com/lentidas/hledger-price-tracker/internal/flags"
	"github.com/lentidas/hledger-price-tracker/internal/stock/search"
)

// TODO Comment functions and such.

type Response interface {
	TypeBody() error
	GenerateOutput(body []byte, begin time.Time, end time.Time, format flags.OutputFormat) (string, error)
}

type RawMetadata struct {
	Information string `json:"1. Information"`
	Symbol      string `json:"2. Symbol"`
	LastRefresh string `json:"3. Last Refreshed"`
	TimeZone    string `json:"4. Time Zone"`
}

type TypedMetadata struct {
	Information string
	Symbol      string
	Currency    string
	LastRefresh time.Time
	TimeZone    string
}

func (typed *TypedMetadata) TypeBody(raw RawMetadata) error {
	lastRefresh, err := time.Parse("2006-01-02", raw.LastRefresh)
	if err != nil {
		return err
	}

	typed.Information = raw.Information
	typed.Symbol = raw.Symbol
	typed.LastRefresh = lastRefresh
	typed.TimeZone = raw.TimeZone

	typed.Currency, err = getCurrency(typed.Symbol)
	if err != nil {
		return fmt.Errorf("[(*TypedMetadata).TypeBody] error getting currency: %w", err)
	}

	return nil
}

type RawMetadataDaily struct {
	Information string `json:"1. Information"`
	Symbol      string `json:"2. Symbol"`
	LastRefresh string `json:"3. Last Refreshed"`
	OutputSize  string `json:"4. Output Size"`
	TimeZone    string `json:"5. Time Zone"`
}

type TypedMetadataDaily struct {
	Information string
	Symbol      string
	Currency    string
	LastRefresh time.Time
	OutputSize  string
	TimeZone    string
}

func (typed *TypedMetadataDaily) TypeBody(raw RawMetadataDaily) error {
	lastRefresh, err := time.Parse("2006-01-02", raw.LastRefresh)
	if err != nil {
		return err
	}

	typed.Information = raw.Information
	typed.Symbol = raw.Symbol
	typed.LastRefresh = lastRefresh
	typed.TimeZone = raw.TimeZone

	typed.Currency, err = getCurrency(typed.Symbol)
	if err != nil {
		return fmt.Errorf("[(*TypedMetadata).TypeBody] error getting currency: %w", err)
	}

	return nil
}

type RawPrices struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

type TypedPrices struct {
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume uint32
}

func (typed *TypedPrices) TypeBody(raw RawPrices) error {
	// TODO Maybe wrap these errors?
	openPrice, err := strconv.ParseFloat(raw.Open, 64)
	if err != nil {
		return err
	}
	highPrice, err := strconv.ParseFloat(raw.High, 64)
	if err != nil {
		return err
	}
	lowPrice, err := strconv.ParseFloat(raw.Low, 64)
	if err != nil {
		return err
	}
	closePrice, err := strconv.ParseFloat(raw.Close, 64)
	if err != nil {
		return err
	}
	volume, err := strconv.ParseUint(raw.Volume, 10, 32)
	if err != nil {
		return err
	}

	typed.Open = openPrice
	typed.High = highPrice
	typed.Low = lowPrice
	typed.Close = closePrice
	typed.Volume = uint32(volume)

	return nil
}

type RawPricesAdjusted struct {
	Open             string `json:"1. open"`
	High             string `json:"2. high"`
	Low              string `json:"3. low"`
	Close            string `json:"4. close"`
	AdjustedClose    string `json:"5. adjusted close"`
	Volume           string `json:"6. volume"`
	DividendAmount   string `json:"7. dividend amount"`
	SplitCoefficient string `json:"8. split coefficient,omitempty"`
}

type TypedPricesAdjusted struct {
	Open             float64
	High             float64
	Low              float64
	Close            float64
	AdjustedClose    float64
	Volume           uint32
	DividendAmount   float64
	SplitCoefficient float64
}

func (typed *TypedPricesAdjusted) TypeBody(raw RawPricesAdjusted) error {
	// TODO Maybe wrap these errors?
	openPrice, err := strconv.ParseFloat(raw.Open, 64)
	if err != nil {
		return err
	}
	highPrice, err := strconv.ParseFloat(raw.High, 64)
	if err != nil {
		return err
	}
	lowPrice, err := strconv.ParseFloat(raw.Low, 64)
	if err != nil {
		return err
	}
	closePrice, err := strconv.ParseFloat(raw.Close, 64)
	if err != nil {
		return err
	}
	adjustedClose, err := strconv.ParseFloat(raw.AdjustedClose, 64)
	if err != nil {
		return err
	}
	volume, err := strconv.ParseUint(raw.Volume, 10, 32)
	if err != nil {
		return err
	}
	dividendAmount, err := strconv.ParseFloat(raw.DividendAmount, 64)
	if err != nil {
		return err
	}
	splitCoefficient, err := strconv.ParseFloat(raw.SplitCoefficient, 64)
	if err != nil {
		return err
	}

	typed.Open = openPrice
	typed.High = highPrice
	typed.Low = lowPrice
	typed.Close = closePrice
	typed.AdjustedClose = adjustedClose
	typed.Volume = uint32(volume)
	typed.DividendAmount = dividendAmount
	typed.SplitCoefficient = splitCoefficient

	return nil
}

// buildURL creates the URL to make the HTTP request to the Alpha Vantage API.
func buildURL(symbol string, format flags.OutputFormat, interval flags.Interval, adjusted bool) (string, error) {
	if internal.ApiKey == "" {
		return "", errors.New("[price.buildURL] api key is required")
	}
	if symbol == "" {
		return "", errors.New("[price.buildURL] no search query provided")
	}
	switch format {
	case flags.OutputFormatHledger, flags.OutputFormatTable, flags.OutputFormatTableLong, flags.OutputFormatJSON, flags.OutputFormatCSV:
		// Do nothing.
	default:
		return "", errors.New("[price.buildURL] invalid output format")
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
			return "", errors.New("[internal.price.buildURL] invalid interval for adjusted prices")
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
			return "", errors.New("[internal.price.buildURL] invalid interval")
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

// createResponseObject creates a Response object with the proper struct that will be used to parse
// the JSON body, depending on the interval and whether the prices are adjusted or not.
func createResponseObject(interval flags.Interval, adjusted bool) (Response, error) {
	var obj Response

	if adjusted {
		switch interval {
		case flags.IntervalDaily:
			obj = &DailyAdjusted{}
		case flags.IntervalWeekly:
			obj = &WeeklyAdjusted{}
		case flags.IntervalMonthly:
			obj = &MonthlyAdjusted{}
		default:
			return obj, errors.New("[internal.stock.createResponseObject] invalid interval for adjusted prices")
		}
	} else {
		switch interval {
		case flags.IntervalDaily:
			obj = &Daily{}
		case flags.IntervalWeekly:
			obj = &Weekly{}
		case flags.IntervalMonthly:
			obj = &Monthly{}
		default:
			return obj, errors.New("[internal.stock.createResponseObject] invalid interval")
		}
	}

	return obj, nil
}

func getCurrency(symbol string) (string, error) {
	// Alpha Vantage does not use the same symbol on their demo examples, so we need to return a default value.
	if internal.DebugMode || internal.ApiKey == "demo" {
		return "NIL", nil
	}

	body, err := search.Execute(symbol, flags.OutputFormatJSON)
	if err != nil {
		return "", err
	}
	var response search.Raw
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return "", fmt.Errorf("[internal.stock.getCurrency] error unmarshalling JSON to get currency: %w", err)
	} else if len(response.BestMatches) == 0 {
		return "", fmt.Errorf("[internal.stock.getCurrency] error getting currency of symbol %s", symbol)
	}

	return response.BestMatches[0].Currency, nil
}

func getDatesNormal(timeSeries map[time.Time]TypedPrices, begin time.Time, end time.Time) []time.Time {
	var dates []time.Time
	for date := range timeSeries {
		if !(date.Before(begin) || date.After(end)) {
			dates = append(dates, date)
		}
	}

	// TODO Implement sorting in chronologically AND reverse-chronologically order.
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	return dates
}

func getDatesAdjusted(timeSeries map[time.Time]TypedPricesAdjusted, begin time.Time, end time.Time) []time.Time {
	var dates []time.Time
	for date := range timeSeries {
		if !(date.Before(begin) || date.After(end)) {
			dates = append(dates, date)
		}
	}

	// TODO Implement sorting in chronologically AND reverse-chronologically order.
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	return dates
}

// TODO Continue implementing unitary tests for this
func Price(symbol string, format flags.OutputFormat, interval flags.Interval, begin string, end string, adjusted bool) (string, error) {
	// Verify function parameters and variables.
	if internal.ApiKey == "" {
		return "", errors.New("[stock.price.Price] API key is required")
	}
	if symbol == "" {
		return "", errors.New("[stock.price.Price] no stock symbol provided")
	}

	url, err := buildURL(symbol, format, interval, adjusted)
	if err != nil {
		return "", err
	}

	body, err := internal.HTTPRequest(url)
	if err != nil {
		return "", err
	}

	response, err := createResponseObject(interval, adjusted)
	if err != nil {
		return "", err
	}

	var beginTime time.Time
	var endTime time.Time = time.Now()
	if begin != "" {
		beginTime, err = time.Parse("2006-01-02", begin)
		if err != nil {
			return "", fmt.Errorf("[stock.Price] failed to parse begin date: %w", err)
		}
	}
	if end != "" {
		endTime, err = time.Parse("2006-01-02", end)
		if err != nil {
			return "", fmt.Errorf("[stock.Price] failed to parse end date: %w", err)
		}
	}

	return response.GenerateOutput(body, beginTime, endTime, format)
}
