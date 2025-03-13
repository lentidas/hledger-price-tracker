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
	"fmt"

	"github.com/spf13/cobra"

	"github.com/lentidas/hledger-price-tracker/internal/flags"
	"github.com/lentidas/hledger-price-tracker/internal/stock/price"
)

var formatPrice = flags.OutputFormatHledger
var interval = flags.IntervalWeekly
var begin string
var end string
var adjusted bool

// priceCmd represents the price command
var priceCmd = &cobra.Command{
	Use:   "price [flags] <stock-symbol>",
	Short: "Get prices for a given stock in a certain time period and interval",
	Long: `
hledger-price-tracker

Command to obtain prices for a given stock in a certain time period and interval.

It returns the open, high, low, close, and volume of the stock for each interval 
in the time period defined. Adjusted close prices are also available.

API documentation:
- https://www.alphavantage.co/documentation/#daily
- https://www.alphavantage.co/documentation/#dailyadj
- https://www.alphavantage.co/documentation/#weekly
- https://www.alphavantage.co/documentation/#weeklyadj
- https://www.alphavantage.co/documentation/#monthly
- https://www.alphavantage.co/documentation/#monthlyadj`,

	// Require the user to provide at least one argument, which is the stock symbol.
	Args: cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		output, err := price.Price(args[0], formatPrice, interval, begin, end, adjusted)
		cobra.CheckErr(err)
		fmt.Print(output)
	},

	// TODO Implement a way to output an error when the API does not find a stock symbol.
}

func init() {
	// Add this subcommand to the `stock` command palette.
	PaletteCmd.AddCommand(priceCmd)

	// Add flags to the `price` subcommand.
	priceCmd.Flags().VarP(&formatPrice, "format", "f", "format of the output (possible values are \"hledger\", \"json\", \"csv\", \"table\", \"table-long\")")
	priceCmd.Flags().VarP(&interval, "interval", "i", "interval between prices (possible values are \"daily\", \"weekly\", \"monthly\")")
	priceCmd.Flags().StringVarP(&begin, "begin", "b", "", "beginning of the time period (format YYYY-MM-DD) (does not apply to  \"json\" or \"csv\" output formats)")
	priceCmd.Flags().StringVarP(&end, "end", "e", "", "end of the time period (format YYYY-MM-DD) (does not apply to  \"json\" or \"csv\" output formats)")
	priceCmd.Flags().BoolVarP(&adjusted, "adjusted", "a", false, "return adjusted close prices")
}
