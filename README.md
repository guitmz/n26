# n26-cli
CLI to get information of your N26 account

# Installation
You can manually build this project or download a binary release.

# Usage
Make sure you have the `N26_USERNAME` environment variable set to your N26 user email.
Example of getting your account balance:
```
$ n26 balance
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
