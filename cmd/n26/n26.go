package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"regexp"

	"github.com/guitmz/n26"
	"github.com/howeyc/gopass"
	"github.com/urfave/cli"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func authentication() *n26.Auth {
	username := os.Getenv("N26_USERNAME")
	if username == "" {
		fmt.Print("N26 username: ")
		fmt.Scanln(&username)
	}
	password := os.Getenv("N26_PASSWORD")
	if password == "" {
		fmt.Print("N26 password: ")
		maskedPass, err := gopass.GetPasswdMasked()
		check(err)
		password = string(maskedPass)
	}
	return &n26.Auth{username, password}
}

// Interface for generic data writer that has a header and data table e.g. table writer and csv writer
type dataWriter interface {
	WriteData(header []string, data [][]string) error
}
type transactionWriter interface {
	WriteTransactions(t *n26.Transactions) error
}

func main() {
	app := cli.NewApp()
	app.Version = "1.4.0"
	app.UsageText = "n26 command [json|csv|statement ID]"
	app.Name = "N26"
	app.Usage = "your N26 Bank financial information on the command line"
	app.Author = "Guilherme Thomazi"
	app.Email = "thomazi@linux.com"
	app.Commands = []cli.Command{
		{
			Name:  "balance",
			Usage: "your balance information",
			Action: func(c *cli.Context) error {
				API := authentication()
				prettyJSON, balance := API.GetBalance(c.Args().First())
				if prettyJSON != "" {
					fmt.Println(prettyJSON)
				} else {
					available := strconv.FormatFloat(balance.AvailableBalance, 'f', -1, 64)
					usable := strconv.FormatFloat(balance.UsableBalance, 'f', -1, 64)
					data := [][]string{[]string{balance.IBAN, balance.BIC, available, usable}}
					NewTableWriter().WriteData([]string{"IBAN", "BIC", "Available Balance", "Usable Balance"}, data)
				}
				return nil
			},
		},
		{
			Name:  "info",
			Usage: "personal information",
			Action: func(c *cli.Context) error {
				API := authentication()
				prettyJSON, info := API.GetInfo(c.Args().First())
				if prettyJSON != "" {
					fmt.Println(prettyJSON)
				} else {
					data := [][]string{[]string{fmt.Sprintf("%s %s", info.FirstName, info.LastName), info.Email, info.MobilePhoneNumber}}
					NewTableWriter().WriteData([]string{"Full Name", "Email", "Mobile Phone Number"}, data)
				}
				return nil
			},
		},
		{
			Name:  "status",
			Usage: "general status of your account",
			Action: func(c *cli.Context) error {
				API := authentication()
				prettyJSON, status := API.GetStatus(c.Args().First())
				if prettyJSON != "" {
					fmt.Println(prettyJSON)
				} else {
					data := [][]string{
						[]string{
							time.Unix(status.Created, 0).String(),
						},
					}
					NewTableWriter().WriteData([]string{"Created"}, data)
				}
				return nil
			},
		},
		{
			Name:  "addresses",
			Usage: "addresses linked to your account",
			Action: func(c *cli.Context) error {
				API := authentication()
				prettyJSON, addresses := API.GetAddresses(c.Args().First())
				if prettyJSON != "" {
					fmt.Println(prettyJSON)
				} else {
					data := [][]string{}
					for _, address := range addresses.Data {
						data = append(data,
							[]string{
								fmt.Sprintf("%s %s", address.AddressLine1, address.StreetName),
								address.HouseNumberBlock,
								address.ZipCode,
								address.CityName,
								address.Type,
							},
						)
					}
					NewTableWriter().WriteData([]string{"Address", "Number", "Zipcode", "City", "Type"}, data)
				}
				return nil
			},
		},
		// {
		// 	Name:  "barzahlen",
		// 	Usage: "barzahlen information",
		// 	Action: func(c *cli.Context) error {
		// 		API.n26Request("/api/barzahlen", &Barzahlen{})
		// 		return nil
		// 	},
		// },
		{
			Name:  "cards",
			Usage: "list your cards information",
			Action: func(c *cli.Context) error {
				API := authentication()
				prettyJSON, cards := API.GetCards(c.Args().First())
				if prettyJSON != "" {
					fmt.Println(prettyJSON)
				} else {
					data := [][]string{}
					for _, card := range *cards {
						data = append(data,
							[]string{
								card.UsernameOnCard,
								card.CardType,
								card.CardProductType,
								card.MaskedPan,
							},
						)
					}
					NewTableWriter().WriteData([]string{"Name on Card", "Type", "Product type", "Number"}, data)
				}
				return nil
			},
		},
		{
			Name:  "limits",
			Usage: "your account limits",
			Action: func(c *cli.Context) error {
				API := authentication()
				prettyJSON, limits := API.GetLimits(c.Args().First())
				if prettyJSON != "" {
					fmt.Println(prettyJSON)
				} else {
					data := [][]string{}
					for _, limit := range *limits {
						amount := strconv.FormatFloat(limit.Amount, 'f', -1, 64)
						data = append(data,
							[]string{
								limit.Limit,
								amount,
							},
						)
					}
					NewTableWriter().WriteData([]string{"Limit", "Amount"}, data)
				}
				return nil
			},
		},
		{
			Name:  "contacts",
			Usage: "your saved contacts",
			Action: func(c *cli.Context) error {
				API := authentication()
				prettyJSON, contacts := API.GetContacts(c.Args().First())
				if prettyJSON != "" {
					fmt.Println(prettyJSON)
				} else {
					data := [][]string{}
					for _, contact := range *contacts {
						data = append(data,
							[]string{
								contact.Name,
								contact.Account.Iban,
								contact.Account.Bic,
								contact.Account.AccountType,
							},
						)
					}
					NewTableWriter().WriteData([]string{"Name", "IBAN", "BIC", "Type"}, data)
				}
				return nil
			},
		},
		{
			Name:      "transactions",
			Usage:     "list your past transactions. Supports CSV output",
			ArgsUsage: "[csv|json|table]",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "from", Usage: "retrieve transactions from this date. " +
					"Also 'to' flag needs to be set. Calendar date in the format yyyy-mm-dd. E.g. 2018-03-01"},
				cli.StringFlag{Name: "to", Usage: "retrieve transactions until this date. " +
					"Also 'from' flag needs to be set. Calendar date in the format yyyy-mm-dd. E.g. 2018-03-31"},
			},
			Action: func(c *cli.Context) (err error) {
				const dateFormat = "2006-01-02"
				API := authentication()
				writer, err := getTransactionWriter(c.Args().First())
				check(err)
				var transactions *n26.Transactions
				if c.IsSet("from") && c.IsSet("to") {
					var from, to n26.TimeStamp
					from.Time, err = time.Parse(dateFormat, c.String("from"))
					check(err)
					to.Time, err = time.Parse(dateFormat, c.String("to"))
					check(err)
					transactions, err = API.GetTransactions(from, to)
				} else {
					transactions, err = API.GetLastTransactions()
				}
				check(err)

				err = writer.WriteTransactions(transactions)
				check(err)
				return
			},
		},
		{
			Name:      "statements",
			Usage:     "your statements. Passing the statement ID as argument, downloads the PDF to the current directory",
			ArgsUsage: "[statement ID]",
			Action: func(c *cli.Context) error {
				API := authentication()
				dateRegex := regexp.MustCompile("statement-[0-9][0-9][0-9][0-9]-(1[0-2]|0[1-9]|\\d)")
				argument := c.Args().First()
				switch {
				case dateRegex.MatchString(argument):
					API.GetStatementPDF(argument)
					fmt.Println(fmt.Sprintf("[+] PDF file %s.pdf downloaded!", argument))
				default:
					prettyJSON, statements := API.GetStatements(argument)
					if prettyJSON != "" {
						fmt.Println(prettyJSON)
					} else {
						data := [][]string{}
						for _, statement := range *statements {
							data = append(data,
								[]string{
									statement.ID,
								},
							)
						}
						NewTableWriter().WriteData([]string{"ID"}, data)
					}
				}
				return nil
			},
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))
	app.Run(os.Args)
}

