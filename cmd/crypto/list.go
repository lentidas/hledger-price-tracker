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

	"github.com/lentidas/hledger-price-tracker/internal/crypto/list"
	"github.com/lentidas/hledger-price-tracker/internal/flags"
)

// Define the output flag and set it to the default value.
var formatList = flags.OutputFormatTable

// listCmd represents the list command.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available digital currencies",
	Long: `
hledger-price-tracker

Command to list all available physical currencies.`,

	Run: func(cmd *cobra.Command, args []string) {
		output, err := list.Execute(formatList)
		cobra.CheckErr(err)
		fmt.Print(output)
	},
}

func init() {
	// Add this subcommand to the `currency` command palette.
	PaletteCmd.AddCommand(listCmd)

	// Add flags to the `list` subcommand.
	listCmd.Flags().VarP(&formatList, "format", "f", "format of the output (possible values are \"csv\" and \"table\")")
}
