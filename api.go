package n26

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

const apiURL = "https://api.tech26.de"

type Auth struct {
	UserName    string
	Password    string
	DeviceToken string
}

type Balance struct {
	AvailableBalance float64 `json:"availableBalance"`
	UsableBalance    float64 `json:"usableBalance"`
	IBAN             string  `json:"iban"`
	BIC              string  `json:"bic"`
	BankName         string  `json:"bankName"`
	Seized           bool    `json:"seized"`
	ID               string  `json:"id"`
}

type PersonalInfo struct {
	ID                        string `json:"id"`
	Email                     string `json:"email"`
	FirstName                 string `json:"firstName"`
	LastName                  string `json:"lastName"`
	KycFirstName              string `json:"kycFirstName"`
	KycLastName               string `json:"kycLastName"`
	Title                     string `json:"title"`
	Gender                    string `json:"gender"`
	BirthDate                 int64  `json:"birthDate"`
	SignupCompleted           bool   `json:"signupCompleted"`
	Nationality               string `json:"nationality"`
	MobilePhoneNumber         string `json:"mobilePhoneNumber"`
	ShadowUserID              string `json:"shadowUserId"`
	TransferWiseTermsAccepted bool   `json:"transferWiseTermsAccepted"`
	IDNowToken                string `json:"idNowToken"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	MfaToken     string `json:"mfaToken"`
}

type Statuses struct {
	ID                           string `json:"id"`
	Created                      int64  `json:"created"`
	Updated                      int64  `json:"updated"`
	SingleStepSignup             int64  `json:"singleStepSignup"`
	EmailValidationInitiated     int64  `json:"emailValidationInitiated"`
	EmailValidationCompleted     int64  `json:"emailValidationCompleted"`
	ProductSelectionCompleted    int64  `json:"productSelectionCompleted"`
	PhonePairingInitiated        int64  `json:"phonePairingInitiated"`
	PhonePairingCompleted        int64  `json:"phonePairingCompleted"`
	KycInitiated                 int64  `json:"kycInitiated"`
	KycCompleted                 int64  `json:"kycCompleted"`
	KycWebIDInitiated            int64  `json:"kycWebIDInitiated"`
	KycWebIDCompleted            int64  `json:"kycWebIDCompleted"`
	CardActivationCompleted      int64  `json:"cardActivationCompleted"`
	PinDefinitionCompleted       int64  `json:"pinDefinitionCompleted"`
	BankAccountCreationInitiated int64  `json:"bankAccountCreationInitiated"`
	BankAccountCreationSucceded  int64  `json:"bankAccountCreationSucceded"`
	FlexAccount                  bool   `json:"flexAccount"`
}

type Addresses struct {
	Paging struct {
		TotalResults int `json:"totalResults"`
	} `json:"paging"`
	Data []struct {
		AddressLine1     string `json:"addressLine1"`
		StreetName       string `json:"streetName"`
		HouseNumberBlock string `json:"houseNumberBlock"`
		ZipCode          string `json:"zipCode"`
		CityName         string `json:"cityName"`
		CountryName      string `json:"countryName"`
		Type             string `json:"type"`
		ID               string `json:"id"`
	} `json:"data"`
}

type Barzahlen struct {
	DepositAllowance           string `json:"depositAllowance"`
	WithdrawAllowance          string `json:"withdrawAllowance"`
	RemainingAmountMonth       string `json:"remainingAmountMonth"`
	FeeRate                    string `json:"feeRate"`
	Cash26WithdrawalsCount     string `json:"cash26WithdrawalsCount"`
	Cash26WithdrawalsSum       string `json:"cash26WithdrawalsSum"`
	AtmWithdrawalsCount        string `json:"atmWithdrawalsCount"`
	AtmWithdrawalsSum          string `json:"atmWithdrawalsSum"`
	MonthlyDepositFeeThreshold string `json:"monthlyDepositFeeThreshold"`
	Success                    bool   `json:"success"`
}

type Cards []struct {
	ID                                  string      `json:"id"`
	PublicToken                         interface{} `json:"publicToken"`
	Pan                                 interface{} `json:"pan"`
	MaskedPan                           string      `json:"maskedPan"`
	ExpirationDate                      TimeStamp   `json:"expirationDate"`
	CardType                            string      `json:"cardType"`
	Status                              string      `json:"status"`
	CardProduct                         interface{} `json:"cardProduct"`
	CardProductType                     string      `json:"cardProductType"`
	PinDefined                          TimeStamp   `json:"pinDefined"`
	CardActivated                       TimeStamp   `json:"cardActivated"`
	UsernameOnCard                      string      `json:"usernameOnCard"`
	ExceetExpressCardDelivery           interface{} `json:"exceetExpressCardDelivery"`
	Membership                          interface{} `json:"membership"`
	ExceetActualDeliveryDate            interface{} `json:"exceetActualDeliveryDate"`
	ExceetExpressCardDeliveryEmailSent  interface{} `json:"exceetExpressCardDeliveryEmailSent"`
	ExceetCardStatus                    interface{} `json:"exceetCardStatus"`
	ExceetExpectedDeliveryDate          interface{} `json:"exceetExpectedDeliveryDate"`
	ExceetExpressCardDeliveryTrackingID interface{} `json:"exceetExpressCardDeliveryTrackingId"`
	CardSettingsID                      interface{} `json:"cardSettingsId"`
	MptsCard                            bool        `json:"mptsCard"`
}

type Limits []struct {
	Limit  string  `json:"limit"`
	Amount float64 `json:"amount"`
}

type Contacts []struct {
	UserID   string `json:"userId"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	Subtitle string `json:"subtitle"`
	Account  struct {
		AccountType string `json:"accountType"`
		Iban        string `json:"iban"`
		Bic         string `json:"bic"`
	} `json:"account"`
}

