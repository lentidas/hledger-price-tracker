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

package rate

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/lentidas/hledger-price-tracker/internal"
	currencyList "github.com/lentidas/hledger-price-tracker/internal/currency/list"
	"github.com/lentidas/hledger-price-tracker/internal/flags"
)

// TODO Add documentation to each of the functions

type Response interface {
	TypeBody() error
	GenerateOutput(body []byte, begin time.Time, end time.Time, format flags.OutputFormat) (string, error)
}

type RawMetadata struct {
	Information   string `json:"1. Information"`
	FromSymbol    string `json:"2. From Symbol"`
	ToSymbol      string `json:"3. To Symbol"`
	LastRefreshed string `json:"4. Last Refreshed"`
	TimeZone      string `json:"5. Time Zone"`
}

type TypedMetadata struct {
	Information   string
	FromSymbol    string
	ToSymbol      string
	LastRefreshed time.Time
	TimeZone      string
}

func (typed *TypedMetadata) TypeBody(raw RawMetadata) error {
	var lastRefreshed time.Time
	var err error

	// Try parsing with date and time format first.
	lastRefreshed, err = time.Parse("2006-01-02 15:04:05", raw.LastRefreshed)
	if err != nil {
		// If that fails, try just the date format.
		lastRefreshed, err = time.Parse("2006-01-02", raw.LastRefreshed)
		if err != nil {
			return fmt.Errorf("[currency.rate.(*TypedMetadataDaily).TypeBody] error parsing last refreshed time: %w", err)
		}
	}

	typed.Information = raw.Information
	typed.FromSymbol = raw.FromSymbol
	typed.ToSymbol = raw.ToSymbol
	typed.LastRefreshed = lastRefreshed
	typed.TimeZone = raw.TimeZone

	return nil
}

type RawMetadataDaily struct {
	Information   string `json:"1. Information"`
	FromSymbol    string `json:"2. From Symbol"`
	ToSymbol      string `json:"3. To Symbol"`
	OutputSize    string `json:"4. Output Size"`
	LastRefreshed string `json:"5. Last Refreshed"`
	TimeZone      string `json:"6. Time Zone"`
}

type TypedMetadataDaily struct {
	Information   string
	FromSymbol    string
	ToSymbol      string
	OutputSize    string
	LastRefreshed time.Time
	TimeZone      string
}

func (typed *TypedMetadataDaily) TypeBody(raw RawMetadataDaily) error {
	var lastRefreshed time.Time
	var err error

	// Try parsing with date and time format first.
	lastRefreshed, err = time.Parse("2006-01-02 15:04:05", raw.LastRefreshed)
	if err != nil {
		// If that fails, try just the date format.
		lastRefreshed, err = time.Parse("2006-01-02", raw.LastRefreshed)
		if err != nil {
			return fmt.Errorf("[currency.rate.(*TypedMetadataDaily).TypeBody] error parsing last refreshed time: %w", err)
		}
	}

	typed.Information = raw.Information
	typed.FromSymbol = raw.FromSymbol
	typed.ToSymbol = raw.ToSymbol
	typed.OutputSize = raw.OutputSize
	typed.LastRefreshed = lastRefreshed
	typed.TimeZone = raw.TimeZone

	return nil
}

type RawPrices struct {
	Open  string `json:"1. open"`
	High  string `json:"2. high"`
	Low   string `json:"3. low"`
	Close string `json:"4. close"`
}

type TypedPrices struct {
	Open  float64
	High  float64
	Low   float64
	Close float64
}

func (typed *TypedPrices) TypeBody(raw RawPrices) error {
	openPrice, err := strconv.ParseFloat(raw.Open, 64)
	if err != nil {
		return fmt.Errorf("[currency.rate.(*TypedPrices).TypeBody] error parsing open price: %w", err)
	}
	highPrice, err := strconv.ParseFloat(raw.High, 64)
	if err != nil {
		return fmt.Errorf("[currency.rate.(*TypedPrices).TypeBody] error parsing high price: %w", err)
	}
	lowPrice, err := strconv.ParseFloat(raw.Low, 64)
	if err != nil {
		return fmt.Errorf("[currency.rate.(*TypedPrices).TypeBody] error parsing low price: %w", err)
	}
	closePrice, err := strconv.ParseFloat(raw.Close, 64)
	if err != nil {
		return fmt.Errorf("[currency.rate.(*TypedPrices).TypeBody] error parsing close price: %w", err)
	}

	typed.Open = openPrice
	typed.High = highPrice
	typed.Low = lowPrice
	typed.Close = closePrice

	return nil
}

