package n26

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func createRequest(path, deviceToken string, body io.Reader) (*http.Request, error) {
	u, _ := url.ParseRequestURI(apiURL)
	u.Path = path
	urlStr := fmt.Sprintf("%v", u)

	req, err := http.NewRequest("POST", urlStr, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Basic YW5kcm9pZDpzZWNyZXQ=")
	req.Header.Set("device-token", deviceToken)
	return req, nil
}

func (t *Token) GetMFAToken(username, password, deviceToken string) error {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", username)
	data.Set("password", password)

	path := "/oauth2/token/"
	req, err := createRequest(path, deviceToken, strings.NewReader(data.Encode()))
	check(err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	check(err)
	if res.StatusCode != 403 {
		return errors.New("Unexpected response from authentication request")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	check(err)

	err = json.Unmarshal(body, t)
	check(err)
	return nil
}

func (t *Token) requestMfaApproval(deviceToken string) error {
	data, err := json.Marshal(map[string]string{
		"challengeType": "oob",
		"mfaToken":      t.MfaToken,
	})
	check(err)

	path := "/api/mfa/challenge"
	req, err := createRequest(path, deviceToken, bytes.NewBuffer(data))
	check(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36")

	res, err := http.DefaultClient.Do(req)
	check(err)
	if res.StatusCode != 201 {
		fmt.Println(res.StatusCode)
		return errors.New("Failed to request MFA approval")
	}

	// retries 12 times every 5 seconds (60 seconds total wait time)
	// until the login is approved in a authorized device (like the users phone)
	for i := 0; i <= 12; i++ {
		status := t.CompleteMfaApproval(deviceToken)
		if status == 400 {
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
	return nil
}

func (t *Token) CompleteMfaApproval(deviceToken string) int {
	data := url.Values{}
	data.Set("grant_type", "mfa_oob")
	data.Set("mfaToken", t.MfaToken)

	path := "/oauth2/token"
	req, err := createRequest(path, deviceToken, strings.NewReader(data.Encode()))
	check(err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	check(err)
	if res.StatusCode == 400 {
		return res.StatusCode
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	check(err)
	err = json.Unmarshal(body, t)
	check(err)

	return res.StatusCode
}