func getTransactionWriter(outType string) (transactionWriter, error) {
	if outType == "json" {
		return jsonWriter{}, nil
	}
	var table dataWriter
	if outType == "csv" {
		var err error
		table, err = NewCsvWriter(os.Stdout)
		if err != nil {
			return nil, err
		}
	} else {
		table = NewTableWriter()
	}
	return transactionToStringWriter{table}, nil
}

type transactionToStringWriter struct {
	out dataWriter
}

func (w transactionToStringWriter) WriteTransactions(transactions *n26.Transactions) error {
	data := [][]string{}
	for _, transaction := range *transactions {
		amount := strconv.FormatFloat(transaction.Amount, 'f', -1, 64)
		var location string
		if transaction.MerchantCity != "" {
			location = transaction.MerchantCity
			if transaction.MerchantCountry != 0 {
				location += ", "
			}
		}
		if transaction.MerchantCountry != 0 {
			location += "Country Code: " + fmt.Sprint(transaction.MerchantCountry)
		}
		data = append(data,
			[]string{
				transaction.VisibleTS.String(),
				transaction.PartnerName,
				transaction.PartnerIban,
				transaction.PartnerBic,
				transaction.MerchantName,
				location,
				amount,
				transaction.CurrencyCode,
				transaction.Type,
			},
		)
	}
	return w.out.WriteData([]string{"Time", "Name", "IBAN", "BIC", "Merchant", "Location", "Amount", "Currency", "Type"},
		data)
}
