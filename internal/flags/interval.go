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

package flags

import (
	"errors"

	"github.com/spf13/cobra"
)

type Interval string

const (
	IntervalDaily   Interval = "daily"
	IntervalWeekly  Interval = "weekly"
	IntervalMonthly Interval = "monthly"
)

// String returns the string representation of the Interval type.
// Used by fmt.Print and Cobra in the help message.
func (i *Interval) String() string {
	return string(*i)
}

// Set sets the value of the Interval type.
// Must have a pointer receiver so it doesn't actually change the value of a copy and not the value itself.
func (i *Interval) Set(value string) error {
	switch value {
	case "daily", "weekly", "monthly":
		*i = Interval(value)
		return nil
	default:
		return errors.New("possible values are \"daily\", \"weekly\", \"monthly\"")
	}
}

// Type is used to describe the expected type for the flag.
func (i *Interval) Type() string {
	return "string"
}

// IntervalCompletion provides completion for the interval flag.
func IntervalCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		"daily\tdaily stock prices",
		"weekly\tweekly stock prices",
		"monthly\tmonthly stock prices",
	}, cobra.ShellCompDirectiveDefault
}
