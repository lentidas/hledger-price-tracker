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

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	// Subcommands and other internal dependencies.
	"github.com/lentidas/hledger-price-tracker/cmd/currency"
	"github.com/lentidas/hledger-price-tracker/cmd/stock"
	"github.com/lentidas/hledger-price-tracker/internal"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hledger-price-tracker",
	Short: "A CLI tool to get market prices for commodities",
	Long: `hledger-price-tracker

hledger-price-tracker is a CLI program written in Go used to generate
market price records for hledger using the Alpha Vantage API.`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("root called") // TODO Remove these debug lines
		fmt.Println()

		// TODO Consider if we print the help message or not. It already seems to be the case when there is an error.
		// Print the help message by default on the root command / bare application.
		// err := cmd.Help()
		// if err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// addSubcommandPalettes adds all subcommand palettes to the root command.
// Although this could be done entirely inside the init() function of the root command, it separated for clarity.
func addSubcommandPalettes() {
	rootCmd.AddCommand(stock.PaletteCmd)
	rootCmd.AddCommand(currency.PaletteCmd)
}

func init() {
	// TODO Validate the timezone flag against a list of valid timezones.
	rootCmd.PersistentFlags().StringVarP(&internal.Timezone, "timezone", "z", "UTC", "Timezone to use for the generated records")
	// TODO Validate the currency flag against a list of valid currencies accepted by the Alpha Vantage API
	//  https://www.alphavantage.co/physical_currency_list/
	//  https://www.alphavantage.co/digital_currency_list/
	rootCmd.PersistentFlags().StringVarP(&internal.DefaultCurrency, "default-currency", "c", "EUR", "Default currency to use for the generated records")

	rootCmd.PersistentFlags().StringVarP(&internal.ApiKey, "api-key", "a", "", "API key for the Alpha Vantage API (required)")

	if err := rootCmd.MarkPersistentFlagRequired("api-key"); err != nil {
		os.Exit(1)
	}

	addSubcommandPalettes()

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hledger-price-tracker.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
