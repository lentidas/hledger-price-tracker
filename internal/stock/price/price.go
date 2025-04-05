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
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/lentidas/hledger-price-tracker/internal"
	"github.com/lentidas/hledger-price-tracker/internal/flags"
	"github.com/lentidas/hledger-price-tracker/internal/stock/search"
)

type Response interface {
	TypeBody() error
	GenerateOutput(body []byte, begin time.Time, end time.Time, format flags.OutputFormat) (string, error)
}

type RawMetadata struct {
	Information   string `json:"1. Information"`
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	TimeZone      string `json:"4. Time Zone"`
}

type TypedMetadata struct {
	Information   string
	Symbol        string
	Currency      string
	LastRefreshed time.Time
	TimeZone      string
}

func (typed *TypedMetadata) TypeBody(raw RawMetadata) error {
	lastRefreshed, err := time.Parse("2006-01-02", raw.LastRefreshed)
	if err != nil {
		return err
	}

	typed.Information = raw.Information
	typed.Symbol = raw.Symbol
	typed.LastRefreshed = lastRefreshed
	typed.TimeZone = raw.TimeZone

	typed.Currency, err = search.GetCurrency(typed.Symbol)
	if err != nil {
		return fmt.Errorf("[stock.price.(*TypedMetadata).TypeBody] error getting currency: %w", err)
	}

	return nil
}

type RawMetadataDaily struct {
	Information   string `json:"1. Information"`
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	OutputSize    string `json:"4. Output Size"`
	TimeZone      string `json:"5. Time Zone"`
}

type TypedMetadataDaily struct {
	Information   string
	Symbol        string
	Currency      string
	LastRefreshed time.Time
	OutputSize    string
	TimeZone      string
}

