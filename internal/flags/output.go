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

package flags

import (
	"errors"

	"github.com/spf13/cobra"
)

type OutputFormat string

const (
	OutputFormatJson       OutputFormat = "json"
	OutputFormatTable      OutputFormat = "table"
	OutputFormatTableShort OutputFormat = "table-short"
	OutputFormatCsv        OutputFormat = "csv"
)

// String returns the string representation of the OutputFormat type.
// Used by fmt.Print and Cobra in the help message.
func (o *OutputFormat) String() string {
	return string(*o)
}

// Set sets the value of the OutputFormat type.
// Must have a pointer receiver so it doesn't actually change the value of a copy and not the value itself.
func (o *OutputFormat) Set(value string) error {
	switch value {
	case "json", "table", "table-short", "csv":
		*o = OutputFormat(value)
		return nil
	default:
		return errors.New("possible values are \"json\", \"table\", \"table-short\", \"csv\"")
	}
}

// Type is used to describe the expected type for the flag.
func (o *OutputFormat) Type() string {
	return "string"
}

func OutputSearchCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		"json\toutput the search results in JSON format",
		"table\toutput the search results in a table format",
		"table-short\toutput the search results in a table format (short version)",
		"csv\toutput the search results in CSV format",
	}, cobra.ShellCompDirectiveDefault
}
