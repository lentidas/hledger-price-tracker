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
	"github.com/spf13/viper"

	// Subcommands and other internal dependencies.
	"github.com/lentidas/hledger-price-tracker/cmd/currency"
	"github.com/lentidas/hledger-price-tracker/cmd/stock"
	"github.com/lentidas/hledger-price-tracker/internal"
)

var cfgFile string

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
		fmt.Println(err) // FIXME Is there a better way to show error messages?
		os.Exit(1)
	}
}

func init() {
	// Initialize configuration with Viper.
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "path to config file (default is $HOME/.hledger-price-tracker.yaml)")

	// TODO Validate the currency flag against a list of valid currencies accepted by the Alpha Vantage API
	//  https://www.alphavantage.co/physical_currency_list/
	//  https://www.alphavantage.co/digital_currency_list/
	rootCmd.PersistentFlags().StringVarP(&internal.DefaultCurrency, "default-currency", "c", "EUR", "default currency to use for the generated records")

	rootCmd.PersistentFlags().StringVarP(&internal.ApiKey, "api-key", "a", "", "API key for the Alpha Vantage API")

	// Add all subcommand palettes to the root command.
	rootCmd.AddCommand(stock.PaletteCmd)
	rootCmd.AddCommand(currency.PaletteCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".hledger-price-tracker" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".hledger-price-tracker")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
