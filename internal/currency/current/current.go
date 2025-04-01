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

package current

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/lentidas/hledger-price-tracker/internal"
	cryptoList "github.com/lentidas/hledger-price-tracker/internal/crypto/list"
	currencyList "github.com/lentidas/hledger-price-tracker/internal/currency/list"
	"github.com/lentidas/hledger-price-tracker/internal/flags"
)

const apiFunctionSearch = "CURRENCY_EXCHANGE_RATE"

type Response interface {
	TypeBody() error
	GenerateOutput(body []byte, format flags.OutputFormat) (string, error)
}

type Raw struct {
	RealtimeCurrencyExchangeRate struct {
		FromCurrencyCode string `json:"1. From_Currency Code"`
		FromCurrencyName string `json:"2. From_Currency Name"`
		ToCurrencyCode   string `json:"3. To_Currency Code"`
		ToCurrencyName   string `json:"4. To_Currency Name"`
		ExchangeRate     string `json:"5. Exchange Rate"`
		LastRefreshed    string `json:"6. Last Refreshed"`
		TimeZone         string `json:"7. Time Zone"`
		BidPrice         string `json:"8. Bid Price"`
		AskPrice         string `json:"9. Ask Price"`
	} `json:"Realtime Currency Exchange Rate"`
}

type Typed struct {
	FromCurrencyCode string
	FromCurrencyName string
	ToCurrencyCode   string
	ToCurrencyName   string
	ExchangeRate     float64
	LastRefreshed    time.Time
	TimeZone         string
	BidPrice         float64
	AskPrice         float64
}

type Current struct {
	Raw   Raw
	Typed Typed
}

func (obj *Current) TypeBody() error {
	lastRefreshed, err := time.Parse("2006-01-02 15:04:05", obj.Raw.RealtimeCurrencyExchangeRate.LastRefreshed)
	if err != nil {
		return fmt.Errorf("[(*Current).TypeBody] error parsing last refreshed date: %w", err)
	}

	exchangeRate, err := strconv.ParseFloat(obj.Raw.RealtimeCurrencyExchangeRate.ExchangeRate, 64)
	if err != nil {
		return fmt.Errorf("[(*Current).TypeBody] error parsing exchange rate: %w", err)
	}
	bidPrice, err := strconv.ParseFloat(obj.Raw.RealtimeCurrencyExchangeRate.BidPrice, 64)
	if err != nil {
		return fmt.Errorf("[(*Current).TypeBody] error parsing bid price: %w", err)
	}
	askPrice, err := strconv.ParseFloat(obj.Raw.RealtimeCurrencyExchangeRate.AskPrice, 64)
	if err != nil {
		return fmt.Errorf("[(*Current).TypeBody] error parsing ask price: %w", err)
	}

	obj.Typed.FromCurrencyCode = obj.Raw.RealtimeCurrencyExchangeRate.FromCurrencyCode
	obj.Typed.FromCurrencyName = obj.Raw.RealtimeCurrencyExchangeRate.FromCurrencyName
	obj.Typed.ToCurrencyCode = obj.Raw.RealtimeCurrencyExchangeRate.ToCurrencyCode
	obj.Typed.ToCurrencyName = obj.Raw.RealtimeCurrencyExchangeRate.ToCurrencyName
	obj.Typed.ExchangeRate = exchangeRate
	obj.Typed.LastRefreshed = lastRefreshed
	obj.Typed.TimeZone = obj.Raw.RealtimeCurrencyExchangeRate.TimeZone
	obj.Typed.BidPrice = bidPrice
	obj.Typed.AskPrice = askPrice

	return nil
}