func (typed *TypedMetadataDaily) TypeBody(raw RawMetadataDaily) error {
	lastRefreshed, err := time.Parse("2006-01-02", raw.LastRefreshed)
	if err != nil {
		return fmt.Errorf("[stock.price.(*TypedMetadataDaily).TypeBody] error parsing last refreshed date: %w", err)
	}

	typed.Information = raw.Information
	typed.Symbol = raw.Symbol
	typed.LastRefreshed = lastRefreshed
	typed.TimeZone = raw.TimeZone

	typed.Currency, err = search.GetCurrency(typed.Symbol)
	if err != nil {
		return fmt.Errorf("[stock.price.(*TypedMetadataDaily).TypeBody] error getting currency: %w", err)
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
	openPrice, err := strconv.ParseFloat(raw.Open, 64)
	if err != nil {
		return fmt.Errorf("[stock.price.(*TypedPrices).TypeBody] error parsing open price: %w", err)
	}
	highPrice, err := strconv.ParseFloat(raw.High, 64)
	if err != nil {
		return fmt.Errorf("[stock.price.(*TypedPrices).TypeBody] error parsing high price: %w", err)
	}
	lowPrice, err := strconv.ParseFloat(raw.Low, 64)
	if err != nil {
		return fmt.Errorf("[stock.price.(*TypedPrices).TypeBody] error parsing low price: %w", err)
	}
	closePrice, err := strconv.ParseFloat(raw.Close, 64)
	if err != nil {
		return fmt.Errorf("[stock.price.(*TypedPrices).TypeBody] error parsing close price: %w", err)
	}
	volume, err := strconv.ParseUint(raw.Volume, 10, 32)
	if err != nil {
		return fmt.Errorf("[stock.price.(*TypedPrices).TypeBody] error parsing volume: %w", err)
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
	SplitCoefficient string `json:"-"`
}

type TypedPricesAdjusted struct {
	Open           float64
	High           float64
	Low            float64
	Close          float64
	AdjustedClose  float64
	Volume         uint32
	DividendAmount float64
}

func (typed *TypedPricesAdjusted) TypeBody(raw RawPricesAdjusted) error {
	openPrice, err := strconv.ParseFloat(raw.Open, 64)
	if err != nil {
		return fmt.Errorf("[stock.price.(*TypedPricesAdjusted).TypeBody] error parsing open price: %w", err)
	}
	highPrice, err := strconv.ParseFloat(raw.High, 64)
	if err != nil {
		return fmt.Errorf("[stock.price.(*TypedPricesAdjusted).TypeBody] error parsing high price: %w", err)
	}
	lowPrice, err := strconv.ParseFloat(raw.Low, 64)
	if err != nil {
		return fmt.Errorf("[stock.price.(*TypedPricesAdjusted).TypeBody] error parsing low price: %w", err)
	}
	closePrice, err := strconv.ParseFloat(raw.Close, 64)
	if err != nil {
		return fmt.Errorf("[stock.price.(*TypedPricesAdjusted).TypeBody] error parsing close price: %w", err)
	}
	adjustedClose, err := strconv.ParseFloat(raw.AdjustedClose, 64)
	if err != nil {
		return fmt.Errorf("[stock.price.(*TypedPricesAdjusted).TypeBody] error parsing adjusted close price: %w", err)
	}
	volume, err := strconv.ParseUint(raw.Volume, 10, 32)
	if err != nil {
		return fmt.Errorf("[stock.price.(*TypedPricesAdjusted).TypeBody] error parsing volume: %w", err)
	}
	dividendAmount, err := strconv.ParseFloat(raw.DividendAmount, 64)
	if err != nil {
		return fmt.Errorf("[stock.price.(*TypedPricesAdjusted).TypeBody] error parsing dividend amount: %w", err)
	}

	typed.Open = openPrice
	typed.High = highPrice
	typed.Low = lowPrice
	typed.Close = closePrice
	typed.AdjustedClose = adjustedClose
	typed.Volume = uint32(volume)
	typed.DividendAmount = dividendAmount

	return nil
}

// buildURL creates the URL to make the HTTP request to the Alpha Vantage API.
func buildURL(symbol string, format flags.OutputFormat, interval flags.Interval, adjusted bool, full bool) (string, error) {
	if internal.ApiKey == "" {
		return "", errors.New("[price.buildURL] API key is required")
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
			return obj, errors.New("[stock.price.createResponseObject] invalid interval for adjusted prices")
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
			return obj, errors.New("[stock.price.createResponseObject] invalid interval")
		}
	}

	return obj, nil
}

// getDatesNormal returns the dates in the time series that are within the specified interval.
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

// getDatesAdjusted returns the dates for the adjusted prices that are within the interval defined by `begin` and `end`.
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

// generateOutputHledger generates the output in hledger format for non-adjusted prices.
func generateOutputHledger(timeSeries map[time.Time]TypedPrices, dates []time.Time, symbol string, currency string) string {
	out := strings.Builder{}
	for _, date := range dates {
		out.WriteString(fmt.Sprintf("P %s \"%s\" %.2f %s\n",
			date.Format("2006-01-02"),
			symbol,
			timeSeries[date].Close,
			currency))
	}
	return out.String()
}

// generateOutputHledgerAdjusted generates the output in hledger format for adjusted prices.
func generateOutputHledgerAdjusted(timeSeries map[time.Time]TypedPricesAdjusted, dates []time.Time, symbol string, currency string) string {
	out := strings.Builder{}
	for _, date := range dates {
		out.WriteString(fmt.Sprintf("P %s \"%s\" %.2f %s\n",
			date.Format("2006-01-02"),
			symbol,
			timeSeries[date].AdjustedClose, // Adjusted close price instead of close price.
			currency))
	}
	return out.String()
}

// generateMetadataTable generates a table with the metadata for a given stock symbol. It is used to display
// the information about the stock before the table with the stock prices.
func generateMetadataTable(symbol string, currency string, lastRefreshed time.Time, timeZone string) string {
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Symbol", "Currency", "Last Refreshed", "Timezone"})
	t.AppendRow(table.Row{symbol, currency, lastRefreshed.Format("2006-01-02"), timeZone})
	return t.Render() + "\n"
}

// generateTimeSeriesTableShort generates a short table with the prices for a given stock symbol.
// It is used to display the stock prices in a compact way.
// Note that for non-adjusted prices this output format is used both in `table` and `table-long`.
func generateTimeSeriesTableShort(timeSeries map[time.Time]TypedPrices, dates []time.Time) string {
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Date", "Open", "High", "Low", "Close", "Volume"})
	for _, date := range dates {
		prices := timeSeries[date]
		t.AppendRow([]interface{}{
			date.Format("2006-01-02"),
			fmt.Sprintf("%.2f", prices.Open),
			fmt.Sprintf("%.2f", prices.High),
			fmt.Sprintf("%.2f", prices.Low),
			fmt.Sprintf("%.2f", prices.Close),
			prices.Volume,
		})
	}
	return t.Render() + "\n"
}

// generateTimeSeriesTableShortAdjusted generates a table with the adjusted prices for a given stock symbol.
// It is used to display the stock prices in a compact way, but only for adjusted prices output.
func generateTimeSeriesTableShortAdjusted(timeSeries map[time.Time]TypedPricesAdjusted, dates []time.Time) string {
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Date", "Open", "High", "Low", "Close", "Adj. Close", "Volume"})
	for _, date := range dates {
		prices := timeSeries[date]
		t.AppendRow([]interface{}{
			date.Format("2006-01-02"),
			fmt.Sprintf("%.2f", prices.Open),
			fmt.Sprintf("%.2f", prices.High),
			fmt.Sprintf("%.2f", prices.Low),
			fmt.Sprintf("%.2f", prices.Close),
			fmt.Sprintf("%.2f", prices.AdjustedClose),
			prices.Volume,
		})
	}
	return t.Render() + "\n"
}

// generateTimeSeriesTableLongAdjusted generates a long table with the adjusted prices for a given stock symbol.
// It is used to display the stock prices in a detailed way, but only for adjusted prices output.
func generateTimeSeriesTableLongAdjusted(timeSeries map[time.Time]TypedPricesAdjusted, dates []time.Time) string {
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Date", "Open", "High", "Low", "Close", "Adj. Close", "Volume", "Dividend Amount"})
	for _, date := range dates {
		prices := timeSeries[date]
		t.AppendRow([]interface{}{
			date.Format("2006-01-02"),
			fmt.Sprintf("%.2f", prices.Open),
			fmt.Sprintf("%.2f", prices.High),
			fmt.Sprintf("%.2f", prices.Low),
			fmt.Sprintf("%.2f", prices.Close),
			fmt.Sprintf("%.2f", prices.AdjustedClose),
			prices.Volume,
			fmt.Sprintf("%.2f", prices.DividendAmount),
		})
	}
	return t.Render() + "\n"
}

// TODO Continue implementing unitary tests for this

// Execute is the core function of the price package. It fetches the stock prices from the Alpha Vantage API for a given
// stock symbol and returns it in the desired format.
func Execute(symbol string, format flags.OutputFormat, interval flags.Interval, begin string, end string, adjusted bool, full bool) (string, error) {
	var err error
	var beginTime time.Time
	var endTime time.Time = time.Now()
	if begin != "" {
		beginTime, err = time.Parse("2006-01-02", begin)
		if err != nil {
			return "", fmt.Errorf("[stock.price.Execute] failed to parse begin date: %w", err)
		}
		if beginTime.After(time.Now()) {
			return "", errors.New("[stock.price.Execute] begin date is in the future")
		}
	}
	if end != "" {
		endTime, err = time.Parse("2006-01-02", end)
		if err != nil {
			return "", fmt.Errorf("[stock.price.Execute] failed to parse end date: %w", err)
		}
	}
	if beginTime.After(endTime) {
		return "", errors.New("[stock.price.Execute] begin date is after end date")
	}

	url, err := buildURL(symbol, format, interval, adjusted, full)
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

	return response.GenerateOutput(body, beginTime, endTime, format)
}
