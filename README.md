[![Build Status](https://travis-ci.org/guitmz/n26.svg?branch=master)](https://travis-ci.org/guitmz/n26) [![Go Report Card](https://goreportcard.com/badge/github.com/guitmz/n26)](https://goreportcard.com/report/github.com/guitmz/n26) [![](https://images.microbadger.com/badges/image/guitmz/n26.svg)](https://microbadger.com/images/guitmz/n26 "Get your own image badge on microbadger.com")

# n26
Go API and CLI to get information of your N26 account

# Installation
- macOS: Available via Homebrew. Just run `brew install guitmz/tools/n26`
- Linux: You can manually build this project or download a binary release.

You can also install with `go get -u github.com/guitmz/n26/cmd/n26` (make sure you have your Go env setup correctly). 

# Docker
A Dockerfile is also provided and the prebuilt image is available for pulling: `docker pull guitmz/n26` or `docker pull guitmz/n26:DESIRED_TAG`

You can run it like:

`$ docker run -e N26_USERNAME="username" -e N26_PASSWORD="password" -e N26_DEVICE_TOKEN="device_token_uuid" guitmz/n26`

or if you want to be asked for your credentials:

`$ docker run -ti -e N26_DEVICE_TOKEN="device_token_uuid" guitmz/n26`

# Authentication
Since 14th of September 2019, N26 requires a login confirmation (2 factor authentication) from the paired phone N26 application to login on devices that are not paired (more details [here](https://n26.com/en-eu/blog/what-is-psd2)). This means you will receive a notification on your phone when you start using this library to request data. This tool checks for your login confirmation every 5 seconds. If you fail to approve the login request within 60 seconds an exception is raised.

### Device Token

Since 17th of June 2020, N26 requires a device_token to differentiate clients. This requires you to specify the `N26_DEVICE_TOKEN` environment variable with an UUID of your choice. Feel free to use any proper UUID generator like https://www.uuidgenerator.net to generate the token.

# Usage
```
NAME:
   N26 - your N26 Bank financial information on the command line

USAGE:
   n26 command [json|csv|statement ID]

VERSION:
   1.5.0

AUTHOR:
   Guilherme Thomazi <thomazi@linux.com>

COMMANDS:
     addresses     addresses linked to your account
     balance       your balance information
     block         blocks a card
     cards         list your cards information
     contacts      your saved contacts
     info          personal information
     limits        your account limits
     spaces        your spaces
     statements    your statements. Passing the statement ID as argument, downloads the PDF to the current directory
     status        general status of your account
     transactions  list your past transactions. Supports CSV output
     unblock       unblocks a card
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
- Improve MFA flow, for now it works but is not really informative
- Make a transfer
- Set card limit
- API docs
- Better error handling
- A terminal UI could also be implemented
- ?

# References
- https://github.com/femueller/python-n26 (MFA reference)
- https://github.com/PierrickP/n26 (API reference)
- https://github.com/Rots (thank you for the PRs!)
