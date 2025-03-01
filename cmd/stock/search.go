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
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <https://www.gnu.org/licenses/>.
 */

package stock

import (
	"fmt"

	"github.com/spf13/cobra"

	// Import internal dependencies.
	"github.com/lentidas/hledger-price-tracker/internal/flags"
	"github.com/lentidas/hledger-price-tracker/internal/stock"
)

// Define the output flag and set it to the default value.
var format = flags.OutputFormatJson

// searchCmd represents the search command.
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "search command", // TODO
	Long:  `TODO`,

	// Require the user to provide at least one argument, which is the query for the stock search.
	Args: cobra.MinimumNArgs(1),

	// TODO Show example with the argument.
	Run: func(cmd *cobra.Command, args []string) {
		output, err := stock.Search(args[0], format)
		cobra.CheckErr(err)
		fmt.Println(output)
	},
}

func init() {
	// Add this subcommand to the `stock` command palette.
	PaletteCmd.AddCommand(searchCmd)

	// Add flags to the `search` subcommand.
	searchCmd.Flags().VarP(&format, "format", "f", "format of the output (possible values are \"json\", \"table\", \"table-short\", \"csv\", default is \"json\")")
}
