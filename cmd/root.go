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

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/lentidas/hledger-price-tracker/cmd/currency"
	"github.com/lentidas/hledger-price-tracker/cmd/stock"
	"github.com/lentidas/hledger-price-tracker/internal"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hledger-price-tracker",
	Short: "A CLI tool to get market prices for commodities",
	Long: `
hledger-price-tracker

hledger-price-tracker is a CLI program written in Go used to generate
market price records for hledger using the Alpha Vantage API.`,
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
	rootCmd.PersistentFlags().StringVarP(&internal.DefaultCurrency, "default-currency", "c", "EUR", "default currency to use for the generated ledger records")
	rootCmd.PersistentFlags().StringVarP(&internal.ApiKey, "api-key", "k", "", "API key to access the Alpha Vantage API")

	// Mark which flags are required for this subcommand.
	requiredPersistentFlags := []string{"api-key"}
	for _, flag := range requiredPersistentFlags {
		if err := rootCmd.MarkPersistentFlagRequired(flag); err != nil {
			// Ignore the error, because Cobra only returns an error if the flag is not previously created.
		}
	}

	// Add command palettes as subgroups to make the message clearer.
	rootCmd.AddGroup(&cobra.Group{ID: "palette", Title: "Command Palettes:"})

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
		viper.AddConfigPath(home + "/.config/hledger-price-tracker")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		// Do nothing if the config file is not found.
	}

	// Set a prefix for environment variables, to avoid conflicts with other programs.
	viper.SetEnvPrefix("HPT") // HPT for hledger-price-tracker

	// Environment variables can't have dashes in them, so bind them to their equivalent keys with underscores,
	// e.g. --api-key to HPT_API_KEY.
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Bind certain flags to environment variables.
	envFlags := []string{"api-key"}
	for _, flag := range envFlags {
		if err := viper.BindEnv(flag); err != nil {
			// Ignore the error, because Cobra only returns an error if the flag is not previously created.
		}
	}

	// Bind each cobra flag to its associated viper configuration (config file and environment variable).
	// The precedence order is: command-line flags -> environment variables -> config file.
	// IMPORTANT: This will only associate the flags from the root command, not the subcommands.
	// NOTE: This os based on the following repository: https://github.com/carolynvs/stingoftheviper/
	rootCmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Apply the viper config value to the flag when the flag is not set and viper has a value.
		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			rootCmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
