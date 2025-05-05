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

package crypto

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/lentidas/hledger-price-tracker/internal/flags"
)

var formatRate = flags.OutputFormatHledger
var interval = flags.IntervalWeekly
var begin string
var end string
var full bool

// rateCmd represents the rate command.
var rateCmd = &cobra.Command{
	Use:   "rate",
	Short: "Get the historical exchange rate between a currency and a cryptocurrency",
	Long: `
hledger-price-tracker

Command to get the exchange rates between a currency and a cryptocurrency
in a certain time period and interval.

It returns the open, high, low, and close exchange rates for each each interval 
in the time period defined.

API documentation:
- https://www.alphavantage.co/documentation/#currency-daily
- https://www.alphavantage.co/documentation/#currency-weekly
- https://www.alphavantage.co/documentation/#currency-monthly`,

	// Require the user to provide at least one argument, which is the currency we want to convert from.
	Args: cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		// TODO
		fmt.Println("rate called")
	},
}

func init() {
	// Add this subcommand to the `currency` command palette.
	PaletteCmd.AddCommand(rateCmd)

	// Add flags to the `rate` subcommand.
	rateCmd.Flags().VarP(&formatRate, "format", "f", "format of the output (possible values are \"hledger\", \"json\", \"csv\", \"table\", \"table-long\")")
	rateCmd.Flags().VarP(&interval, "interval", "i", "interval between prices (possible values are \"daily\", \"weekly\", \"monthly\")")
	rateCmd.Flags().StringVarP(&begin, "begin", "b", "", "beginning of the time period (format YYYY-MM-DD) (does not apply to  \"json\" or \"csv\" output formats)")
	rateCmd.Flags().StringVarP(&end, "end", "e", "", "end of the time period (format YYYY-MM-DD) (does not apply to  \"json\" or \"csv\" output formats)")
	rateCmd.Flags().BoolVar(&full, "full", false, "for daily intervals, return all the data, otherwise return only the last 100 data points (does nothing for weekly or monthly intervals)")
}
