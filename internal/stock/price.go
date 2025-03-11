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
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lentidas/hledger-price-tracker/internal"
	"github.com/lentidas/hledger-price-tracker/internal/flags"
)

const apiFunctionTimeSeriesDaily = "TIME_SERIES_DAILY"
const apiFunctionTimeSeriesDailyAdjusted = "TIME_SERIES_DAILY_ADJUSTED" // Requires premium API key.
const apiFunctionTimeSeriesWeekly = "TIME_SERIES_WEEKLY"
const apiFunctionTimeSeriesWeeklyAdjusted = "TIME_SERIES_WEEKLY_ADJUSTED"
const apiFunctionTimeSeriesMonthly = "TIME_SERIES_MONTHLY"
const apiFunctionTimeSeriesMonthlyAdjusted = "TIME_SERIES_MONTHLY_ADJUSTED"

type PriceResponseMetadataRaw struct {
	Information string `json:"1. Information"`
	Symbol      string `json:"2. Symbol"`
	LastRefresh string `json:"3. Last Refreshed"`
	TimeZone    string `json:"4. Time Zone"`
}

type PriceResponseMetadataTyped struct {
	Information string
	Symbol      string
	Currency    string
	LastRefresh time.Time
	TimeZone    string
}

func (typed *PriceResponseMetadataTyped) TypeBody(raw PriceResponseMetadataRaw) error {
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
		return fmt.Errorf("[internal.stock.PriceResponseMetadataTyped.TypeBody] error getting currency: %v", err)
	}

	return nil
}

type PriceResponseMetadataDailyRaw struct {
	Information string `json:"1. Information"`
	Symbol      string `json:"2. Symbol"`
	LastRefresh string `json:"3. Last Refreshed"`
	OutputSize  string `json:"4. Output Size"`
	TimeZone    string `json:"5. Time Zone"`
}

type PriceResponseMetadataDailyTyped struct {
	Information string
	Symbol      string
	Currency    string
	LastRefresh time.Time
	OutputSize  string
	TimeZone    string
}

func (typed *PriceResponseMetadataDailyTyped) TypeBody(raw PriceResponseMetadataDailyRaw) error {
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
		return fmt.Errorf("[internal.stock.PriceResponseMetadataTyped.TypeBody] error getting currency: %v", err)
	}

	return nil
}

type PriceResponsePricesRaw struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

type PriceResponsePricesTyped struct {
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume uint32
}