type Transactions []struct {
	ID                   string    `json:"id"`
	UserID               string    `json:"userId"`
	Type                 string    `json:"type"`
	Amount               float64   `json:"amount"`
	CurrencyCode         string    `json:"currencyCode"`
	OriginalAmount       float64   `json:"originalAmount,omitempty"`
	OriginalCurrency     string    `json:"originalCurrency,omitempty"`
	ExchangeRate         float64   `json:"exchangeRate,omitempty"`
	MerchantCity         string    `json:"merchantCity,omitempty"`
	VisibleTS            TimeStamp `json:"visibleTS"`
	Mcc                  int       `json:"mcc,omitempty"`
	MccGroup             int       `json:"mccGroup,omitempty"`
	MerchantName         string    `json:"merchantName,omitempty"`
	Recurring            bool      `json:"recurring"`
	AccountID            string    `json:"accountId"`
	Category             string    `json:"category"`
	CardID               string    `json:"cardId,omitempty"`
	UserCertified        TimeStamp `json:"userCertified"`
	Pending              bool      `json:"pending"`
	TransactionNature    string    `json:"transactionNature"`
	CreatedTS            TimeStamp `json:"createdTS"`
	MerchantCountry      int       `json:"merchantCountry,omitempty"`
	SmartLinkID          string    `json:"smartLinkId"`
	LinkID               string    `json:"linkId"`
	Confirmed            TimeStamp `json:"confirmed"`
	PartnerBic           string    `json:"partnerBic,omitempty"`
	PartnerBcn           string    `json:"partnerBcn,omitempty"`
	PartnerAccountIsSepa bool      `json:"partnerAccountIsSepa,omitempty"`
	PartnerName          string    `json:"partnerName,omitempty"`
	PartnerIban          string    `json:"partnerIban,omitempty"`
	PartnerAccountBan    string    `json:"partnerAccountBan,omitempty"`
	ReferenceText        string    `json:"referenceText,omitempty"`
	UserAccepted         int64     `json:"userAccepted,omitempty"`
	SmartContactID       string    `json:"smartContactId,omitempty"`
}

type Statements []struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	VisibleTS int64  `json:"visibleTS"`
	Month     int    `json:"month"`
	Year      int    `json:"year"`
}

