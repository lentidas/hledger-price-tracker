# Changelog

## 1.0.0 (2025-03-19)


### ⚠ BREAKING CHANGES

* implement a completely different paradigm

### Features

* add basic documentation for `stock search` command ([7979d44](https://github.com/lentidas/hledger-price-tracker/commit/7979d4411e9180b6b7ba3d52cf5423bc0fe6a10a))
* add error validation and unitary tests to buildSearchURL() ([de3e77b](https://github.com/lentidas/hledger-price-tracker/commit/de3e77b3d47c5cb57a47c5efbfd6ea752e89a590))
* add hledger as a valid output format ([733f8eb](https://github.com/lentidas/hledger-price-tracker/commit/733f8ebc3afaa53c0a396a9e229f2a9902829e11))
* add root command prototype ([c128c93](https://github.com/lentidas/hledger-price-tracker/commit/c128c939113efac35fe93b7086a3d9894a862612))
* add stock and currency palettes with search subcommand ([a8fdcc8](https://github.com/lentidas/hledger-price-tracker/commit/a8fdcc84c52a6bc521146ab608509e8f5f5d8ea0))
* add version subcommand ([302fcb0](https://github.com/lentidas/hledger-price-tracker/commit/302fcb0384fe512e4401ce4abc7571ba3a250050))
* complicate the search structures to make it more generic ([4da887e](https://github.com/lentidas/hledger-price-tracker/commit/4da887eafb61eb160e1241bcf1f16b873a54cf44))
* create command to output configuration ([7e8f089](https://github.com/lentidas/hledger-price-tracker/commit/7e8f089873c79977a8c1d371441141d10060fc9b))
* create common function to perform HTTP requests ([f3e7944](https://github.com/lentidas/hledger-price-tracker/commit/f3e794410ee4a05ca2baae35b974494bd7287d9b))
* create generic and typed JSON body and inherited method ([99ceebd](https://github.com/lentidas/hledger-price-tracker/commit/99ceebd1b9681b6f4d034ef263bd134102e53e58))
* create Go struct to represent JSON response ([c50ca24](https://github.com/lentidas/hledger-price-tracker/commit/c50ca2410abab8da969c4fcf8db68be5fde3396b))
* divide Search() into smaller functions to improve tests ([a6ed1ed](https://github.com/lentidas/hledger-price-tracker/commit/a6ed1edcc7a1ee0e869997bc39d0952dd7d245aa))
* finally embed response inside a generic struct ([abade3e](https://github.com/lentidas/hledger-price-tracker/commit/abade3ed4892382a2b2f1f2cdcaa01ca180d369f))
* finally finished implementing a working version of price ([c290e70](https://github.com/lentidas/hledger-price-tracker/commit/c290e708afa09b48fd657edd31dea0ddbb718c4d))
* greatly simplify code because OOP in Go is a b*tch ([00c9a49](https://github.com/lentidas/hledger-price-tracker/commit/00c9a49b70a9001dc3759c3f691361b01e1b1694))
* group palette commands in a command group for clarity ([13f0cb1](https://github.com/lentidas/hledger-price-tracker/commit/13f0cb130f4955aa7a603c6a2a812d2aae3734c9))
* implement `stock price` subcommand ([c30ef03](https://github.com/lentidas/hledger-price-tracker/commit/c30ef039c24ce2cb059244d1e82c3c5e9b35926e))
* implement a completely different paradigm ([3f9cfa8](https://github.com/lentidas/hledger-price-tracker/commit/3f9cfa8ab28adc704276d5dc185289ab6950d109))
* implement configuration file and setting API key ([38de3e7](https://github.com/lentidas/hledger-price-tracker/commit/38de3e7923ea18b860cfe603c90465fdc92bb8fb))
* implement CSV and table output to console ([3cf989e](https://github.com/lentidas/hledger-price-tracker/commit/3cf989e6bfe3a5559eabf6ad1f7bca4cf9e08cd9))
* implement JSON unmarshalling for price command ([8b85757](https://github.com/lentidas/hledger-price-tracker/commit/8b857572d3e26c6af30278b204e57b0cb142df8d))
* implement output flag and reorganize internal modules ([69c42a2](https://github.com/lentidas/hledger-price-tracker/commit/69c42a2a1f2a49ef274b029699725c2b7c8e1723))
* implement overloading API key from config or env ([9873869](https://github.com/lentidas/hledger-price-tracker/commit/9873869af9510568bd0617b2dc47ae4d6ccc09f9))
* implement the `stock search` command ([dac6296](https://github.com/lentidas/hledger-price-tracker/commit/dac62966ee4f17dc9e0d551b4e01318ca310d62c))
* initial commit ([73d1f81](https://github.com/lentidas/hledger-price-tracker/commit/73d1f81c5487b431b42c65b58209161fb65045df))
* initialize Viper dependency and start configuration ([b20ced4](https://github.com/lentidas/hledger-price-tracker/commit/b20ced4cf3b973e0285c97c2e62a3dcc05fa5d38))
* prepare the stock price structures ([5be2a68](https://github.com/lentidas/hledger-price-tracker/commit/5be2a6860bfb50b742fc308b816fefefbb9c3b39))
* start implementing stock price logic ([c452c78](https://github.com/lentidas/hledger-price-tracker/commit/c452c78d8a11db269d5a4e68dc4880c63cf1f40d))


### Bug Fixes

* cleanup variable declarations after last commit ([f2feb19](https://github.com/lentidas/hledger-price-tracker/commit/f2feb19b482289ebcec46de2ed54cd16d67f7d5e))
* found a way that works for price but with code repetition ([be41055](https://github.com/lentidas/hledger-price-tracker/commit/be410554f98223613617d3784c18a4beca2ec0bd))
* make price structs un-exported ([f7dea00](https://github.com/lentidas/hledger-price-tracker/commit/f7dea00dff9d457174cd19a035d3ae5e19d9ea30))
* make searchResponseTyped implement interface ([12d61b8](https://github.com/lentidas/hledger-price-tracker/commit/12d61b840da3de965722a8deb35caf066747f703))
* remove run from root command (show help automatically) ([f1f39d6](https://github.com/lentidas/hledger-price-tracker/commit/f1f39d69a97a073f13e819bea7f04dc6c45281e1))
* rename entrypoint and add comments ([2f88143](https://github.com/lentidas/hledger-price-tracker/commit/2f88143bcd151b674422646a34eb8814a72ba075))
* use Cobra function for error checking ([61fddf6](https://github.com/lentidas/hledger-price-tracker/commit/61fddf632999ec47dd50e09f1deb605e62ee8cf7))


### Documentation

* add description for the search command ([2232b40](https://github.com/lentidas/hledger-price-tracker/commit/2232b40bd0e2098811d4da08e0266b06608da266))
* add new line in top of long command description ([b620379](https://github.com/lentidas/hledger-price-tracker/commit/b620379af6cf00345f8ad7a467a9bee2ade05ee2))
* add note to link to Stack Overflow answer ([0e7251e](https://github.com/lentidas/hledger-price-tracker/commit/0e7251ed5f35d20dc43e94c05d31f3f8d764319c))
* add usage of `price` subcommand ([800bc6d](https://github.com/lentidas/hledger-price-tracker/commit/800bc6d0fde0a5f12b840255123f765dd2d3e4c7))
* document the configuration behaviour ([bf5f36a](https://github.com/lentidas/hledger-price-tracker/commit/bf5f36aadc726049b952b3350509696a6a694368))
* fix bold in README.md ([c2fb020](https://github.com/lentidas/hledger-price-tracker/commit/c2fb020d6097d0f9afc7555b7e2760104a6cb21b))
* remove unnecessary comment ([0a19653](https://github.com/lentidas/hledger-price-tracker/commit/0a19653fe70942ac4526d7cfd5cf331a23313f28))
* remove whitespace in copyright header ([68f8085](https://github.com/lentidas/hledger-price-tracker/commit/68f808513b77b264898fb911d98bc2b313cb77f3))
* small correction on the copyright message ([fd1bc8c](https://github.com/lentidas/hledger-price-tracker/commit/fd1bc8cb534494e75a3ae58c84c79bdae57a4a65))


### Miscellaneous Chores

* add CODEOWNERS and .renovaterc.json ([f3b32fc](https://github.com/lentidas/hledger-price-tracker/commit/f3b32fc3889e63022e519e7b5e9c9196c3896904))
* add initial code style configuration ([bed7222](https://github.com/lentidas/hledger-price-tracker/commit/bed7222c3424b9f86387c089d49fde7766494e00))
* backup ([41e6a5c](https://github.com/lentidas/hledger-price-tracker/commit/41e6a5c29f4626e96725bc3f6ce1a60ab38e7a98))
* commit before trying getters/setters approach ([a0e7fc1](https://github.com/lentidas/hledger-price-tracker/commit/a0e7fc19e731296f5bb7a3c23c6e58751405240f))
* **deps:** update actions/create-github-app-token action to v1.11.6 ([#11](https://github.com/lentidas/hledger-price-tracker/issues/11)) ([3ef151c](https://github.com/lentidas/hledger-price-tracker/commit/3ef151cadea8a87d3cb5f1c2fc36df3e7b242fa6))
* **deps:** update module github.com/spf13/viper to v1.20.0 ([#9](https://github.com/lentidas/hledger-price-tracker/issues/9)) ([4321ce8](https://github.com/lentidas/hledger-price-tracker/commit/4321ce87c490b4ced63d73c75497e7bf7e7b4278))
* fix CODEOWNERS ([9200500](https://github.com/lentidas/hledger-price-tracker/commit/9200500a5ac75cbd576ba4b9ded790557b640c69))
* merge pull request [#4](https://github.com/lentidas/hledger-price-tracker/issues/4) from lentidas/feat/config-file-api-key ([38de3e7](https://github.com/lentidas/hledger-price-tracker/commit/38de3e7923ea18b860cfe603c90465fdc92bb8fb))
* merge pull request [#6](https://github.com/lentidas/hledger-price-tracker/issues/6) from lentidas/feat/search-command ([dac6296](https://github.com/lentidas/hledger-price-tracker/commit/dac62966ee4f17dc9e0d551b4e01318ca310d62c))
* merge pull request [#7](https://github.com/lentidas/hledger-price-tracker/issues/7) from lentidas/feat/stock-price-command-revamped ([c30ef03](https://github.com/lentidas/hledger-price-tracker/commit/c30ef039c24ce2cb059244d1e82c3c5e9b35926e))
* run `go mod tidy` ([c813d27](https://github.com/lentidas/hledger-price-tracker/commit/c813d27c8a28de2cfb4d5adf55a098357654d44d))


### Code Refactoring

* remove test not implemented ([ad2c249](https://github.com/lentidas/hledger-price-tracker/commit/ad2c2492ac21fd1d1bbd373665d12b9bd2b78ce6))
* rename variable to avoid confusion in search.go ([bdd4510](https://github.com/lentidas/hledger-price-tracker/commit/bdd4510a6afbee3a408ea47aaec1e1258397d272))


### Tests

* implement first unitary tests ([a37639c](https://github.com/lentidas/hledger-price-tracker/commit/a37639c686678c4f3d2d7f7b84f31adc2772f24b))
* organize tests a little bit better and test main.go ([be8da0b](https://github.com/lentidas/hledger-price-tracker/commit/be8da0b057d8563d518b88f3c46c886475d78779))


### Continuous Integration

* add basic GoReleaser configuration and workflow ([6511d3f](https://github.com/lentidas/hledger-price-tracker/commit/6511d3f65de1e2f3ab528624dd3c9d6008ebf6a1))
* add comment ([626a643](https://github.com/lentidas/hledger-price-tracker/commit/626a643e85db34c6c0d13a50b56b1386c6bdcefb))
* add commit linter ([4e35edb](https://github.com/lentidas/hledger-price-tracker/commit/4e35edbbb1a48e5a7d5380ce45d968ea3be46d09))
* add custom Renovate configuration ([980049a](https://github.com/lentidas/hledger-price-tracker/commit/980049ad089c05ad80c89a8be6afbee1e801861d))
* add Semantic Release workflow ([3a5420f](https://github.com/lentidas/hledger-price-tracker/commit/3a5420fa069a80b0c5b79f345dd4dc51052bf986))
* add token to check out step ([1f55e6a](https://github.com/lentidas/hledger-price-tracker/commit/1f55e6a3bc12ceafafd12772d18164fa0f83f71a))
* add way to launch Semantic Release manually ([fa71d3f](https://github.com/lentidas/hledger-price-tracker/commit/fa71d3f85397174764135ef1d3dd64425141cc52))
* consolidate release in a single workflow ([9843b70](https://github.com/lentidas/hledger-price-tracker/commit/9843b705919adbb6f70aab59385df953d9206bd0))
* create workaround to get release notes ([008b9e5](https://github.com/lentidas/hledger-price-tracker/commit/008b9e5b96861e1c44b013a971abb669146c5eec))
* dynamically generate the app name and e-mail ([b9602ef](https://github.com/lentidas/hledger-price-tracker/commit/b9602ef01a4720a4a6825476855223316d27ef1f))
* fix syntax ([bdadce1](https://github.com/lentidas/hledger-price-tracker/commit/bdadce16c80089497cc64f95c4c004035544b56d))
* fix the artifact path ([9a77b5e](https://github.com/lentidas/hledger-price-tracker/commit/9a77b5e2886c2f434c47db5eec4f169dba93ad4e))
* give a type to the variable ([4904002](https://github.com/lentidas/hledger-price-tracker/commit/490400248ae40f053d281578b4f648f2b23a6953))
* overload both environment variables ([a39ddb6](https://github.com/lentidas/hledger-price-tracker/commit/a39ddb632e1213f2398908abb9994e24e4fa5b3b))
* remove [skip ci] that was added to the global config ([3638fb9](https://github.com/lentidas/hledger-price-tracker/commit/3638fb9bcfcf310eec7e619b5a99b8c689cb33ae))
* rename env variable ([09fcb7d](https://github.com/lentidas/hledger-price-tracker/commit/09fcb7d237afd8fb4ce4ffd66cb2aba22a0be095))
* try a different trigger ([48cfd70](https://github.com/lentidas/hledger-price-tracker/commit/48cfd7038042da7104badf788816296e4cb9471c))
* try and use the app with the owner set ([ba7a939](https://github.com/lentidas/hledger-price-tracker/commit/ba7a9395dd39e20af7a379ea4b9ae640da7f5864))
* try different approach to trigger a release ([b26ef07](https://github.com/lentidas/hledger-price-tracker/commit/b26ef07c003220bce92c7857f8f7848b0022a3f5))
* try setting up the repository explicitly ([4008e76](https://github.com/lentidas/hledger-price-tracker/commit/4008e7683bcb58101d2f6ef5227ad5e2c4b96bcf))
* try using last working version ([d80c25f](https://github.com/lentidas/hledger-price-tracker/commit/d80c25f7acbeb163dff9aba5551148f0206d532b))
* use another environment variable ([d469d5f](https://github.com/lentidas/hledger-price-tracker/commit/d469d5f4d1d63d8d38505f22138d94c98882e1ed))
* use application bot email and name for commits ([6ceb9a6](https://github.com/lentidas/hledger-price-tracker/commit/6ceb9a629aca8abb8a8a833c9cf94de0dcf280e4))
* use default GITHUB_TOKEN ([6024560](https://github.com/lentidas/hledger-price-tracker/commit/602456003daa71f00632013ac3da758436354728))
* use Release Please instead of Semantic Release ([90064b8](https://github.com/lentidas/hledger-price-tracker/commit/90064b80594fde3f31e67ce0f46cd4bd4d7f6c24))
* use variables to set application bot email and name ([8ea2479](https://github.com/lentidas/hledger-price-tracker/commit/8ea2479779cd1a65ab9f438ad00525f213da31d6))
