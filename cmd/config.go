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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command.
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show information about the application's configuration file",
	Long: `hledger-price-tracker config

This command shows information about the application's configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
