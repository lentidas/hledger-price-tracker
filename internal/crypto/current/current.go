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

package current

import (
	currencyCurrent "github.com/lentidas/hledger-price-tracker/internal/currency/current"
	"github.com/lentidas/hledger-price-tracker/internal/flags"
)

func Execute(from string, to string, format flags.OutputFormat) (string, error) {
	// The exchange rate is given by the same API function for cryptocurrencies and currencies,
	// so we can use the same function from the analogous module.
	return currencyCurrent.Execute(from, to, format)
}
