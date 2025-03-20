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
	"os"

	"github.com/spf13/cobra"
)

// PaletteCmd represents the currency command palette.
var PaletteCmd = &cobra.Command{
	Use:     "crypto",
	GroupID: "palette",
	Short:   "Palette command that groups all subcommands related to cryptocurrencies",
	Long: `
hledger-price-tracker

Palette command that groups all subcommands related to cryptocurrencies and
their exchange rates.`,

	Run: func(cmd *cobra.Command, args []string) {
		// Print the help message for this command palette.
		err := cmd.Help()
		if err != nil {
			os.Exit(1)
		}
	},
}

func init() {}
