package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/urfave/cli"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (t *Token) requestToken(usr string, pass string) string {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Add("username", usr)
	data.Add("password", pass)

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = "/oauth/token"
	urlStr := fmt.Sprintf("%v", u)

	req, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	req.Header.Add("Authorization", "Basic YW5kcm9pZDpzZWNyZXQ=")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	check(json.Unmarshal(body, t))

	return t.AccessToken
}

func n26Request(t Token, endpoint string, scope interface{}) {
	u, _ := url.ParseRequestURI(apiURL)
	u.Path = endpoint
	urlStr := fmt.Sprintf("%v", u)

	req, _ := http.NewRequest("GET", urlStr, nil)
	req.Header.Add("Authorization", "bearer "+t.AccessToken)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	check(json.Unmarshal(body, &scope))
	response, _ := json.MarshalIndent(&scope, "", "  ")
	fmt.Print(string(response) + "\n")
}

func authentication() Token {
	fmt.Print("N26 password: ")
	pass, err := gopass.GetPasswdMasked()
	check(err)

	token := Token{}
	token.requestToken(os.Getenv("N26_USERNAME"), string(pass))
	return token
}

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Name = "N26"
	app.Usage = "your N26 Bank financial information on the command line"
	app.Author = "Guilherme Thomazi"
	app.Email = "thomazi@linux.com"
	app.Commands = []cli.Command{
		{
			Name:  "balance",
			Usage: "your balance information",
			Action: func(c *cli.Context) error {
				n26Request(authentication(), "/api/accounts", &Balance{})
				return nil
			},
		},
		{
			Name:  "info",
			Usage: "personal information",
			Action: func(c *cli.Context) error {
				n26Request(authentication(), "/api/me", &PersonalInfo{})
				return nil
			},
		},
		{
			Name:  "status",
			Usage: "general status of your account",
			Action: func(c *cli.Context) error {
				n26Request(authentication(), "/api/me/statuses", &Statuses{})
				return nil
			},
		},
		{
			Name:  "addresses",
			Usage: "addresses linked to your account",
			Action: func(c *cli.Context) error {
				n26Request(authentication(), "/api/addresses", &Addresses{})
				return nil
			},
		},
		{
			Name:  "barzahlen",
			Usage: "barzahlen information",
			Action: func(c *cli.Context) error {
				n26Request(authentication(), "/api/barzahlen", &Barzahlen{})
				return nil
			},
		},
		{
			Name:  "cards",
			Usage: "list your cards information",
			Action: func(c *cli.Context) error {
				n26Request(authentication(), "/api/v2/cards", &Cards{})
				return nil
			},
		},
		{
			Name:  "limits",
			Usage: "your account limits",
			Action: func(c *cli.Context) error {
				n26Request(authentication(), "/api/settings/account/limits", &Limits{})
				return nil
			},
		},
		{
			Name:  "contacts",
			Usage: "your saved contacts",
			Action: func(c *cli.Context) error {
				n26Request(authentication(), "/api/smrt/contacts", &Contacts{})
				return nil
			},
		},
		{
			Name:  "transactions",
			Usage: "your past transactions",
			Action: func(c *cli.Context) error {
				n26Request(authentication(), "/api/smrt/transactions", &Transactions{})
				return nil
			},
		},
		{
			Name:  "statements",
			Usage: "your statements",
			Action: func(c *cli.Context) error {
				n26Request(authentication(), "/api/statements", &Statements{})
				return nil
			},
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))
	app.Run(os.Args)
}
