package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"time"

	"regexp"

	"github.com/guitmz/n26"
	"github.com/howeyc/gopass"
	"github.com/urfave/cli"
)

const (
	appVersion = "1.5.2"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func authentication() (*n26.Client, error) {
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
	deviceToken := os.Getenv("N26_DEVICE_TOKEN")
	if deviceToken == "" {
		fmt.Print("N26 device token (must be in uuid format): ")
		fmt.Scanln(&deviceToken)
	}
	return n26.NewClient(n26.Auth{UserName: username, Password: password, DeviceToken: deviceToken})
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
	app.Version = appVersion
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
				API, err := authentication()
				check(err)
				prettyJSON, balance := API.GetBalance(c.Args().First())
				if prettyJSON != "" {
					fmt.Println(prettyJSON)
				} else {
					available := strconv.FormatFloat(balance.AvailableBalance, 'f', -1, 64)
					usable := strconv.FormatFloat(balance.UsableBalance, 'f', -1, 64)
					data := [][]string{{balance.IBAN, balance.BIC, available, usable}}
					NewTableWriter().WriteData([]string{"IBAN", "BIC", "Available Balance", "Usable Balance"}, data)
				}
				return nil
			},
		},
		{
			Name:  "info",
			Usage: "personal information",
			Action: func(c *cli.Context) error {
				API, err := authentication()
				check(err)
				prettyJSON, info := API.GetInfo(c.Args().First())
				if prettyJSON != "" {
					fmt.Println(prettyJSON)
				} else {
					data := [][]string{{fmt.Sprintf("%s %s", info.FirstName, info.LastName), info.Email, info.MobilePhoneNumber}}
					NewTableWriter().WriteData([]string{"Full Name", "Email", "Mobile Phone Number"}, data)
				}
				return nil
			},
		},
		{
			Name:  "status",
			Usage: "general status of your account",
			Action: func(c *cli.Context) error {
				API, err := authentication()
				check(err)
				prettyJSON, status := API.GetStatus(c.Args().First())
				if prettyJSON != "" {
					fmt.Println(prettyJSON)
				} else {
					data := [][]string{
						{
							time.Unix(status.Created/1000, 0).String(),
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
				API, err := authentication()
				check(err)
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
				API, err := authentication()
				check(err)
				prettyJSON, cards := API.GetCards(c.Args().First())
				if prettyJSON != "" {
					fmt.Println(prettyJSON)
				} else {
					data := [][]string{}
					for _, card := range *cards {
						data = append(data,
							[]string{
								card.ID,
								card.UsernameOnCard,
								card.CardType,
								card.CardProductType,
								card.MaskedPan,
								card.ExpirationDate.String(),
								card.Status,
							},
						)
					}
					NewTableWriter().WriteData([]string{"ID", "Name on Card", "Type", "Product type", "Number", "Expiration Date", "Status"}, data)
				}
				return nil
			},
		},
		{
			Name:  "limits",
			Usage: "your account limits",
			Action: func(c *cli.Context) error {
				API, err := authentication()
				check(err)
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
				API, err := authentication()
				check(err)
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
			Usage:     "list your past transactions. Supports CSV output.",
			ArgsUsage: "[csv|json|table|smartcsv]",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "limit", Value: "10", Usage: "retrieve last N transactions. Default to 10."},
				cli.StringFlag{Name: "from", Usage: "retrieve transactions from this date. " +
					"Also 'to' flag needs to be set. Calendar date in the format yyyy-mm-dd. E.g. 2018-03-01"},
				cli.StringFlag{Name: "to", Usage: "retrieve transactions until this date. " +
					"Also 'from' flag needs to be set. Calendar date in the format yyyy-mm-dd. E.g. 2018-03-31"},
			},
			Action: func(c *cli.Context) (err error) {
				const dateFormat = "2006-01-02"
				var from, to n26.TimeStamp
				if c.IsSet("from") {
					from.Time, err = time.Parse(dateFormat, c.String("from"))
					check(err)
				}
				if c.IsSet("to") {
					to.Time, err = time.Parse(dateFormat, c.String("to"))
					check(err)
				}
				API, err := authentication()
				check(err)

				if c.Args().First() == "smartcsv" {
					if from.IsZero() || to.IsZero() {
						fmt.Println("Start and end time must be set for smart CSV!")
						return nil
					}
					err = API.GetSmartStatementCsv(from, to, func(r io.Reader) error {
						_, err := io.Copy(os.Stdout, r)
						return err
					})
					return
				}
				writer, err := getTransactionWriter(c.Args().First())
				check(err)
				limit := c.String("limit")
				var transactions *n26.Transactions
				if !from.IsZero() && !to.IsZero() {
					transactions, err = API.GetTransactions(from, to, limit)
				} else {
					transactions, err = API.GetLastTransactions(limit)
				}
				check(err)

				err = writer.WriteTransactions(transactions)
				check(err)
				return
			},
		},
		{
			Name:      "statements",
			Usage:     "your statements. Passing one or more space separated statement IDs as argument, downloads the PDF to the current directory",
			ArgsUsage: "[statement ID]",
			Action: func(c *cli.Context) error {
				API, err := authentication()
				check(err)
				dateRegex := regexp.MustCompile("statement-[0-9][0-9][0-9][0-9]-(1[0-2]|0[1-9]|\\d)")
				for _, argument := range c.Args() {
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
				}
				return nil
			},
		},
		{
			Name:      "block",
			Usage:     "blocks a card",
			ArgsUsage: "[card ID]",
			Action: func(c *cli.Context) error {
				API, err := authentication()
				check(err)
				API.BlockCard(c.Args().First())
				return nil
			},
		},
		{
			Name:      "unblock",
			Usage:     "unblocks a card",
			ArgsUsage: "[card ID]",
			Action: func(c *cli.Context) error {
				API, err := authentication()
				check(err)
				API.UnblockCard(c.Args().First())
				return nil
			},
		},
		{
			Name:  "spaces",
			Usage: "your spaces",
			Action: func(c *cli.Context) error {
				API, err := authentication()
				check(err)
				prettyJSON, spaces := API.GetSpaces(c.Args().First())
				if prettyJSON != "" {
					fmt.Println(prettyJSON)
				} else {
					data := [][]string{}
					for _, space := range spaces.Spaces {
						data = append(data,
							[]string{
								space.Name,
								strconv.FormatFloat(space.Balance.AvailableBalance, 'f', -1, 64),
							},
						)
					}
					fmt.Printf("\nYour total balance is: %s\n", strconv.FormatFloat(spaces.TotalBalance, 'f', -1, 64))
					fmt.Printf("You still have %d available spaces to create and use\n\n", spaces.UserFeatures.AvailableSpaces)
					NewTableWriter().WriteData([]string{"Name", "Balance"}, data)
				}
				return nil
			},
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))
	err := app.Run(os.Args)
	check(err)
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