func buildURL(from string, to string, format flags.OutputFormat, interval flags.Interval, full bool) (string, error) {
	if internal.ApiKey == "" {
		return "", errors.New("[currency.rate.buildURL] API key is required")
	}

	fromBoolCurrency, err := currencyList.CurrencyExists(from)
	if err != nil {
		return "", err
	} else if !fromBoolCurrency {
		return "", errors.New("[currency.rate.buildURL] from currency is not valid")
	}

	toBoolCurrency, err := currencyList.CurrencyExists(to)
	if err != nil {
		return "", err
	} else if !toBoolCurrency {
		return "", errors.New("[currency.rate.buildURL] to currency is not valid")
	}

	if from == to {
		return "", errors.New("[currency.rate.buildURL] from and to currencies must be different")
	}

	switch format {
	case flags.OutputFormatHledger, flags.OutputFormatTable, flags.OutputFormatTableLong, flags.OutputFormatJSON, flags.OutputFormatCSV:
		// Do nothing.
	default:
		return "", errors.New("[currency.rate.buildURL] invalid output format")
	}

	url := strings.Builder{}
	url.WriteString(internal.ApiBaseUrl)
	url.WriteString("function=")

	switch interval {
	case flags.IntervalDaily:
		url.WriteString(apiFunctionCurrencyRateDaily)
	case flags.IntervalWeekly:
		url.WriteString(apiFunctionCurrencyRateWeekly)
	case flags.IntervalMonthly:
		url.WriteString(apiFunctionCurrencyRateMonthly)
	default:
		return "", errors.New("[currency.rate.buildURL] invalid interval for adjusted prices")
	}

	url.WriteString("&from_symbol=")
	url.WriteString(from)
	url.WriteString("&to_symbol=")
	url.WriteString(to)

	// Print entire time series if daily, because user can then limit the interval with `begin` and `end`.
	if interval == flags.IntervalDaily && full {
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
// the JSON body, depending on the interval.
func createResponseObject(interval flags.Interval) (Response, error) {
	var obj Response

	switch interval {
	case flags.IntervalDaily:
		obj = &Daily{}
	case flags.IntervalWeekly:
		obj = &Weekly{}
	case flags.IntervalMonthly:
		obj = &Monthly{}
	default:
		return obj, errors.New("[currency.rate.createResponseObject] invalid interval")
	}

	return obj, nil
}

// getDates returns the dates in the time series that are within the specified interval.
func getDates(timeSeries map[time.Time]TypedPrices, begin time.Time, end time.Time) []time.Time {
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

// generateOutputHledger generates the output in hledger format for non-adjusted prices.
func generateOutputHledger(timeSeries map[time.Time]TypedPrices, dates []time.Time, from string, to string) string {
	out := strings.Builder{}
	for _, date := range dates {
		out.WriteString(fmt.Sprintf("P %s \"%s\" %.2f \"%s\"\n",
			date.Format("2006-01-02"),
			from,
			timeSeries[date].Close,
			to))
	}
	return out.String()
}

func generateMetadataTable(from string, to string, lastRefreshed time.Time) string {
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"From", "To", "Last Refreshed"})
	t.AppendRow(table.Row{from, to, lastRefreshed.Format("2006-01-02 15:04:05")})
	return t.Render() + "\n"
}

func generateTimeSeriesTable(timeSeries map[time.Time]TypedPrices, dates []time.Time) string {
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Date", "Open", "High", "Low", "Close"})
	for _, date := range dates {
		prices := timeSeries[date]
		t.AppendRow([]interface{}{
			date.Format("2006-01-02"),
			fmt.Sprintf("%.2f", prices.Open),
			fmt.Sprintf("%.2f", prices.High),
			fmt.Sprintf("%.2f", prices.Low),
			fmt.Sprintf("%.2f", prices.Close),
		})
	}
	return t.Render() + "\n"
}

func Execute(from string, to string, format flags.OutputFormat, interval flags.Interval, begin string, end string, full bool) (string, error) {
	beginTime, endTime, err := internal.ValidateDates(begin, end)
	if err != nil {
		return "", err
	}

	url, err := buildURL(from, to, format, interval, full)
	if err != nil {
		return "", err
	}

	body, err := internal.HTTPRequest(url)
	if err != nil {
		return "", err
	}

	response, err := createResponseObject(interval)
	if err != nil {
		return "", err
	}

	return response.GenerateOutput(body, beginTime, endTime, format)
}
