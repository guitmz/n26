[![Build Status](https://travis-ci.org/guitmz/n26.svg?branch=master)](https://travis-ci.org/guitmz/n26)

# n26
Go API and CLI to get information of your N26 account

# Installation
- macOS: Available via Homebrew. Just run `brew install guitmz/n26/n26`
- Linux: You can manually build this project or download a binary release.

You can also install with `go get -u github.com/guitmz/n26/cmd/n26` (make sure you have your Go env setup correctly). 

# Docker
A Dockerfile is also provided and the prebuilt image is available for pulling: `docker pull guitmz/n26`

# Usage
```
NAME:
   N26 - your N26 Bank financial information on the command line

USAGE:
   n26 command [json|statement ID]

VERSION:
   1.3.0

AUTHOR:
   Guilherme Thomazi <thomazi@linux.com>

COMMANDS:
     addresses     addresses linked to your account
     balance       your balance information
     cards         list your cards information
     contacts      your saved contacts
     info          personal information
     limits        your account limits
     statements    your statements
     status        general status of your account
     transactions  your past transactions
     help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

You can have the `N26_USERNAME` and `N26_PASSWORD` environment variables set to your N26 user email and password. If you don't, you will be prompt for this information, so it's not mandatory.
Example of getting your account balance:
```
$ n26 balance
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

And `csv` for transactions.

You can run `n26 help` for usage description.

# Missing features
- Make a transfer
- Block/unblock cards
- Set card limit
- API docs
- ?
