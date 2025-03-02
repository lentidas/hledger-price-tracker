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
> **Only the global flags have corresponding settings available in the configuration file. **For any subcommand flag you will need to specify it in the command-line.

> [!WARNING]
> Consider setting the configuration file permissions to read-only for your user (`600`) to avoid leaking your API key. An alternative is to use the environment variable `HPT_API_KEY` to set the API key.

> [!NOTE]
> Only the `api-key` setting can be set through the environment variable `HPT_API_KEY`. **Other settings do not have corresponding environment variables.**