func (obj *Current) GenerateOutput(body []byte, format flags.OutputFormat) (string, error) {
	switch format {
	case flags.OutputFormatCSV:
		return "", errors.New("[(*Current).GenerateOutput] CSV output format not supported")
	case flags.OutputFormatJSON:
		return string(body), nil
	case flags.OutputFormatHledger, flags.OutputFormatTable, flags.OutputFormatTableLong:
		// Parse the JSON body into the Raw struct.
		err := json.Unmarshal(body, &obj.Raw)
		if err != nil {
			return "", fmt.Errorf("[(*Current).GenerateOutput] failure to unmarshal JSON body: %w", err)
		}

		// Cast the attributes into proper types.
		err = obj.TypeBody()
		if err != nil {
			return "", fmt.Errorf("[(*Current).GenerateOutput] error casting response attributes: %w", err)
		}

		if format == flags.OutputFormatHledger {
			output := fmt.Sprintf("P %s \"%s\" %.2f \"%s\"\n",
				obj.Typed.LastRefreshed.Format("2006-01-02"),
				obj.Typed.FromCurrencyCode,
				obj.Typed.ExchangeRate,
				obj.Typed.ToCurrencyCode)
			return output, nil
		} else {
			tMetadata := table.NewWriter()
			tMetadata.SetStyle(table.StyleLight)
			tMetadata.AppendHeader(table.Row{"From", "To", "Last Refreshed"})
			tMetadata.AppendRow(table.Row{
				fmt.Sprintf("%s (%s)", obj.Typed.FromCurrencyName, obj.Typed.FromCurrencyCode),
				fmt.Sprintf("%s (%s)", obj.Typed.ToCurrencyName, obj.Typed.ToCurrencyCode),
				obj.Typed.LastRefreshed.Format("2006-01-02 15:04:05"),
			})

			tData := table.NewWriter()
			tData.SetStyle(table.StyleLight)
			if format == flags.OutputFormatTable {
				tData.AppendHeader(table.Row{
					obj.Typed.FromCurrencyCode,
					obj.Typed.ToCurrencyCode,
				})
				tData.AppendRow(table.Row{"1", obj.Typed.ExchangeRate})
			} else {
				tData.AppendHeader(table.Row{
					obj.Typed.FromCurrencyCode,
					obj.Typed.ToCurrencyCode,
					"Bid Price",
					"Ask Price",
				})
				tData.AppendRow(table.Row{
					"1",
					obj.Typed.ExchangeRate,
					obj.Typed.BidPrice,
					obj.Typed.AskPrice,
				})
			}

			return tMetadata.Render() + "\n" + tData.Render() + "\n", nil
		}
	default:
		return "", errors.New("[(*Current).GenerateOutput] invalid output format")
	}
}

func buildURL(from string, to string) (string, error) {
	if internal.ApiKey == "" {
		return "", errors.New("[currency/crypto.current.buildURL] API key is required")
	}
	if from == "" {
		return "", errors.New("[currency/crypto.current.buildURL] from currency is required")
	}
	if to == "" {
		return "", errors.New("[currency/crypto.current.buildURL] to currency is required")
	}

	// Validate the currency/crypto code.
	fromBoolCurrency, fromErrorCurrency := currencyList.CurrencyExists(from)
	toBoolCurrency, toErrorCurrency := currencyList.CurrencyExists(to)
	fromBoolCrypto, fromErrorCrypto := cryptoList.CryptoExists(from)
	toBoolCrypto, toErrorCrypto := cryptoList.CryptoExists(to)
	if fromErrorCurrency != nil {
		return "", fromErrorCurrency
	} else if toErrorCurrency != nil {
		return "", toErrorCurrency
	} else if fromErrorCrypto != nil {
		return "", fromErrorCrypto
	} else if toErrorCrypto != nil {
		return "", toErrorCrypto
	}
	if !fromBoolCurrency && !fromBoolCrypto {
		return "", errors.New("[currency/crypto.current.buildURL] from currency is not valid")
	}
	if !toBoolCurrency && !toBoolCrypto {
		return "", errors.New("[currency/crypto.current.buildURL] to currency is not valid")
	}

	if from == to {
		return "", errors.New("[currency/crypto.current.buildURL] from and to currencies must be different")
	}

	url := strings.Builder{}
	url.WriteString(internal.ApiBaseUrl)
	url.WriteString("function=")
	url.WriteString(apiFunctionSearch)
	url.WriteString("&from_currency=")
	url.WriteString(from)
	url.WriteString("&to_currency=")
	url.WriteString(to)
	url.WriteString("&apikey=")
	url.WriteString(internal.ApiKey)

	return url.String(), nil
}

func Execute(from string, to string, format flags.OutputFormat) (string, error) {
	url, err := buildURL(from, to)
	if err != nil {
		return "", err
	}

	body, err := internal.HTTPRequest(url)
	if err != nil {
		return "", err
	}

	response := Current{}

	return response.GenerateOutput(body, format)
}