type Spaces struct {
	Spaces []struct {
		Balance struct {
			AvailableBalance float64     `json:"availableBalance"`
			OverdraftAmount  interface{} `json:"overdraftAmount"`
		} `json:"balance"`
		Color          string      `json:"color"`
		Goal           interface{} `json:"goal"`
		ID             string      `json:"id"`
		ImageURL       string      `json:"imageUrl"`
		IsCardAttached bool        `json:"isCardAttached"`
		IsPrimary      bool        `json:"isPrimary"`
		Name           string      `json:"name"`
	} `json:"spaces"`
	TotalBalance float64 `json:"totalBalance"`
	UserFeatures struct {
		AvailableSpaces int  `json:"availableSpaces"`
		CanUpgrade      bool `json:"canUpgrade"`
	} `json:"userFeatures"`
}

type Client http.Client

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func NewClient(a Auth) (*Client, error) {
	token := &Token{}
	err := token.GetMFAToken(a.UserName, a.Password, a.DeviceToken)
	check(err)
	err = token.requestMfaApproval(a.DeviceToken)
	check(err)

	tokenSource := &TokenSource{
		AccessToken: token.AccessToken,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	return (*Client)(oauthClient), nil
}

func (c *Client) n26RawRequest(requestMethod, endpoint string, params map[string]string, callback func(io.Reader) error) error {
	var req *http.Request
	var err error

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = endpoint
	u.RawQuery = mapToQuery(params).Encode()

	switch requestMethod {
	case http.MethodGet:
		req, err = http.NewRequest(http.MethodGet, u.String(), nil)
		check(err)
	case http.MethodPost:
		req, err = http.NewRequest(http.MethodPost, u.String(), nil)
		check(err)
	}

	res, err := (*http.Client)(c).Do(req)
	check(err)
	defer res.Body.Close()
	return callback(res.Body)
}

func (c *Client) n26Request(requestMethod, endpoint string, params map[string]string) []byte {
	var body []byte
	err := c.n26RawRequest(requestMethod, endpoint, params, func(r io.Reader) error {
		var err error
		body, err = ioutil.ReadAll(r)
		return err
	})
	check(err)
	return body
}

func mapToQuery(params map[string]string) url.Values {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	return values
}

func (auth *Client) GetBalance(retType string) (string, *Balance) {
	body := auth.n26Request(http.MethodGet, "/api/accounts", nil)
	balance := &Balance{}
	check(json.Unmarshal(body, &balance))
	identedJSON, _ := json.MarshalIndent(&balance, "", "  ")
	if retType == "json" {
		return string(identedJSON), balance
	}
	return "", balance
}

func (auth *Client) GetInfo(retType string) (string, *PersonalInfo) {
	body := auth.n26Request(http.MethodGet, "/api/me", nil)
	info := &PersonalInfo{}
	check(json.Unmarshal(body, &info))
	identedJSON, _ := json.MarshalIndent(&info, "", "  ")
	if retType == "json" {
		return string(identedJSON), info
	}
	return "", info
}

func (auth *Client) GetStatus(retType string) (string, *Statuses) {
	body := auth.n26Request(http.MethodGet, "/api/me/statuses", nil)
	status := &Statuses{}
	check(json.Unmarshal(body, &status))
	identedJSON, _ := json.MarshalIndent(&status, "", "  ")
	if retType == "json" {
		return string(identedJSON), status
	}
	return "", status
}

func (auth *Client) GetAddresses(retType string) (string, *Addresses) {
	body := auth.n26Request(http.MethodGet, "/api/addresses", nil)
	addresses := &Addresses{}
	check(json.Unmarshal(body, &addresses))
	identedJSON, _ := json.MarshalIndent(&addresses, "", "  ")
	if retType == "json" {
		return string(identedJSON), addresses
	}
	return "", addresses
}

func (auth *Client) GetCards(retType string) (string, *Cards) {
	body := auth.n26Request(http.MethodGet, "/api/v2/cards", nil)
	cards := &Cards{}
	check(json.Unmarshal(body, &cards))
	identedJSON, _ := json.MarshalIndent(&cards, "", "  ")
	if retType == "json" {
		return string(identedJSON), cards
	}
	return "", cards
}

func (auth *Client) GetLimits(retType string) (string, *Limits) {
	body := auth.n26Request(http.MethodGet, "/api/settings/account/limits", nil)
	limits := &Limits{}
	check(json.Unmarshal(body, &limits))
	identedJSON, _ := json.MarshalIndent(&limits, "", "  ")
	if retType == "json" {
		return string(identedJSON), limits
	}
	return "", limits
}

func (auth *Client) GetContacts(retType string) (string, *Contacts) {
	body := auth.n26Request(http.MethodGet, "/api/smrt/contacts", nil)
	contacts := &Contacts{}
	check(json.Unmarshal(body, &contacts))
	identedJSON, _ := json.MarshalIndent(&contacts, "", "  ")
	if retType == "json" {
		return string(identedJSON), contacts
	}
	return "", contacts
}

func (auth *Client) GetLastTransactions(limit string) (*Transactions, error) {
	return auth.GetTransactions(TimeStamp{}, TimeStamp{}, limit)
}

// Get transactions for the given time window.
// Use the zero values for the time stamps if no restrictions are
// desired (use the defaults on the server)
func (auth *Client) GetTransactions(from, to TimeStamp, limit string) (*Transactions, error) {
	params := map[string]string{
		"limit": limit,
	}
	//Filter is applied only if both values are set
	if !from.IsZero() && !to.IsZero() {
		params["from"] = fmt.Sprint(from.AsMillis())
		params["to"] = fmt.Sprint(to.AsMillis())
	}
	body := auth.n26Request(http.MethodGet, "/api/smrt/transactions", params)
	transactions := &Transactions{}
	if err := json.Unmarshal(body, &transactions); err != nil {
		return nil, err
	}
	return transactions, nil
}

// Get transactions for the given time window as N26 CSV file. Stored as 'smrt_statement.csv'
func (auth *Client) GetSmartStatementCsv(from, to TimeStamp, reader func(io.Reader) error) error {
	//Filter is applied only if both values are set
	if from.IsZero() || to.IsZero() {
		return errors.New("Start and end time must be set")
	}
	return auth.n26RawRequest(http.MethodGet, fmt.Sprintf("/api/smrt/reports/%v/%v/statements", from.AsMillis(), to.AsMillis()), nil, reader)
}

func (auth *Client) GetStatements(retType string) (string, *Statements) {
	body := auth.n26Request(http.MethodGet, "/api/statements", nil)
	statements := &Statements{}
	check(json.Unmarshal(body, &statements))
	identedJSON, _ := json.MarshalIndent(&statements, "", "  ")
	if retType == "json" {
		return string(identedJSON), statements
	}
	return "", statements
}

func (auth *Client) GetStatementPDF(ID string) {
	body := auth.n26Request(http.MethodGet, fmt.Sprintf("/api/statements/%s", ID), nil)
	ioutil.WriteFile(
		fmt.Sprintf("%s.pdf", ID),
		body,
		0750,
	)
}

func (auth *Client) BlockCard(ID string) {
	_ = auth.n26Request(http.MethodPost, fmt.Sprintf("/api/cards/%s/block", ID), nil)
	fmt.Printf("\nYour card with ID: %s is DISABLED\n\n", ID)
}

func (auth *Client) UnblockCard(ID string) {
	_ = auth.n26Request(http.MethodPost, fmt.Sprintf("/api/cards/%s/unblock", ID), nil)
	fmt.Printf("\nYour card with ID: %s is ACTIVE\n\n", ID)
}

func (auth *Client) GetSpaces(retType string) (string, *Spaces) {
	body := auth.n26Request(http.MethodGet, "/api/spaces", nil)
	spaces := &Spaces{}
	check(json.Unmarshal(body, &spaces))
	identedJSON, _ := json.MarshalIndent(&spaces, "", "  ")
	if retType == "json" {
		return string(identedJSON), spaces
	}
	return "", spaces
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