func (typed *PriceResponsePricesTyped) TypeBody(raw PriceResponsePricesRaw) error {
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

type PriceResponsePricesAdjustedRaw struct {
	Open             string `json:"1. open"`
	High             string `json:"2. high"`
	Low              string `json:"3. low"`
	Close            string `json:"4. close"`
	AdjustedClose    string `json:"5. adjusted close"`
	Volume           string `json:"6. volume"`
	DividendAmount   string `json:"7. dividend amount"`
	SplitCoefficient string `json:"8. split coefficient,omitempty"`
}

type PriceResponsePricesAdjustedTyped struct {
	Open             float64
	High             float64
	Low              float64
	Close            float64
	AdjustedClose    float64
	Volume           uint32
	DividendAmount   float64
	SplitCoefficient float64
}

func (typed *PriceResponsePricesAdjustedTyped) TypeBody(raw PriceResponsePricesAdjustedRaw) error {
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

type PriceResponseDailyRaw struct {
	MetaData   PriceResponseMetadataDailyRaw     `json:"Meta Data"`
	TimeSeries map[string]PriceResponsePricesRaw `json:"Time Series (Daily)"`
}

type PriceResponseDailyTyped struct {
	MetaData   PriceResponseMetadataDailyTyped
	TimeSeries map[time.Time]PriceResponsePricesTyped
}

type PriceResponseDaily struct {
	Raw   PriceResponseDailyRaw
	Typed PriceResponseDailyTyped
}

func (obj *PriceResponseDaily) UnmarshalJSON(body []byte) error {
	err := json.Unmarshal(body, &obj.Raw)
	if err != nil {
		return fmt.Errorf("[internal.stock.UnmarshalJSON] error unmarshalling JSON: %v", err)
	} else {
		return nil
	}
}

func (obj *PriceResponseDaily) TypeBody() error {
	err := obj.Typed.MetaData.TypeBody(obj.Raw.MetaData)
	if err != nil {
		return fmt.Errorf("[internal.stock.TypeBody] failure to cast metadata body: %v", err)
	}

	obj.Typed.TimeSeries = make(map[time.Time]PriceResponsePricesTyped)

	for date, prices := range obj.Raw.TimeSeries {
		var typedPrices PriceResponsePricesTyped
		err = typedPrices.TypeBody(prices)
		if err != nil {
			return fmt.Errorf("[internal.stock.TypeBody]: failure to cast price body: %v", err)
		}

		date, err := time.Parse("2006-01-02", date)
		if err != nil {
			return fmt.Errorf("[internal.stock.TypeBody]: failure to cast price date: %v", err)
		}

		obj.Typed.TimeSeries[date] = typedPrices
	}

	return nil
}

type PriceResponseDailyAdjustedRaw struct {
	MetaData   PriceResponseMetadataDailyRaw             `json:"Meta Data"`
	TimeSeries map[string]PriceResponsePricesAdjustedRaw `json:"Time Series (Daily)"`
}

type PriceResponseDailyAdjustedTyped struct {
	MetaData   PriceResponseMetadataDailyTyped
	TimeSeries map[time.Time]PriceResponsePricesAdjustedTyped
}

type PriceResponseDailyAdjusted struct {
	Raw   PriceResponseDailyAdjustedRaw
	Typed PriceResponseDailyAdjustedTyped
}

func (obj *PriceResponseDailyAdjusted) UnmarshalJSON(body []byte) error {
	err := json.Unmarshal(body, &obj.Raw)
	if err != nil {
		return fmt.Errorf("[internal.stock.UnmarshalJSON] error unmarshalling JSON: %v", err)
	} else {
		return nil
	}
}

func (obj *PriceResponseDailyAdjusted) TypeBody() error {
	err := obj.Typed.MetaData.TypeBody(obj.Raw.MetaData)
	if err != nil {
		return fmt.Errorf("[internal.stock.TypeBody] failure to cast metadata body: %v", err)
	}

	obj.Typed.TimeSeries = make(map[time.Time]PriceResponsePricesAdjustedTyped)

	for date, prices := range obj.Raw.TimeSeries {
		var typedPrices PriceResponsePricesAdjustedTyped
		err = typedPrices.TypeBody(prices)
		if err != nil {
			return fmt.Errorf("[internal.stock.TypeBody]: failure to cast price body: %v", err)
		}

		date, err := time.Parse("2006-01-02", date)
		if err != nil {
			return fmt.Errorf("[internal.stock.TypeBody]: failure to cast price date: %v", err)
		}

		obj.Typed.TimeSeries[date] = typedPrices
	}

	return nil
}

type PriceResponseWeeklyRaw struct {
	MetaData   PriceResponseMetadataRaw          `json:"Meta Data"`
	TimeSeries map[string]PriceResponsePricesRaw `json:"Weekly Time Series"`
}

type PriceResponseWeeklyTyped struct {
	MetaData   PriceResponseMetadataTyped
	TimeSeries map[time.Time]PriceResponsePricesTyped
}

type PriceResponseWeekly struct {
	Raw   PriceResponseWeeklyRaw
	Typed PriceResponseWeeklyTyped
}

func (obj *PriceResponseWeekly) UnmarshalJSON(body []byte) error {
	err := json.Unmarshal(body, &obj.Raw)
	if err != nil {
		return fmt.Errorf("[internal.stock.UnmarshalJSON] error unmarshalling JSON: %v", err)
	} else {
		return nil
	}
}

func (obj *PriceResponseWeekly) TypeBody() error {
	err := obj.Typed.MetaData.TypeBody(obj.Raw.MetaData)
	if err != nil {
		return fmt.Errorf("[internal.stock.TypeBody] failure to cast metadata body: %v", err)
	}

	obj.Typed.TimeSeries = make(map[time.Time]PriceResponsePricesTyped)

	for date, prices := range obj.Raw.TimeSeries {
		var typedPrices PriceResponsePricesTyped
		err = typedPrices.TypeBody(prices)
		if err != nil {
			return fmt.Errorf("[internal.stock.TypeBody]: failure to cast price body: %v", err)
		}

		date, err := time.Parse("2006-01-02", date)
		if err != nil {
			return fmt.Errorf("[internal.stock.TypeBody]: failure to cast price date: %v", err)
		}

		obj.Typed.TimeSeries[date] = typedPrices
	}

	return nil
}

type PriceResponseWeeklyAdjustedRaw struct {
	MetaData   PriceResponseMetadataRaw                  `json:"Meta Data"`
	TimeSeries map[string]PriceResponsePricesAdjustedRaw `json:"Weekly Adjusted Time Series"`
}

type PriceResponseWeeklyAdjustedTyped struct {
	MetaData   PriceResponseMetadataTyped
	TimeSeries map[time.Time]PriceResponsePricesAdjustedTyped
}

type PriceResponseWeeklyAdjusted struct {
	Raw   PriceResponseWeeklyAdjustedRaw
	Typed PriceResponseWeeklyAdjustedTyped
}

func (obj *PriceResponseWeeklyAdjusted) UnmarshalJSON(body []byte) error {
	err := json.Unmarshal(body, &obj.Raw)
	if err != nil {
		return fmt.Errorf("[internal.stock.UnmarshalJSON] error unmarshalling JSON: %v", err)
	} else {
		return nil
	}
}

func (obj *PriceResponseWeeklyAdjusted) TypeBody() error {
	err := obj.Typed.MetaData.TypeBody(obj.Raw.MetaData)
	if err != nil {
		return fmt.Errorf("[internal.stock.TypeBody] failure to cast metadata body: %v", err)
	}

	obj.Typed.TimeSeries = make(map[time.Time]PriceResponsePricesAdjustedTyped)

	for date, prices := range obj.Raw.TimeSeries {
		var typedPrices PriceResponsePricesAdjustedTyped
		err = typedPrices.TypeBody(prices)
		if err != nil {
			return fmt.Errorf("[internal.stock.TypeBody]: failure to cast price body: %v", err)
		}

		date, err := time.Parse("2006-01-02", date)
		if err != nil {
			return fmt.Errorf("[internal.stock.TypeBody]: failure to cast price date: %v", err)
		}

		obj.Typed.TimeSeries[date] = typedPrices
	}

	return nil
}

type PriceResponseMonthlyRaw struct {
	MetaData   PriceResponseMetadataRaw          `json:"Meta Data"`
	TimeSeries map[string]PriceResponsePricesRaw `json:"Monthly Time Series"`
}

type PriceResponseMonthlyTyped struct {
	MetaData   PriceResponseMetadataTyped
	TimeSeries map[time.Time]PriceResponsePricesTyped
}

type PriceResponseMonthly struct {
	Raw   PriceResponseMonthlyRaw
	Typed PriceResponseMonthlyTyped
}

func (obj *PriceResponseMonthly) UnmarshalJSON(body []byte) error {
	err := json.Unmarshal(body, &obj.Raw)
	if err != nil {
		return fmt.Errorf("[internal.stock.UnmarshalJSON] error unmarshalling JSON: %v", err)
	} else {
		return nil
	}
}

func (obj *PriceResponseMonthly) TypeBody() error {
	err := obj.Typed.MetaData.TypeBody(obj.Raw.MetaData)
	if err != nil {
		return fmt.Errorf("[internal.stock.TypeBody] failure to cast metadata body: %v", err)
	}

	obj.Typed.TimeSeries = make(map[time.Time]PriceResponsePricesTyped)

	for date, prices := range obj.Raw.TimeSeries {
		var typedPrices PriceResponsePricesTyped
		err = typedPrices.TypeBody(prices)
		if err != nil {
			return fmt.Errorf("[internal.stock.TypeBody]: failure to cast price body: %v", err)
		}

		date, err := time.Parse("2006-01-02", date)
		if err != nil {
			return fmt.Errorf("[internal.stock.TypeBody]: failure to cast price date: %v", err)
		}

		obj.Typed.TimeSeries[date] = typedPrices
	}

	return nil
}

type PriceResponseMonthlyAdjustedRaw struct {
	MetaData   PriceResponseMetadataRaw                  `json:"Meta Data"`
	TimeSeries map[string]PriceResponsePricesAdjustedRaw `json:"Monthly Adjusted Time Series"`
}

type PriceResponseMonthlyAdjustedTyped struct {
	MetaData   PriceResponseMetadataTyped
	TimeSeries map[time.Time]PriceResponsePricesAdjustedTyped
}

type PriceResponseMonthlyAdjusted struct {
	Raw   PriceResponseMonthlyAdjustedRaw
	Typed PriceResponseMonthlyAdjustedTyped
}

func (obj *PriceResponseMonthlyAdjusted) UnmarshalJSON(body []byte) error {
	err := json.Unmarshal(body, &obj.Raw)
	if err != nil {
		return fmt.Errorf("[internal.stock.UnmarshalJSON] error unmarshalling JSON: %v", err)
	} else {
		return nil
	}
}

func (obj *PriceResponseMonthlyAdjusted) TypeBody() error {
	err := obj.Typed.MetaData.TypeBody(obj.Raw.MetaData)
	if err != nil {
		return fmt.Errorf("[internal.stock.TypeBody] failure to cast metadata body: %v", err)
	}

	obj.Typed.TimeSeries = make(map[time.Time]PriceResponsePricesAdjustedTyped)

	for date, prices := range obj.Raw.TimeSeries {
		var typedPrices PriceResponsePricesAdjustedTyped
		err = typedPrices.TypeBody(prices)
		if err != nil {
			return fmt.Errorf("[internal.stock.TypeBody]: failure to cast price body: %v", err)
		}

		date, err := time.Parse("2006-01-02", date)
		if err != nil {
			return fmt.Errorf("[internal.stock.TypeBody]: failure to cast price date: %v", err)
		}

		obj.Typed.TimeSeries[date] = typedPrices
	}

	return nil
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

func getCurrency(symbol string) (string, error) {
	// Alpha Vantage does not use the same symbol on their demo examples, so we need to return a default value.
	if internal.DebugMode || internal.ApiKey == "demo" {
		return "NIL", nil
	}

	body, err := Search(symbol, flags.OutputFormatJSON)
	if err != nil {
		return "", err
	}
	var search SearchResponseRaw
	err = json.Unmarshal([]byte(body), &search)
	if err != nil {
		return "", fmt.Errorf("[internal.stock.getCurrency] error unmarshalling JSON to get currency: %v", err)
	} else if len(search.BestMatches) == 0 {
		return "", fmt.Errorf("[internal.stock.getCurrency] error getting currency of symbol %s", symbol)
	}

	return search.BestMatches[0].Currency, nil
}

// createResponseObject creates a internal.JSONResponse object with the proper struct that will be used to parse
// the JSON body, depending on the interval and whether the prices are adjusted or not.
func createResponseObject(body []byte, interval flags.Interval, adjusted bool) (internal.JSONResponse, error) {
	var obj internal.JSONResponse

	if adjusted {
		switch interval {
		case flags.IntervalDaily:
			obj = &PriceResponseDailyAdjusted{}
		case flags.IntervalWeekly:
			obj = &PriceResponseWeeklyAdjusted{}
		case flags.IntervalMonthly:
			obj = &PriceResponseMonthlyAdjusted{}
		default:
			return obj, errors.New("[internal.stock.createResponseObject] invalid interval for adjusted prices")
		}
	} else {
		switch interval {
		case flags.IntervalDaily:
			obj = &PriceResponseDaily{}
		case flags.IntervalWeekly:
			obj = &PriceResponseWeekly{}
		case flags.IntervalMonthly:
			obj = &PriceResponseMonthly{}
		default:
			return obj, errors.New("[internal.stock.createResponseObject] invalid interval")
		}
	}

	// Parse the JSON body.
	err := obj.UnmarshalJSON(body)
	if err != nil {
		return nil, fmt.Errorf("[internal.stock.generatePriceOutput] failed to unmarshal JSON body: %v", err)
	}

	// Cast the attributes into proper types.
	err = obj.TypeBody()
	if err != nil {
		return nil, fmt.Errorf("[internal.stock.generatePriceOutput] failed to cast body attributes: %v", err)
	}

	return obj, nil
}

// func getDates(begin string, end string, response internal.JSONResponse) ([]time.Time, error) {
// 	// Convert the `begin` and `end` strings into time.Time objects.
// 	var beginTime, endTime time.Time
// 	var err error
// 	if begin != "" {
// 		beginTime, err = time.Parse("2006-01-02", begin)
// 		if err != nil {
// 			return nil, fmt.Errorf("[internal.stock.getDates] failed to parse begin date: %v", err)
// 		}
// 	}
// 	if end != "" {
// 		endTime, err = time.Parse("2006-01-02", end)
// 		if err != nil {
// 			return nil, fmt.Errorf("[internal.stock.getDates] failed to parse end date: %v", err)
// 		}
// 	}
//
// 	// Extract the keys and sort them
// 	var keys []time.Time
// 	for k := range obj. {
// 		keys = append(keys, k)
// 	}
// 	sort.Slice(keys, func(i, j int) bool {
// 		return keys[i].Before(keys[j])
// 	})
// }

// generatePriceOutput generates the output in the desired format by either outputting the raw string of the body or by
// parsing its JSON content then creating a proper table.
// TODO Create unitary tests for this function.
func generatePriceOutput(body []byte, format flags.OutputFormat, interval flags.Interval, begin string, end string, adjusted bool) (string, error) {
	switch format {
	case flags.OutputFormatHledger:
		// Create the object with the JSON response.
		obj, err := createResponseObject(body, interval, adjusted)
		if err != nil {
			return "", fmt.Errorf("[internal.stock.createResponseObject] failed to create response object: %v", err)
		}

		fmt.Println("obj:", obj) // TODO Remove this

		return "hledger", nil
	case flags.OutputFormatJSON, flags.OutputFormatCSV:
		return string(body), nil
	case flags.OutputFormatTable, flags.OutputFormatTableLong:
		// Create the object with the JSON response.
		obj, err := createResponseObject(body, interval, adjusted)
		if err != nil {
			return "", fmt.Errorf("[internal.stock.createResponseObject] failed to create response object: %v", err)
		}

		fmt.Println("obj:", obj) // TODO Remove this

		return "table", nil
	}

	// TODO
	return "", nil
}

// TODO Continue implementing unitary tests for this
func Price(symbol string, format flags.OutputFormat, interval flags.Interval, begin string, end string, adjusted bool) (string, error) {
	// TODO Remove these debug lines
	fmt.Println("[DEBUG Price()]")
	fmt.Println("symbol:", symbol)
	fmt.Println("format:", format)
	fmt.Println("interval:", interval)
	fmt.Println("begin:", begin)
	fmt.Println("end:", end)
	fmt.Println("adjusted:", adjusted)
	fmt.Println("[END_DEBUG]")

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

	body, err := internal.HTTPRequest(url)
	if err != nil {
		return "", err
	}

	return generatePriceOutput(body, format, interval, begin, end, adjusted)
}
