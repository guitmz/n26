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
	"github.com/olekukonko/tablewriter"
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

func main() {
	table := tablewriter.NewWriter(os.Stdout)
	app := cli.NewApp()
	app.Version = "1.1.0"
	app.UsageText = "n26 command [json|statement ID]"
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
					table.SetHeader([]string{"IBAN", "BIC", "Available Balance", "Usable Balance"})
					table.AppendBulk(data)
					table.Render()
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
					table.SetHeader([]string{"Full Name", "Email", "Mobile Phone Number"})
					table.AppendBulk(data)
					table.Render()
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
					table.SetHeader([]string{"Created"})
					table.AppendBulk(data)
					table.Render()
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
					table.SetHeader([]string{"Address", "Number", "Zipcode", "City", "Type"})
					table.AppendBulk(data)
					table.Render()
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
					table.SetHeader([]string{"Name on Card", "Type", "Product type", "Number"})
					table.AppendBulk(data)
					table.Render()
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
					table.SetHeader([]string{"Limit", "Amount"})
					table.AppendBulk(data)
					table.Render()
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
					table.SetHeader([]string{"Name", "IBAN", "BIC", "Type"})
					table.AppendBulk(data)
					table.Render()
				}
				return nil
			},
		},
		{
			Name:  "transactions",
			Usage: "your past transactions",
			Action: func(c *cli.Context) error {
				API := authentication()
				prettyJSON, transactions := API.GetTransactions(c.Args().First())
				if prettyJSON != "" {
					fmt.Println(prettyJSON)
				} else {
					data := [][]string{}
					for _, transaction := range *transactions {
						amount := strconv.FormatFloat(transaction.Amount, 'f', -1, 64)
						data = append(data,
							[]string{
								transaction.PartnerName,
								transaction.PartnerIban,
								transaction.PartnerBic,
								amount,
								transaction.CurrencyCode,
								transaction.Type,
							},
						)
					}
					table.SetHeader([]string{"Name", "IBAN", "BIC", "Amount", "Currency", "Type"})
					table.AppendBulk(data)
					table.Render()
				}
				return nil
			},
		},
		{
			Name:  "statements",
			Usage: "your statements",
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
						table.SetHeader([]string{"ID"})
						table.AppendBulk(data)
						table.Render()
					}
				}
				return nil
			},
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))
	app.Run(os.Args)
}
