# n26
CLI to get information of your N26 account

# Installation
- macOS: Available via Homebrew. Just run `brew install guitmz/tools/n26`
- Linux: You can manually build this project or download a binary release.

You can also install with `go get -u github.com/guitmz/n26` (make sure you have your Go env setup correctly)

# Usage
You can have the `N26_USERNAME` and `N26_PASSWORD` environment variables set to your N26 user email and password. If you don't, you will be prompt for this information, so it's not mandatory.
Example of getting your account balance:
```
+------------------------+-------------+-------------------+----------------+
|          IBAN          |     BIC     | AVAILABLE BALANCE | USABLE BALANCE |
+------------------------+-------------+-------------------+----------------+
| DE74100XXXXXXXXXXXXXXX | NTSXXXXXXXX |              88.8 |           88.8 |
+------------------------+-------------+-------------------+----------------+
```

You can also use the `json` option to output it as JSON with more information:
```
$ n26 balance json
N26 password: ********
{
  "availableBalance": 107.5,
  "usableBalance": 107.5,
  "iban": "DEXXXXXXXXXXXXXX",
  "bic": "NTXXXXXXXXXXX",
  "bankName": "N26 Bank",
  "seized": false,
  "id": "11111-1scasda-1112312-adasdasdasdas"
}
```
You can run `n26 help` for usage description.

# Missing features
- Make a transfer
- Block/unblock cards
- Set card limit
- ?
