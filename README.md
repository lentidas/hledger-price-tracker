<h1 align="center" style="margin-top: 0px;">hledger-price-tracker</h1>

<div align="center">

![Release](https://img.shields.io/github/v/release/lentidas/hledger-price-tracker?style=for-the-badge) ![License](https://img.shields.io/github/license/lentidas/hledger-price-tracker?style=for-the-badge)
![Tests](https://img.shields.io/github/actions/workflow/status/lentidas/hledger-price-tracker/go-tests.yaml?style=for-the-badge&label=Tests)

</div>

CLI program to generate market price records for [hledger](https://hledger.org/), using the [Alpha Vantage API](https://www.alphavantage.co/documentation/) and written in [Go](https://go.dev/).

> [!NOTE]
> This is my first Go project. Despite my best efforts, the code may not follow best practices, but I'm open to suggestions and improvements.
> 
> Also note that maybe some of it might be specific to how I use hledger, how I set my commodities, and so on.

## Table of Contents

- [Table of Contents](#table-of-contents)
- [Configuration](#configuration)
- [Usage](#usage)
  - [`currency`](#currency)
    - [`currency list`](#currency-list)
    - [`currency current`](#currency-current)
    - [`currency rate`](#currency-rate)
  - [`crypto`](#crypto)
    - [`crypto list`](#crypto-list)
    - [`crypto current`](#crypto-current)
  - [`stock`](#stock)
    - [`stock search`](#stock-search)
    - [`stock price`](#stock-price)

## Configuration

The program reads a configuration file in the YAML format. The default path is `~/.config/hledger-price-tracker/config.yaml`, but you can specify a different path using the `--config` or `-c` flag. The configuration file is **optional**, and you can accomplish the same behaviour through command-line flags.

The precedence order is: command-line flags -> environment variables -> configuration file.

A simple configuration file should look like this:

```yaml
api-key: YOUR_API_KEY
default-currency: EUR
```

> [!IMPORTANT]
> **Only the global flags have corresponding settings available in the configuration file.** For any subcommand flag you will need to specify it in the command-line.

> [!WARNING]
> Consider setting the configuration file permissions to read-only for your user (`600`) to avoid leaking your API key. An alternative is to use the environment variable `HPT_API_KEY` to set the API key.

> [!NOTE]
> Only the `api-key` setting can be set through the environment variable `HPT_API_KEY`. **Other settings do not have corresponding environment variables.**

## Usage

### `currency`

#### `currency list`

This command lists all the physical currencies available on the Alpha Vantage API. It was implemented to allow the user to check what is the 3-letter code for the currency they want to use, since this code is what the API expects.

```shell
hledger-price-tracker currency list
```
```
┌──────┬─────────────────────────────────────┐
│ CODE │ CURRENCY NAME                       │
├──────┼─────────────────────────────────────┤
│ AED  │ United Arab Emirates Dirham         │
│ YER  │ Yemeni Rial                         │
│ ...  │ ...                                 │
│ ZMW  │ Zambian Kwacha                      │
│ ZWL  │ Zimbabwean Dollar                   │
└──────┴─────────────────────────────────────┘
```

> [!NOTE]
> Actually, this command does not talk to the API directly, it simply parses [this CSV file](https://www.alphavantage.co/physical_currency_list/) from Alpha Vantage's documentation.

> [!TIP]
> You can use the `--format` or `-f` flag to get the same CSV output. *Other formats are not available.*

#### `currency current`

This command is used to get the current exchange rate between two currencies.

It expects at least one argument, which is the currency for which you want to get the price to.

The second argument is the currency you want to convert to. It defaults to `EUR` or the currency given to the `--currency` flag (the flag is always overridden by the second argument).

> [!NOTE]
> The commands `currency current` and `crypto current` are interchangeable, as there is a single API endpoint for getting the current exchange rate between two currencies, either physical or digital ones.

```shell
hledger-price-tracker currency current USD JPY --api-key demo
```
```
P 2025-04-05 "USD" 146.94 "JPY"
```

You can specify a different output format using the `--format` or `-f` flag. The available formats are `hledger` (default), `table`, `table-long` (table with more information), and `json`.
The `json` output is nothing more than the raw body of the response from the Alpha Vantage API.

```shell
hledger-price-tracker currency current USD JPY --api-key demo --format table
```
```
┌─────┬─────────┬─────────────────────┐
│ USD │     JPY │ LAST REFRESHED      │
├─────┼─────────┼─────────────────────┤
│ 1   │ 146.935 │ 2025-04-05 18:50:34 │
└─────┴─────────┴─────────────────────┘
```

#### `currency rate`

This command gets the historical exchange rates between two physical currencies, either daily, weekly, or monthly. The default interval is weekly and the default output is given in hledger syntax.

The following command gets the weekly exchange rates from EUR to USD.

```shell
hledger-price-tracker currency rate EUR USD --api-key demo
```
```
P 2015-04-12 "EUR" 1.06 "USD"
P 2015-04-19 "EUR" 1.08 "USD"
P 2015-04-26 "EUR" 1.09 "USD"

...

P 2025-03-16 "EUR" 1.09 "USD"
P 2025-03-23 "EUR" 1.08 "USD"
P 2025-03-30 "EUR" 1.08 "USD"
P 2025-04-04 "EUR" 1.10 "USD"
```

There is also the `--format` or `-f` flag to specify a different output format. The available formats are `hledger` (default), `table`, `table-long`, `json`, and `csv`.

The `hledger` format always uses the closing price of the day. The other formats show more information.

The following example shows the exchange rates from EUR to USD stock in the weekly interval for the first months of 2025, specified using the `--begin` and `--end` flags.

```shell
hledger-price-tracker currency rate EUR USD --api-key demo --format table --begin 2025-01-01 --end 2025-03-31
```
```
┌──────┬──────────┬─────────────────────┐
│ FROM │ CURRENCY │ LAST REFRESHED      │
├──────┼──────────┼─────────────────────┤
│ EUR  │ USD      │ 2025-04-04 21:25:00 │
└──────┴──────────┴─────────────────────┘
┌────────────┬──────┬──────┬──────┬───────┐
│ DATE       │ OPEN │ HIGH │ LOW  │ CLOSE │
├────────────┼──────┼──────┼──────┼───────┤
│ 2025-01-05 │ 1.04 │ 1.04 │ 1.02 │ 1.03  │
│ 2025-01-12 │ 1.04 │ 1.04 │ 1.02 │ 1.02  │
│ 2025-01-19 │ 1.03 │ 1.04 │ 1.02 │ 1.03  │
│ 2025-01-26 │ 1.04 │ 1.05 │ 1.03 │ 1.05  │
│ 2025-02-02 │ 1.04 │ 1.05 │ 1.02 │ 1.02  │
│ 2025-02-09 │ 1.03 │ 1.04 │ 1.03 │ 1.03  │
│ 2025-02-16 │ 1.03 │ 1.05 │ 1.03 │ 1.05  │
│ 2025-02-23 │ 1.05 │ 1.05 │ 1.04 │ 1.05  │
│ 2025-03-02 │ 1.05 │ 1.05 │ 1.04 │ 1.04  │
│ 2025-03-09 │ 1.05 │ 1.09 │ 1.05 │ 1.09  │
│ 2025-03-16 │ 1.08 │ 1.09 │ 1.08 │ 1.09  │
│ 2025-03-23 │ 1.09 │ 1.10 │ 1.08 │ 1.08  │
│ 2025-03-30 │ 1.08 │ 1.08 │ 1.07 │ 1.08  │
└────────────┴──────┴──────┴──────┴───────┘
```

### `crypto`

#### `crypto list`

This command lists all the digital currencies available on the Alpha Vantage API. It was implemented as a way to check what is the 3-letter code for the currency the user wants to use, since this code is what the API expects.

```shell
hledger-price-tracker crypto list
```
```
┌───────────┬──────────────────────────────────┐
│ CODE      │ CURRENCY NAME                    │
├───────────┼──────────────────────────────────┤
│ 1ST       │ FirstBlood                       │
│ 2GIVE     │ GiveCoin                         │
│ ...       │ ...                              │
│ ZLA       │ Zilla                            │
│ ZRX       │ 0x                               │
└───────────┴──────────────────────────────────┘
```

> [!NOTE]
> Actually, this command does not talk to the API directly, it simply parses [this CSV file](https://www.alphavantage.co/digital_currency_list/) from Alpha Vantage's documentation.

> [!TIP]
> You can use the `--format` or `-f` flag to get the same CSV output. *Other formats are not available.*

#### `crypto current`

This command is used to get the current exchange rate between the specified cryptocurrency and currency.

It expects at least one argument, which is the cryptocurrency for which you want to get the price to.

The second argument is the currency you want to convert to, for example, USD or EUR. It defaults to `EUR` or the currency given to the `--currency` flag (the flag is always overridden by the second argument).

> [!NOTE]
> The commands `crypto current` and `currency current` are interchangeable, as there is a single API endpoint for getting the current exchange rate between two currencies, either physical or digital ones.

```shell
hledger-price-tracker crypto current BTC --api-key demo
```
```
P 2025-04-05 "BTC" 75670.94 "EUR"
```

You can specify a different output format using the `--format` or `-f` flag. The available formats are `hledger` (default), `table`, `table-long` (table with more information), and `json`.
The `json` output is nothing more than the raw body of the response from the Alpha Vantage API.

```shell
hledger-price-tracker crypto current BTC --api-key demo --format table
```
```
┌─────┬──────────┬─────────────────────┐
│ BTC │      EUR │ LAST REFRESHED      │
├─────┼──────────┼─────────────────────┤
│ 1   │ 75682.92 │ 2025-04-05 18:44:29 │
└─────┴──────────┴─────────────────────┘
```

<!--

#### `crypto rate`

TODO

-->

### `stock`

#### `stock search`

This command allows you to search for a stock symbol with a given query. The query can be a part of the company name or the stock symbol itself.

```shell
hledger-price-tracker stock search tesco --api-key demo
```
```
┌───┬──────────┬───────────────────────┬────────┬────────────────┬──────────┬─────────────┐
│ # │ SYMBOL   │ NAME                  │ TYPE   │ REGION         │ CURRENCY │ MATCH SCORE │
├───┼──────────┼───────────────────────┼────────┼────────────────┼──────────┼─────────────┤
│ 1 │ TSCO.LON │ Tesco PLC             │ Equity │ United Kingdom │ GBX      │      72.73% │
│ 2 │ TSCDF    │ Tesco plc             │ Equity │ United States  │ USD      │      71.43% │
│ 3 │ TSCDY    │ Tesco plc             │ Equity │ United States  │ USD      │      71.43% │
│ 4 │ TCO2.FRK │ TESCO PLC ADR/1 LS-05 │ Equity │ Frankfurt      │ EUR      │      54.55% │
│ 5 │ TCO0.FRK │ TESCO PLC LS-0633333  │ Equity │ Frankfurt      │ EUR      │      54.55% │
└───┴──────────┴───────────────────────┴────────┴────────────────┴──────────┴─────────────┘
```

You can specify a different output format using the `--format` or `-f` flag. The available formats are `table` (default), `table-long` (table with more information), `json`, and `csv`.
Both `json` and `csv` output nothing more than the raw body of the response from the Alpha Vantage API.

```shell
hledger-price-tracker stock search tesco --api-key demo --format table-long
```
```
┌───┬──────────┬───────────────────────┬────────┬────────────────┬─────────────┬──────────────┬──────────┬──────────┬─────────────┐
│ # │ SYMBOL   │ NAME                  │ TYPE   │ REGION         │ MARKET OPEN │ MARKET CLOSE │ TIMEZONE │ CURRENCY │ MATCH SCORE │
├───┼──────────┼───────────────────────┼────────┼────────────────┼─────────────┼──────────────┼──────────┼──────────┼─────────────┤
│ 1 │ TSCO.LON │ Tesco PLC             │ Equity │ United Kingdom │       08:00 │        16:30 │ UTC+01   │ GBX      │      72.73% │
│ 2 │ TSCDF    │ Tesco plc             │ Equity │ United States  │       09:30 │        16:00 │ UTC-04   │ USD      │      71.43% │
│ 3 │ TSCDY    │ Tesco plc             │ Equity │ United States  │       09:30 │        16:00 │ UTC-04   │ USD      │      71.43% │
│ 4 │ TCO2.FRK │ TESCO PLC ADR/1 LS-05 │ Equity │ Frankfurt      │       08:00 │        20:00 │ UTC+02   │ EUR      │      54.55% │
│ 5 │ TCO0.FRK │ TESCO PLC LS-0633333  │ Equity │ Frankfurt      │       08:00 │        20:00 │ UTC+02   │ EUR      │      54.55% │
└───┴──────────┴───────────────────────┴────────┴────────────────┴─────────────┴──────────────┴──────────┴──────────┴─────────────┘
```
```shell
hledger-price-tracker stock search tesco --api-key demo --format json
```
```
{
    "bestMatches": [
        {
            "1. symbol": "TSCO.LON",
            "2. name": "Tesco PLC",
            "3. type": "Equity",
            "4. region": "United Kingdom",
            "5. marketOpen": "08:00",
            "6. marketClose": "16:30",
            "7. timezone": "UTC+01",
            "8. currency": "GBX",
            "9. matchScore": "0.7273"
        },
        {
            "1. symbol": "TSCDF",
            "2. name": "Tesco plc",
            "3. type": "Equity",
            "4. region": "United States",
            "5. marketOpen": "09:30",
            "6. marketClose": "16:00",
            "7. timezone": "UTC-04",
            "8. currency": "USD",
            "9. matchScore": "0.7143"
        },
        {
            "1. symbol": "TSCDY",
            "2. name": "Tesco plc",
            "3. type": "Equity",
            "4. region": "United States",
            "5. marketOpen": "09:30",
            "6. marketClose": "16:00",
            "7. timezone": "UTC-04",
            "8. currency": "USD",
            "9. matchScore": "0.7143"
        },
        {
            "1. symbol": "TCO2.FRK",
            "2. name": "TESCO PLC ADR/1 LS-05",
            "3. type": "Equity",
            "4. region": "Frankfurt",
            "5. marketOpen": "08:00",
            "6. marketClose": "20:00",
            "7. timezone": "UTC+02",
            "8. currency": "EUR",
            "9. matchScore": "0.5455"
        },
        {
            "1. symbol": "TCO0.FRK",
            "2. name": "TESCO PLC LS-0633333",
            "3. type": "Equity",
            "4. region": "Frankfurt",
            "5. marketOpen": "08:00",
            "6. marketClose": "20:00",
            "7. timezone": "UTC+02",
            "8. currency": "EUR",
            "9. matchScore": "0.5455"
        }
    ]
}
```
```shell
hledger-price-tracker stock search BA --api-key demo --format csv
```
```
symbol,name,type,region,marketOpen,marketClose,timezone,currency,matchScore
BA,Boeing Company,Equity,United States,09:30,16:00,UTC-04,USD,1.0000
BA.LON,BAE Systems plc,Equity,United Kingdom,08:00,16:30,UTC+01,GBX,0.6667
BA05.LON,BA05,Equity,United Kingdom,08:00,16:30,UTC+01,GBP,0.6667
BA29.LON,BA29,Equity,United Kingdom,08:00,16:30,UTC+01,GBP,0.6667
BA69.LON,BA69,Equity,United Kingdom,08:00,16:30,UTC+01,GBP,0.6667
BA3.FRK,Brooks Automation Inc,Equity,Frankfurt,08:00,20:00,UTC+02,EUR,0.5714
BAAPL,null,Equity,United States,09:30,16:00,UTC-04,USD,0.5714
BAAPV,null,Equity,United States,09:30,16:00,UTC-04,USD,0.5714
BAAAAX,Building America Strategy Port CDA USD Ser 21/1Q MNT CASH,Mutual Fund,United States,09:30,16:00,UTC-04,USD,0.5000
BAAAFX,Building America Strgy Portf CDA USD Ser 2022/2Q MNT CASH,Mutual Fund,United States,09:30,16:00,UTC-04,USD,0.5000
```

#### `stock price`

This command allows you to get the price of a stock symbol, either daily, weekly, or monthly. The default interval is weekly and the default output is given in hledger syntax.

The following example shows the price of IBM stock in the weekly interval for the entirety of the available data.

```shell
hledger-price-tracker stock price IBM --api-key demo
```
```
P 1999-11-12 "IBM" 95.87 NIL
P 1999-11-19 "IBM" 103.94 NIL
P 1999-11-26 "IBM" 105.00 NIL
P 1999-12-03 "IBM" 111.87 NIL
P 1999-12-10 "IBM" 109.00 NIL
P 1999-12-17 "IBM" 110.00 NIL
P 1999-12-23 "IBM" 108.62 NIL
P 1999-12-31 "IBM" 107.87 NIL
P 2000-01-07 "IBM" 113.50 NIL
P 2000-01-14 "IBM" 119.62 NIL
P 2000-01-21 "IBM" 121.50 NIL
P 2000-01-28 "IBM" 111.56 NIL
P 2000-02-04 "IBM" 115.62 NIL
P 2000-02-11 "IBM" 115.37 NIL
P 2000-02-18 "IBM" 112.50 NIL
P 2000-02-25 "IBM" 108.00 NIL
P 2000-03-03 "IBM" 108.00 NIL
P 2000-03-10 "IBM" 105.25 NIL
P 2000-03-17 "IBM" 110.00 NIL

...

P 2024-10-18 "IBM" 232.20 NIL
P 2024-10-25 "IBM" 214.67 NIL
P 2024-11-01 "IBM" 208.25 NIL
P 2024-11-08 "IBM" 213.72 NIL
P 2024-11-15 "IBM" 204.99 NIL
P 2024-11-22 "IBM" 222.97 NIL
P 2024-11-29 "IBM" 227.41 NIL
P 2024-12-06 "IBM" 238.04 NIL
P 2024-12-13 "IBM" 230.82 NIL
P 2024-12-20 "IBM" 223.36 NIL
P 2024-12-27 "IBM" 222.78 NIL
P 2025-01-03 "IBM" 222.65 NIL
P 2025-01-10 "IBM" 219.75 NIL
P 2025-01-17 "IBM" 224.79 NIL
P 2025-01-24 "IBM" 224.80 NIL
P 2025-01-31 "IBM" 255.70 NIL
P 2025-02-07 "IBM" 252.34 NIL
P 2025-02-14 "IBM" 261.28 NIL
P 2025-02-21 "IBM" 261.48 NIL
P 2025-02-28 "IBM" 252.44 NIL
P 2025-03-07 "IBM" 261.54 NIL
P 2025-03-12 "IBM" 249.63 NIL
```

> [!NOTE]
> Note the that the currency is displayed as `NIL` when using the `demo` API key. This is normal, and with a real API key, the currency will be displayed correctly.
>
> Also note how the stock symbol is displayed between double quotes. This is to avoid any issues with the hledger syntax, because some stock symbols contain special characters or numbers (p.e. `IBMB34.SAO` for IBM traded in São Paulo's exchange).

> [!IMPORTANT]
> The `--currency` flag has no effect on the output of the `stock price` subcommand, because the currency is defined by the stock symbol itself. For example, `IBM` is traded in USD, and `IBM.FRK` is traded in EUR, despite being the same publicly-traded company.

There is also the `--format` or `-f` flag to specify a different output format. The available formats are `hledger` (default), `table`, `table-long`, `json`, and `csv`.

The `hledger` format always uses the closing price of the stock. The other formats show more information.

The following example shows the price of IBM stock in the weekly interval for the first months of 2025, specified using the `--begin` and `--end` flags.

```shell
hledger-price-tracker stock price IBM --api-key demo --format table --begin 2025-01-01 --end 2025-03-21
```
```
┌────────┬──────────┬────────────────┬────────────┐
│ SYMBOL │ CURRENCY │ LAST REFRESHED │ TIMEZONE   │
├────────┼──────────┼────────────────┼────────────┤
│ IBM    │ NIL      │ 2025-03-12     │ US/Eastern │
└────────┴──────────┴────────────────┴────────────┘
┌────────────┬────────┬────────┬────────┬────────┬──────────┐
│ DATE       │ OPEN   │ HIGH   │ LOW    │ CLOSE  │   VOLUME │
├────────────┼────────┼────────┼────────┼────────┼──────────┤
│ 2025-01-03 │ 220.54 │ 223.66 │ 217.60 │ 222.65 │ 10819153 │
│ 2025-01-10 │ 223.00 │ 226.71 │ 216.80 │ 219.75 │ 12337094 │
│ 2025-01-17 │ 217.89 │ 225.96 │ 214.61 │ 224.79 │ 18990367 │
│ 2025-01-24 │ 224.99 │ 227.45 │ 220.35 │ 224.80 │ 15594637 │
│ 2025-01-31 │ 222.19 │ 261.80 │ 219.84 │ 255.70 │ 39048997 │
│ 2025-02-07 │ 252.40 │ 265.72 │ 251.84 │ 252.34 │ 30149848 │
│ 2025-02-14 │ 250.86 │ 261.94 │ 246.87 │ 261.28 │ 19898073 │
│ 2025-02-21 │ 261.93 │ 265.09 │ 259.83 │ 261.48 │ 18534169 │
│ 2025-02-28 │ 261.50 │ 263.85 │ 246.54 │ 252.44 │ 25541761 │
│ 2025-03-07 │ 254.74 │ 261.96 │ 245.18 │ 261.54 │ 22284160 │
│ 2025-03-12 │ 261.56 │ 266.45 │ 245.53 │ 249.63 │ 17606010 │
└────────────┴────────┴────────┴────────┴────────┴──────────┘
```

> [!NOTE]
> Since the `json` and `csv` outputs are the raw body of the response from the Alpha Vantage API, the `--begin` and `--end` flags do not have any effect.

Alpha Vantage also provides the option to get an adjusted closing price, along with the given dividends. You can use the `--adjusted` or `-a` flag to get this information.

```shell
hledger-price-tracker stock price IBM --api-key demo --format table-long --begin 2024-01-01 --end 2024-12-31 --adjusted
```
```
┌────────┬──────────┬────────────────┬────────────┐
│ SYMBOL │ CURRENCY │ LAST REFRESHED │ TIMEZONE   │
├────────┼──────────┼────────────────┼────────────┤
│ IBM    │ NIL      │ 2025-03-12     │ US/Eastern │
└────────┴──────────┴────────────────┴────────────┘
┌────────────┬────────┬────────┬────────┬────────┬────────────┬──────────┬─────────────────┐
│ DATE       │ OPEN   │ HIGH   │ LOW    │ CLOSE  │ ADJ. CLOSE │   VOLUME │ DIVIDEND AMOUNT │
├────────────┼────────┼────────┼────────┼────────┼────────────┼──────────┼─────────────────┤
│ 2024-01-05 │ 162.83 │ 163.29 │ 158.67 │ 159.16 │ 152.58     │ 14822074 │ 0.00            │
│ 2024-01-12 │ 158.69 │ 165.98 │ 157.88 │ 165.80 │ 158.95     │ 17643392 │ 0.00            │
│ 2024-01-19 │ 165.80 │ 171.58 │ 165.04 │ 171.48 │ 164.39     │ 19864308 │ 0.00            │
│ 2024-01-26 │ 172.82 │ 196.90 │ 172.40 │ 187.42 │ 179.67     │ 56232762 │ 0.00            │
│ 2024-02-02 │ 187.46 │ 189.46 │ 182.71 │ 185.79 │ 178.11     │ 28283876 │ 0.00            │
│ 2024-02-09 │ 185.51 │ 187.18 │ 181.49 │ 186.34 │ 180.25     │ 22784812 │ 1.66            │
│ 2024-02-16 │ 185.90 │ 188.95 │ 182.26 │ 187.64 │ 181.50     │ 21745006 │ 0.00            │
│ 2024-02-23 │ 187.64 │ 188.77 │ 178.75 │ 185.72 │ 179.65     │ 17487852 │ 0.00            │
│ 2024-03-01 │ 185.60 │ 188.38 │ 182.62 │ 188.20 │ 182.05     │ 21955379 │ 0.00            │
│ 2024-03-08 │ 187.76 │ 198.73 │ 187.60 │ 195.95 │ 189.54     │ 29085296 │ 0.00            │
│ 2024-03-15 │ 195.09 │ 199.18 │ 190.70 │ 191.07 │ 184.82     │ 27466323 │ 0.00            │
│ 2024-03-22 │ 191.70 │ 193.98 │ 190.01 │ 190.84 │ 184.60     │ 23968505 │ 0.00            │
│ 2024-03-28 │ 190.26 │ 191.93 │ 188.50 │ 190.96 │ 184.72     │ 15383298 │ 0.00            │
│ 2024-04-05 │ 190.00 │ 193.28 │ 187.34 │ 189.14 │ 182.96     │ 12808073 │ 0.00            │
│ 2024-04-12 │ 189.24 │ 191.25 │ 181.69 │ 182.27 │ 176.31     │ 14955313 │ 0.00            │
│ 2024-04-19 │ 185.57 │ 187.48 │ 180.17 │ 181.58 │ 175.64     │ 16929550 │ 0.00            │
│ 2024-04-26 │ 182.45 │ 184.68 │ 165.66 │ 167.13 │ 161.67     │ 42329269 │ 0.00            │
│ 2024-05-03 │ 167.40 │ 168.22 │ 162.62 │ 165.71 │ 160.29     │ 22536194 │ 0.00            │
│ 2024-05-10 │ 166.50 │ 170.26 │ 165.88 │ 167.15 │ 163.31     │ 17421523 │ 1.67            │
│ 2024-05-17 │ 167.50 │ 169.63 │ 166.48 │ 169.03 │ 165.15     │ 15933303 │ 0.00            │
│ 2024-05-24 │ 169.00 │ 175.46 │ 168.38 │ 170.89 │ 166.96     │ 18410125 │ 0.00            │
│ 2024-05-31 │ 170.44 │ 171.09 │ 163.84 │ 166.85 │ 163.02     │ 15594186 │ 0.00            │
│ 2024-06-07 │ 166.54 │ 171.31 │ 163.53 │ 170.01 │ 166.10     │ 14102396 │ 0.00            │
│ 2024-06-14 │ 169.55 │ 172.47 │ 166.81 │ 169.21 │ 165.32     │ 16222067 │ 0.00            │
│ 2024-06-21 │ 168.76 │ 174.96 │ 167.50 │ 172.46 │ 168.50     │ 21531360 │ 0.00            │
│ 2024-06-28 │ 175.00 │ 178.46 │ 170.41 │ 172.95 │ 168.98     │ 18850478 │ 0.00            │
│ 2024-07-05 │ 173.45 │ 177.98 │ 173.38 │ 176.02 │ 171.97     │  9939255 │ 0.00            │
│ 2024-07-12 │ 176.41 │ 184.16 │ 174.45 │ 182.83 │ 178.63     │ 16071235 │ 0.00            │
│ 2024-07-19 │ 183.38 │ 189.47 │ 181.95 │ 183.25 │ 179.04     │ 17829469 │ 0.00            │
│ 2024-07-26 │ 183.40 │ 196.26 │ 182.86 │ 191.75 │ 187.34     │ 25458498 │ 0.00            │
│ 2024-08-02 │ 193.18 │ 194.55 │ 185.70 │ 189.12 │ 184.77     │ 20594367 │ 0.00            │
│ 2024-08-09 │ 184.55 │ 192.88 │ 181.81 │ 191.45 │ 188.68     │ 18895865 │ 1.67            │
│ 2024-08-16 │ 191.25 │ 194.35 │ 189.00 │ 193.78 │ 190.98     │ 11330854 │ 0.00            │
│ 2024-08-23 │ 193.84 │ 197.92 │ 193.72 │ 196.10 │ 193.26     │ 11022549 │ 0.00            │
│ 2024-08-30 │ 196.00 │ 202.17 │ 195.90 │ 202.13 │ 199.21     │ 15570283 │ 0.00            │
│ 2024-09-06 │ 201.91 │ 205.95 │ 199.34 │ 200.74 │ 197.84     │ 13519865 │ 0.00            │
│ 2024-09-13 │ 201.94 │ 216.09 │ 201.43 │ 214.79 │ 211.68     │ 21518747 │ 0.00            │
│ 2024-09-20 │ 215.88 │ 218.84 │ 210.37 │ 217.70 │ 214.55     │ 28532770 │ 0.00            │
│ 2024-09-27 │ 218.00 │ 224.15 │ 217.27 │ 220.84 │ 217.65     │ 16300165 │ 0.00            │
│ 2024-10-04 │ 220.65 │ 226.08 │ 215.80 │ 226.00 │ 222.73     │ 17778630 │ 0.00            │
│ 2024-10-11 │ 225.38 │ 235.83 │ 225.02 │ 233.26 │ 229.89     │ 18398213 │ 0.00            │
│ 2024-10-18 │ 233.57 │ 237.37 │ 230.17 │ 232.20 │ 228.84     │ 18477394 │ 0.00            │
│ 2024-10-25 │ 231.21 │ 233.34 │ 214.38 │ 214.67 │ 211.57     │ 31380820 │ 0.00            │
│ 2024-11-01 │ 215.50 │ 216.25 │ 203.51 │ 208.25 │ 205.24     │ 26467891 │ 0.00            │
│ 2024-11-08 │ 207.65 │ 216.70 │ 205.57 │ 213.72 │ 210.63     │ 15846890 │ 0.00            │
│ 2024-11-15 │ 214.40 │ 215.41 │ 204.07 │ 204.99 │ 203.63     │ 19438346 │ 1.67            │
│ 2024-11-22 │ 207.00 │ 227.20 │ 205.37 │ 222.97 │ 221.49     │ 21386866 │ 0.00            │
│ 2024-11-29 │ 223.35 │ 230.36 │ 222.65 │ 227.41 │ 225.90     │ 17274177 │ 0.00            │
│ 2024-12-06 │ 227.50 │ 238.38 │ 225.51 │ 238.04 │ 236.46     │ 18743737 │ 0.00            │
│ 2024-12-13 │ 238.00 │ 239.35 │ 227.80 │ 230.82 │ 229.28     │ 20886084 │ 0.00            │
│ 2024-12-20 │ 230.73 │ 231.03 │ 220.03 │ 223.36 │ 221.87     │ 28267440 │ 0.00            │
│ 2024-12-27 │ 222.81 │ 225.40 │ 221.08 │ 222.78 │ 221.30     │  9272351 │ 0.00            │
└────────────┴────────┴────────┴────────┴────────┴────────────┴──────────┴─────────────────┘
```

<!-- TODO Add information on how to contribute to the project, namely how to initialize the repository locally, initialize Go, etc. -->
