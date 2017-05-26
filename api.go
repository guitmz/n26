package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const apiURL = "https://api.tech26.de"

type Auth struct {
	UserName string
	Password string
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
	ExpirationDate                      int64       `json:"expirationDate"`
	CardType                            string      `json:"cardType"`
	Status                              string      `json:"status"`
	CardProduct                         interface{} `json:"cardProduct"`
	CardProductType                     string      `json:"cardProductType"`
	PinDefined                          interface{} `json:"pinDefined"`
	CardActivated                       interface{} `json:"cardActivated"`
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
	ID                   string  `json:"id"`
	UserID               string  `json:"userId"`
	Type                 string  `json:"type"`
	Amount               float64 `json:"amount"`
	CurrencyCode         string  `json:"currencyCode"`
	OriginalAmount       float64 `json:"originalAmount,omitempty"`
	OriginalCurrency     string  `json:"originalCurrency,omitempty"`
	ExchangeRate         float64 `json:"exchangeRate,omitempty"`
	MerchantCity         string  `json:"merchantCity,omitempty"`
	VisibleTS            int64   `json:"visibleTS"`
	Mcc                  int     `json:"mcc,omitempty"`
	MccGroup             int     `json:"mccGroup,omitempty"`
	MerchantName         string  `json:"merchantName,omitempty"`
	Recurring            bool    `json:"recurring"`
	AccountID            string  `json:"accountId"`
	Category             string  `json:"category"`
	CardID               string  `json:"cardId,omitempty"`
	UserCertified        int64   `json:"userCertified"`
	Pending              bool    `json:"pending"`
	TransactionNature    string  `json:"transactionNature"`
	CreatedTS            int64   `json:"createdTS"`
	MerchantCountry      int     `json:"merchantCountry,omitempty"`
	SmartLinkID          string  `json:"smartLinkId"`
	LinkID               string  `json:"linkId"`
	Confirmed            int64   `json:"confirmed"`
	PartnerBic           string  `json:"partnerBic,omitempty"`
	PartnerBcn           string  `json:"partnerBcn,omitempty"`
	PartnerAccountIsSepa bool    `json:"partnerAccountIsSepa,omitempty"`
	PartnerName          string  `json:"partnerName,omitempty"`
	PartnerIban          string  `json:"partnerIban,omitempty"`
	PartnerAccountBan    string  `json:"partnerAccountBan,omitempty"`
	ReferenceText        string  `json:"referenceText,omitempty"`
	UserAccepted         int64   `json:"userAccepted,omitempty"`
	SmartContactID       string  `json:"smartContactId,omitempty"`
}

type Statements []struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	VisibleTS int64  `json:"visibleTS"`
	Month     int    `json:"month"`
	Year      int    `json:"year"`
}

func (auth Auth) requestToken() string {
	token := &Token{}
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Add("username", auth.UserName)
	data.Add("password", auth.Password)

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = "/oauth/token"
	urlStr := fmt.Sprintf("%v", u)

	req, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	req.Header.Add("Authorization", "Basic YW5kcm9pZDpzZWNyZXQ=")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	check(err)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	check(json.Unmarshal(body, token))

	return token.AccessToken
}

func (auth Auth) n26Request(endpoint string) []byte {
	u, _ := url.ParseRequestURI(apiURL)
	u.Path = endpoint
	urlStr := fmt.Sprintf("%v", u)

	req, _ := http.NewRequest("GET", urlStr, nil)
	req.Header.Add("Authorization", "bearer "+auth.requestToken())

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}

func (auth Auth) getBalance(retType string) (string, *Balance) {
	body := auth.n26Request("/api/accounts")
	balance := &Balance{}
	check(json.Unmarshal(body, &balance))
	identedJSON, _ := json.MarshalIndent(&balance, "", "  ")
	if retType == "json" {
		return string(identedJSON), balance
	}
	return "", balance
}

func (auth Auth) getInfo(retType string) (string, *PersonalInfo) {
	body := auth.n26Request("/api/me")
	info := &PersonalInfo{}
	check(json.Unmarshal(body, &info))
	identedJSON, _ := json.MarshalIndent(&info, "", "  ")
	if retType == "json" {
		return string(identedJSON), info
	}
	return "", info
}

func (auth Auth) getStatus(retType string) (string, *Statuses) {
	body := auth.n26Request("/api/me/statuses")
	status := &Statuses{}
	check(json.Unmarshal(body, &status))
	identedJSON, _ := json.MarshalIndent(&status, "", "  ")
	if retType == "json" {
		return string(identedJSON), status
	}
	return "", status
}

func (auth Auth) getAddresses(retType string) (string, *Addresses) {
	body := auth.n26Request("/api/addresses")
	addresses := &Addresses{}
	check(json.Unmarshal(body, &addresses))
	identedJSON, _ := json.MarshalIndent(&addresses, "", "  ")
	if retType == "json" {
		return string(identedJSON), addresses
	}
	return "", addresses
}

func (auth Auth) getCards(retType string) (string, *Cards) {
	body := auth.n26Request("/api/v2/cards")
	cards := &Cards{}
	check(json.Unmarshal(body, &cards))
	identedJSON, _ := json.MarshalIndent(&cards, "", "  ")
	if retType == "json" {
		return string(identedJSON), cards
	}
	return "", cards
}

func (auth Auth) getLimits(retType string) (string, *Limits) {
	body := auth.n26Request("/api/settings/account/limits")
	limits := &Limits{}
	check(json.Unmarshal(body, &limits))
	identedJSON, _ := json.MarshalIndent(&limits, "", "  ")
	if retType == "json" {
		return string(identedJSON), limits
	}
	return "", limits
}

func (auth Auth) getContacts(retType string) (string, *Contacts) {
	body := auth.n26Request("/api/smrt/contacts")
	contacts := &Contacts{}
	check(json.Unmarshal(body, &contacts))
	identedJSON, _ := json.MarshalIndent(&contacts, "", "  ")
	if retType == "json" {
		return string(identedJSON), contacts
	}
	return "", contacts
}

func (auth Auth) getTransactions(retType string) (string, *Transactions) {
	body := auth.n26Request("/api/smrt/transactions")
	transactions := &Transactions{}
	check(json.Unmarshal(body, &transactions))
	identedJSON, _ := json.MarshalIndent(&transactions, "", "  ")
	if retType == "json" {
		return string(identedJSON), transactions
	}
	return "", transactions
}

func (auth Auth) getStatements(retType string) (string, *Statements) {
	body := auth.n26Request("/api/statements")
	statements := &Statements{}
	check(json.Unmarshal(body, &statements))
	identedJSON, _ := json.MarshalIndent(&statements, "", "  ")
	if retType == "json" {
		return string(identedJSON), statements
	}
	return "", statements
}

func (auth Auth) getStatementPDF(ID string) {
	body := auth.n26Request(fmt.Sprintf("%s%s", "/api/statements/", ID))
	ioutil.WriteFile(
		fmt.Sprintf("%s.pdf", ID),
		body,
		0750,
	)
}
