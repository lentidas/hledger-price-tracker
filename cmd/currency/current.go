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

package currency

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/lentidas/hledger-price-tracker/internal"
	"github.com/lentidas/hledger-price-tracker/internal/currency/current"
	"github.com/lentidas/hledger-price-tracker/internal/flags"
)

// Define the output flag and set it to the default value.
var formatCurrent = flags.OutputFormatHledger

// currentCmd represents the current command.
var currentCmd = &cobra.Command{
	Use:   "current [flags] <from-currency> [<to-currency>]",
	Short: "Get the current exchange rate between two currencies",
	Long: `
hledger-price-tracker

Command to get the current exchange rate between two currencies.
Can be used to convert between either physical or crypto currencies
(i.e. performs the exact same operation as the analogous command in the 
\'crypto\' command palette).

API documentation: https://www.alphavantage.co/documentation/#currency-exchange`,

	// Require the user to provide at least one argument, which is the currency we want to convert from.
	Args: cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		var to string
		if len(args) < 2 {
			to = internal.DefaultCurrency
		} else {
			to = args[1]
		}
		output, err := current.Execute(args[0], to, formatCurrent)
		cobra.CheckErr(err)
		fmt.Print(output)
	},
}

func init() {
	// Add this subcommand to the `currency` command palette.
	PaletteCmd.AddCommand(currentCmd)

	// Add flags to the `current` subcommand.
	currentCmd.Flags().VarP(&formatCurrent, "format", "f", "format of the output (possible values are \"hledger\", \"json\", \"table\", \"table-long\")")
}
