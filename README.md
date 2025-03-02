# hledger-price-tracker

CLI program to generate market price records for [hledger](https://hledger.org/), using the [Alpha Vantage API](https://www.alphavantage.co/documentation/) and written in [Go](https://go.dev/).

> [!NOTE]
> This is my first Go project. Despite my best efforts, the code may not follow best practices, but I'm open to suggestions and improvements.
> 
> Also note that maybe some of it might be specific to how I use hledger, how I set my commodities, and so on.

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

### `stock search`

This command allows you to search for a stock symbol with a given query. The query can be a part of the company name or the stock symbol itself.

```shell
$ hledger-price-tracker stock search tesco --api-key demo
+---+----------+-----------------------+--------+----------------+----------+-------------+
| # | SYMBOL   | NAME                  | TYPE   | REGION         | CURRENCY | MATCH SCORE |
+---+----------+-----------------------+--------+----------------+----------+-------------+
| 1 | TSCO.LON | Tesco PLC             | Equity | United Kingdom | GBX      | 0.7273      |
| 2 | TSCDF    | Tesco plc             | Equity | United States  | USD      | 0.7143      |
| 3 | TSCDY    | Tesco plc             | Equity | United States  | USD      | 0.7143      |
| 4 | TCO2.FRK | TESCO PLC ADR/1 LS-05 | Equity | Frankfurt      | EUR      | 0.5455      |
| 5 | TCO0.FRK | TESCO PLC LS-0633333  | Equity | Frankfurt      | EUR      | 0.5455      |
+---+----------+-----------------------+--------+----------------+----------+-------------+
```

You can specify a different output format using the `--output` or `-o` flag. The available formats are `table` (default), `table-long` (table with more informations), `json`, and `csv`.
Both `json` and `csv` output nothing more than the raw body of the response from the Alpha Vantage API.

```shell
$ hledger-price-tracker stock search tesco --api-key demo --format table-long
+---+----------+-----------------------+--------+----------------+-------------+--------------+----------+----------+-------------+
| # | SYMBOL   | NAME                  | TYPE   | REGION         | MARKET OPEN | MARKET CLOSE | TIMEZONE | CURRENCY | MATCH SCORE |
+---+----------+-----------------------+--------+----------------+-------------+--------------+----------+----------+-------------+
| 1 | TSCO.LON | Tesco PLC             | Equity | United Kingdom | 08:00       | 16:30        | UTC+01   | GBX      | 0.7273      |
| 2 | TSCDF    | Tesco plc             | Equity | United States  | 09:30       | 16:00        | UTC-04   | USD      | 0.7143      |
| 3 | TSCDY    | Tesco plc             | Equity | United States  | 09:30       | 16:00        | UTC-04   | USD      | 0.7143      |
| 4 | TCO2.FRK | TESCO PLC ADR/1 LS-05 | Equity | Frankfurt      | 08:00       | 20:00        | UTC+02   | EUR      | 0.5455      |
| 5 | TCO0.FRK | TESCO PLC LS-0633333  | Equity | Frankfurt      | 08:00       | 20:00        | UTC+02   | EUR      | 0.5455      |
+---+----------+-----------------------+--------+----------------+-------------+--------------+----------+----------+-------------+

$ hledger-price-tracker stock search tesco --api-key demo --format json
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

$ hledger-price-tracker stock search BA --api-key demo --format csv
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